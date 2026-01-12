class SecurityPostureEngine:
    def __init__(
        self,
        mfa_enabled: bool,
        encryption_enabled: bool,
        regular_patching: bool,
        backups_enabled: bool,
        employee_training: bool,
        incident_response_plan: bool
    ):
        self.mfa_enabled = mfa_enabled
        self.encryption_enabled = encryption_enabled
        self.regular_patching = regular_patching
        self.backups_enabled = backups_enabled
        self.employee_training = employee_training
        self.incident_response_plan = incident_response_plan

    def calculate_security_score(self) -> int:
        score = 0

        if self.mfa_enabled:
            score += 20

        if self.encryption_enabled:
            score += 20

        if self.regular_patching:
            score += 15

        if self.backups_enabled:
            score += 15

        if self.employee_training:
            score += 15

        if self.incident_response_plan:
            score += 15

        return min(score, 100)

    def security_level(self) -> str:
        score = self.calculate_security_score()

        if score >= 80:
            return "STRONG"
        elif score >= 50:
            return "MODERATE"
        else:
            return "WEAK"
