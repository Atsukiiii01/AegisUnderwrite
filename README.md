# AegisUnderwrite

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![Release](https://img.shields.io/badge/Release-v1.0.0-blue?style=flat-square)

**A concurrent, fault-tolerant OSINT underwriting engine designed for automated attack surface quantification.**

AegisUnderwrite dynamically aggregates and calculates cybersecurity risk across Identity, Domain, and Infrastructure vectors. Built entirely in Go, it features a context-aware intelligence router, native concurrency, and a zero-dependency persistence layer for immutable audit trailing.

---

## Table of Contents
- [System Architecture](#system-architecture)
- [Core Features](#core-features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Disclaimer](#disclaimer)

---

## System Architecture

The engine operates on a fan-out/fan-in concurrency model, bound by strict context timeouts to ensure network latency never blocks execution.

| Component | Function |
| :--- | :--- |
| **Intelligence Router** | Dynamically classifies inputs (Email, Domain, IP) to execute only relevant modules, preventing wasted API calls. |
| **Goroutine Fan-Out** | Triggers simultaneous execution across all valid OSINT providers. |
| **Fault Tolerance** | Isolates module failures. If an API key is missing or a service times out, the engine gracefully captures the error and calculates risk using surviving data. |
| **Risk Aggregation** | Anchors the final risk score to the highest critical vulnerability found, applying a compounding 10% penalty for secondary exposures. |
| **Persistence Layer** | Writes completed scan ledgers to a local `modernc.org/sqlite` database. |

---

## Core Features

* **Zero-Dependency Audit Logging:** Native Go SQLite integration requires no CGO compilation, ensuring the binary remains highly portable across architectures.
* **Modular Provider Interface:** Designed with the Open-Closed Principle. New OSINT sources can be integrated by adhering to a strict `Provider` interface without modifying the core engine.
* **Dynamic Intelligence Routing:** Never sends an IP address to an email breach database.

---

## Installation

Ensure you have [Go](https://golang.org/dl/) 1.21 or later installed.

```bash
# Clone the repository
git clone https://github.com/Atsukiiii01/AegisUnderwrite.git
cd AegisUnderwrite

# Compile the native binary
go build -o aegis main.go
```

---

## Configuration

AegisUnderwrite utilizes environment variables to authenticate with premium threat intelligence providers. 

Export the following keys into your environment. *(Note: If a key is missing, the engine will safely bypass that specific module).*

```bash
export SHODAN_API_KEY="your_shodan_key"
export AEGIS_PREMIUM_IDENTITY_KEY="your_premium_identity_key"
```

---

## Usage

Aegis operates via a strictly typed Command Line Interface. 

### 1. Active Scanning
Initiate a concurrent analysis against a target. The engine will automatically detect the input type (IP, Email, Domain) and route it to the appropriate providers.

```bash
./aegis scan -target 8.8.8.8
```

**Expected Output:**
```json
{
  "target": "8.8.8.8",
  "target_type": "IP",
  "overall_risk": 20,
  "risk_level": "LOW",
  "findings": [
    {
      "provider_name": "infrastructure_shodan",
      "target": "8.8.8.8",
      "raw_data": {
        "cves": null,
        "open_ports": [443, 53],
        "organization": "Google LLC"
      },
      "risk_score": 20
    }
  ]
}
```

### 2. Audit History
Retrieve the historical intelligence ledger for a specific target to analyze risk drift over time.

```bash
./aegis history -target test@company.com
```

**Expected Output:**
```text
[SYSTEM] Retrieving audit history for: test@company.com
[2026-05-17T01:01:06Z] test@company.com | Provider: identity_xposedornot_free | Risk Score: 15
[2026-05-17T01:01:06Z] test@company.com | Provider: mock_breach_db | Risk Score: 85
```

---

## Disclaimer

This tool is designed strictly for automated attack surface quantification, defensive research, and authorized underwriting assessments. Users are responsible for adhering to all applicable laws and terms of service for the integrated APIs. Ensure explicit authorization before actively scanning third-party infrastructure.
