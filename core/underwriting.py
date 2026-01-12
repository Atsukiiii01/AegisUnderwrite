class UnderwritingEngine:
    def __init__(self, risk_score: int, security_score: int):
        self.risk_score = risk_score
        self.security_score = security_score

    def underwriting_decision(self):
        """
        Decision logic:
        - Combine risk and security realistically
        - Reject only in extreme cases
        """

        # Normalize scores
        effective_risk = self.risk_score - (self.security_score * 0.4)

        # Clamp
        effective_risk = max(0, min(100, effective_risk))

        # =========================
        # REJECTION RULE (RARE)
        # =========================
        if self.risk_score >= 90 and self.security_score <= 20:
            return {
                "status": "REJECTED",
                "risk_score": self.risk_score,
                "security_score": self.security_score,
                "risk_tier": "CRITICAL"
            }

        # =========================
        # RISK TIERS
        # =========================
        if effective_risk <= 30:
            tier = "LOW"
            base_premium = 40000
        elif effective_risk <= 60:
            tier = "MEDIUM"
            base_premium = 70000
        else:
            tier = "HIGH"
            base_premium = 120000

        # =========================
        # PREMIUM ADJUSTMENTS
        # =========================
        premium = int(base_premium + (effective_risk * 500))

        coverage = 10_000_000
        if tier == "HIGH":
            coverage = 5_000_000

        return {
            "status": "APPROVED",
            "risk_score": self.risk_score,
            "security_score": self.security_score,
            "risk_tier": tier,
            "annual_premium_inr": premium,
            "coverage_limit_inr": coverage
        }
