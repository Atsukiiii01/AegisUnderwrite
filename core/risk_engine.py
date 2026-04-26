class RiskAssessmentEngine:
    def __init__(self, industry, employee_count, data_sensitivity, previous_incidents, tech_findings=None, has_ssl=True):
        self.industry = industry.lower()
        self.employee_count = employee_count
        self.data_sensitivity = data_sensitivity.lower()
        self.previous_incidents = previous_incidents
        self.tech_findings = tech_findings or []
        self.has_ssl = has_ssl
        self.score = 0
        self.reasons = []

    def calculate_risk_score(self):
        # 1. Base Score (Foundational Risk)
        base_scores = {"healthcare": 25, "finance": 30, "government": 35}
        current_risk = base_scores.get(self.industry, 15)
        self.reasons.append(f"Base Industry Risk ({self.industry}): {current_risk}")
        
        # 2. Technical Penalties (The Truth Layer)
        tech_penalty = 0
        if not self.has_ssl:
            tech_penalty += 15
            self.reasons.append("CRITICAL: No SSL/HTTPS detected: +15")
        
        port_weights = {3389: 30, 445: 25, 21: 15} # Weighted by exploitability
        for finding in self.tech_findings:
            penalty = port_weights.get(finding['port'], 10)
            tech_penalty += penalty
            self.reasons.append(f"CRITICAL PORT OPEN: {finding['service']} ({finding['port']}): +{penalty}")

        # 3. Compounding Multipliers (Scaling Risk)
        emp_multiplier = 1.0 if self.employee_count < 50 else (1.2 if self.employee_count <= 200 else 1.5)
        incident_multiplier = 1.0 + (self.previous_incidents * 0.25) 

        # Final Formula: (Base + Tech) * Scale * History
        final_score = (current_risk + tech_penalty) * emp_multiplier * incident_multiplier
        
        self.score = round(min(final_score, 100), 2)
        return self.score