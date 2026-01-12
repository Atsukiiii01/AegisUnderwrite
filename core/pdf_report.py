from reportlab.lib.pagesizes import A4
from reportlab.pdfgen import canvas
from datetime import datetime
import os

def generate_underwriting_pdf(data: dict) -> str:
    filename = f"underwriting_{data['username']}_{int(datetime.utcnow().timestamp())}.pdf"
    filepath = os.path.join("reports", filename)

    os.makedirs("reports", exist_ok=True)

    c = canvas.Canvas(filepath, pagesize=A4)
    width, height = A4

    y = height - 50

    c.setFont("Helvetica-Bold", 16)
    c.drawString(50, y, "AegisUnderwrite – Cyber Insurance Report")
    y -= 40

    c.setFont("Helvetica", 11)
    c.drawString(50, y, f"Generated On: {datetime.utcnow()}")
    y -= 30

    fields = [
        ("Username", data["username"]),
        ("Industry", data["industry"]),
        ("Employee Count", data["employee_count"]),
        ("Data Sensitivity", data["data_sensitivity"]),
        ("Previous Incidents", data["previous_incidents"]),
        ("Risk Score", data["risk_score"]),
        ("Security Score", data["security_score"]),
        ("Risk Tier", data["risk_tier"]),
        ("Status", data["status"]),
        ("Annual Premium", f"₹{data['premium']}"),
        ("Coverage Limit", f"₹{data['coverage']}"),
    ]

    for label, value in fields:
        c.drawString(50, y, f"{label}: {value}")
        y -= 20

    c.showPage()
    c.save()

    return filepath
