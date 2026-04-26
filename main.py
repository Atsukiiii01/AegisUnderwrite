from flask import Flask, render_template, request
from core.security_posture import PostureScanner
from core.risk_engine import RiskAssessmentEngine
import os

# Point Flask to your exact frontend structure
template_dir = os.path.abspath('frontend/templates')
static_dir = os.path.abspath('frontend/static')

app = Flask(__name__, template_folder=template_dir, static_folder=static_dir)

@app.route('/')
def index():
    return render_template('index.html')

@app.route('/underwrite', methods=['POST'])
def underwrite():
    # 1. Capture exact fields from your HTML form
    target_domain = request.form.get('domain')
    industry = request.form.get('industry')
    emp_count = int(request.form.get('employee_count', 100))
    sensitivity = request.form.get('data_sensitivity', 'medium')
    incidents = int(request.form.get('previous_incidents', 0))

    if not target_domain:
        return "Error: Target Domain is required.", 400

    # 2. Execute Objective Technical Audit
    scanner = PostureScanner(target_domain)
    tech_findings = scanner.scan()
    has_ssl = scanner.check_ssl()

    # 3. Calculate Risk via Engine
    engine = RiskAssessmentEngine(
        industry=industry,
        employee_count=emp_count,
        data_sensitivity=sensitivity,
        previous_incidents=incidents,
        tech_findings=tech_findings,
        has_ssl=has_ssl
    )
    
    risk_score = engine.calculate_risk_score()
    
    # 4. Filter Explanations for your HTML format
    # The UI wants Business Risk and Security Findings split into two lists
    security_explanation = [r for r in engine.reasons if "PORT" in r or "SSL" in r]
    risk_explanation = [r for r in engine.reasons if "PORT" not in r and "SSL" not in r]
    
    # If there are no technical findings, ensure the UI shows a clean bill of health
    if not security_explanation:
        security_explanation.append("No critical open ports or SSL misconfigurations detected.")

    # Calculate Security Score (Inverse of Risk)
    security_score = max(0, 100 - risk_score)

    # 5. The Financial Underwriting Logic
    # Base coverage starts at 50 Lakhs INR
    base_coverage = 5000000 
    
    # High-risk industries or high sensitivity data require higher coverage caps
    if industry in ['finance', 'healthcare'] or sensitivity == 'high':
        base_coverage = 15000000 # 1.5 Cr INR

    # Base premium is 1.5% of coverage
    base_premium = base_coverage * 0.015 
    
    # Premium scales aggressively based on the risk score calculated by the engine
    final_premium = int(base_premium * (risk_score / 35)) 
    
    # Underwriting Decision Gate: Deny coverage if risk is completely unmanageable
    if risk_score > 80 or not has_ssl:
        status = "DENIED"
        final_premium = 0
        base_coverage = 0
    else:
        status = "APPROVED"

    decision = {
        "status": status,
        "annual_premium_inr": f"{final_premium:,}",
        "coverage_limit_inr": f"{base_coverage:,}"
    }

    # 6. Push to your Result Template
    return render_template(
        'result.html', 
        risk_score=risk_score,
        security_score=security_score,
        decision=decision,
        security_explanation=security_explanation,
        risk_explanation=risk_explanation
    )

if __name__ == '__main__':
    print("[*] AegisUnderwrite Engine Active.")
    print("[*] Access the dashboard at: http://127.0.0.1:5000")
    app.run(debug=True, port=5000)