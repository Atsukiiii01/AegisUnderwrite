class UnderwritingEngine:
    def __init__(self, risk_score, security_score):
        self.risk_score = risk_score
        self.security_score = security_score

    def underwriting_decision(self):
        base_premium = 50000 + (self.risk_score * 1000)
        if self.security_score >= 70:
            mod, tier = 0.75, "Preferred"
        elif self.security_score >= 40:
            mod, tier = 1.0, "Standard"
        else:
            mod, tier = 1.25, "High Risk"

        final_premium = int(base_premium * mod)
        coverage = final_premium * 100

        if self.risk_score >= 75 and self.security_score < 40:
            status = "REJECTED"
            exp = ["Risk too high and verifiable security controls insufficient."]
        else:
            status = "APPROVED"
            exp = [f"Approved at {tier} tier based on OSINT security scan."]

        return {
            "status": status, "risk_tier": tier,
            "annual_premium_inr": final_premium, "coverage_limit_inr": coverage,
            "explanation": exp
        }