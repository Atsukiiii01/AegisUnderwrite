from fastapi import FastAPI, Request, Form, Depends
from fastapi.responses import HTMLResponse, RedirectResponse, FileResponse
from fastapi.templating import Jinja2Templates
from fastapi.staticfiles import StaticFiles
from sqlalchemy.orm import Session
from jose import jwt, JWTError
from datetime import datetime
import json
import uuid
import os

from reportlab.lib.pagesizes import A4
from reportlab.pdfgen import canvas

from core.database import get_db, engine, Base
from core.models import User, UnderwritingHistory
from core.auth import hash_password, verify_password
from core.risk_engine import RiskAssessmentEngine
from core.security_posture import SecurityPostureEngine
from core.underwriting import UnderwritingEngine

# =====================
# CONFIG
# =====================
SECRET_KEY = "CHANGE_THIS_SECRET"
ALGORITHM = "HS256"

DEFAULT_ADMIN_USERNAME = "admin"
DEFAULT_ADMIN_PASSWORD = "admin123"

PDF_DIR = "generated_pdfs"
os.makedirs(PDF_DIR, exist_ok=True)

# =====================
# APP SETUP
# =====================
app = FastAPI()
app.mount("/static", StaticFiles(directory="frontend/static"), name="static")
templates = Jinja2Templates(directory="frontend/templates")

# =====================
# DB INIT
# =====================
Base.metadata.create_all(bind=engine)

def create_default_admin():
    db = next(get_db())
    if not db.query(User).filter(User.username == DEFAULT_ADMIN_USERNAME).first():
        db.add(User(
            username=DEFAULT_ADMIN_USERNAME,
            password_hash=hash_password(DEFAULT_ADMIN_PASSWORD),
            role="ADMIN",
            is_active=True,
            must_change_password=False
        ))
        db.commit()
    db.close()

create_default_admin()

# =====================
# AUTH HELPERS
# =====================
def create_token(user: User):
    return jwt.encode(
        {"sub": user.username, "role": user.role},
        SECRET_KEY,
        algorithm=ALGORITHM
    )

def get_logged_in_user(request: Request, db: Session):
    token = request.cookies.get("access_token")
    if not token:
        return None
    try:
        payload = jwt.decode(token, SECRET_KEY, algorithms=[ALGORITHM])
        return db.query(User).filter(User.username == payload["sub"], User.is_active == True).first()
    except JWTError:
        return None

# =====================
# ROUTES
# =====================
@app.get("/", response_class=HTMLResponse)
def home(request: Request, db: Session = Depends(get_db)):
    user = get_logged_in_user(request, db)
    if not user:
        return RedirectResponse("/login", 302)
    return templates.TemplateResponse("index.html", {"request": request, "user": user})

# ---------- LOGIN ----------
@app.get("/login", response_class=HTMLResponse)
def login_page(request: Request):
    return templates.TemplateResponse("login.html", {"request": request, "error": None})

@app.post("/login")
def login(request: Request, username: str = Form(...), password: str = Form(...), db: Session = Depends(get_db)):
    user = db.query(User).filter(User.username == username).first()
    if not user or not verify_password(password, user.password_hash):
        return templates.TemplateResponse("login.html", {"request": request, "error": "Invalid credentials"}, 401)
    if not user.is_active:
        return templates.TemplateResponse("login.html", {"request": request, "error": "Account deactivated"}, 403)

    token = create_token(user)
    resp = RedirectResponse("/", 302)
    resp.set_cookie("access_token", token, httponly=True)
    return resp

@app.post("/logout")
def logout():
    r = RedirectResponse("/login", 302)
    r.delete_cookie("access_token")
    return r

# ---------- UNDERWRITE ----------
@app.post("/underwrite", response_class=HTMLResponse)
def underwrite(
    request: Request,
    industry: str = Form(...),
    employee_count: int = Form(...),
    data_sensitivity: str = Form(...),
    previous_incidents: int = Form(...),
    mfa_enabled: bool = Form(False),
    encryption_enabled: bool = Form(False),
    regular_patching: bool = Form(False),
    backups_enabled: bool = Form(False),
    employee_training: bool = Form(False),
    incident_response_plan: bool = Form(False),
    db: Session = Depends(get_db)
):
    user = get_logged_in_user(request, db)
    if not user:
        return RedirectResponse("/login", 302)

    risk = RiskAssessmentEngine(industry, employee_count, data_sensitivity, True, previous_incidents)
    security = SecurityPostureEngine(
        mfa_enabled, encryption_enabled, regular_patching,
        backups_enabled, employee_training, incident_response_plan
    )

    uw = UnderwritingEngine(risk.calculate_risk_score(), security.calculate_security_score())
    decision = uw.underwriting_decision()

    db.add(UnderwritingHistory(
        user_id=user.id,
        industry=industry,
        employee_count=employee_count,
        data_sensitivity=data_sensitivity,
        previous_incidents=previous_incidents,
        risk_score=decision.get("risk_score"),
        security_score=decision.get("security_score"),
        status=decision.get("status"),
        risk_tier=decision.get("risk_tier"),
        annual_premium=decision.get("annual_premium_inr"),
        coverage_limit=decision.get("coverage_limit_inr")
    ))
    db.commit()

    return templates.TemplateResponse("result.html", {"request": request, "decision": decision})

# ---------- PDF DOWNLOAD ----------
@app.post("/underwrite/pdf")
def download_pdf(decision_json: str = Form(...)):
    decision = json.loads(decision_json)

    filename = f"{PDF_DIR}/underwriting_{uuid.uuid4().hex}.pdf"
    c = canvas.Canvas(filename, pagesize=A4)
    width, height = A4

    y = height - 50
    c.setFont("Helvetica-Bold", 16)
    c.drawString(50, y, "AegisUnderwrite – Underwriting Report")

    c.setFont("Helvetica", 11)
    y -= 40
    for key, value in decision.items():
        c.drawString(50, y, f"{key.replace('_',' ').title()}: {value}")
        y -= 20

    c.showPage()
    c.save()

    return FileResponse(filename, media_type="application/pdf", filename="Underwriting_Report.pdf")
