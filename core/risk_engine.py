class RiskAssessmentEngine:
    def __init__(
        self,
        industry: str,
        employee_count: int,
        data_sensitivity: str,
        cloud_usage: bool,
        previous_incidents: int
    ):
        self.industry = industry.lower()
        self.employee_count = employee_count
        self.data_sensitivity = data_sensitivity.lower()
        self.cloud_usage = cloud_usage
        self.previous_incidents = previous_incidents

    def calculate_risk_score(self) -> int:
        score = 0

        # Industry risk
        industry_risk = {
            "finance": 25,
            "healthcare": 25,
            "technology": 20,
            "retail": 15,
            "education": 10,
            "manufacturing": 15
        }
        score += industry_risk.get(self.industry, 10)

        # Company size risk
        if self.employee_count > 1000:
            score += 20
        elif self.employee_count > 200:
            score += 15
        else:
            score += 10

        # Data sensitivity risk
        if self.data_sensitivity == "high":
            score += 25
        elif self.data_sensitivity == "medium":
            score += 15
        else:
            score += 5

        # Cloud exposure risk
        if self.cloud_usage:
            score += 10

        # Previous incidents
        score += min(self.previous_incidents * 5, 20)

        return min(score, 100)

    def risk_level(self) -> str:
        score = self.calculate_risk_score()

        if score >= 70:
            return "HIGH"
        elif score >= 40:
            return "MEDIUM"
        else:
            return "LOW"
