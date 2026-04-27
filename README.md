# AegisUnderwrite: Telemetry-to-Underwriting Pipeline

AegisUnderwrite functions by replacing trust with objective, infrastructural verification. The core of this system is the automated data pipeline flowing from active target reconnaissance (`PostureScanner`) directly into compounding financial liability calculations (`RiskAssessmentEngine`).

## The Pipeline Architecture

### Phase 1: Objective Telemetry (`PostureScanner`)
Instead of relying on self-reported security questionnaires, the engine initiates active reconnaissance on the target's perimeter.
- **SSL/TLS Verification:** Validates the presence of base-level encryption (Port 443) to mitigate MitM and credential sniffing risks.
- **Critical Port Probing:** Executes dependency-free socket scans targeting specific, high-liability exposures (e.g., RDP/3389, SMB/445, FTP/21, SSH/22).
- **Output:** A structured telemetry payload containing verified technical vulnerabilities.

### Phase 2: Compounding Risk Matrix (`RiskAssessmentEngine`)
The technical telemetry is ingested by the risk engine, which maps technical vulnerabilities to business risk using exponential scaling rather than linear scoring.
- **Base Industry Risk:** Establishes the foundational liability based on the vertical (e.g., Finance, Healthcare).
- **The Truth Penalty:** Maps the raw output of the `PostureScanner` to immediate risk spikes. Missing SSL or exposed RDP triggers critical vulnerability multipliers.
- **Scale & History Multipliers:** The combined base and technical risk is multiplied by the target's attack surface area (employee count) and historical operational security (previous incidents).

### Phase 3: The Financial Orchestrator (`main.py`)
The calculated risk score dictates the automated underwriting decision:
- **Hard Gates:** Terminal risk scores (>80) or missing basic encryption (SSL) result in an immediate `DENIED` status. 
- **Premium Scaling:** Approved targets receive dynamically generated annual INR premiums scaled proportionally to their calculated risk score and data sensitivity coverage limits.
