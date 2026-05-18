# AegisUnderwrite 

![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)
![Python Version](https://img.shields.io/badge/Python-3.10+-3776AB?style=flat-square&logo=python)
![Release](https://img.shields.io/badge/Release-v1.4.0-blue?style=flat-square)

**A concurrent, fault-tolerant OSINT underwriting engine designed for automated attack surface quantification and executive risk translation.**

AegisUnderwrite dynamically aggregates and calculates cybersecurity risk across Identity, Domain, and Infrastructure vectors. Built entirely in Go, it features a context-aware intelligence router, native concurrency, and a zero-dependency persistence layer. The v1.4.0 architecture decouples data collection from analysis, using a Python-based LLM presentation layer to translate mathematically strict risk scores into actionable CISO-level strategy.

---

## Table of Contents
- [System Architecture](#system-architecture)
- [Core Features](#core-features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [AI Executive Wrapper](#ai-executive-wrapper)
- [Disclaimer](#disclaimer)

---

## System Architecture

The engine operates on a strict boundary between the deterministic data layer and the non-deterministic presentation layer. 

| Component | Function |
| :--- | :--- |
| **Data Layer (Go)** | Highly concurrent fan-out/fan-in network execution. Identifies targets, queries OSINT providers, and calculates absolute risk math. |
| **Presentation Layer (Python)** | Ingests the pure JSON payload from the Go engine and utilizes the Gemini API to generate executive summaries and remediation steps. |
| **Fault Tolerance** | Isolates module failures. If an API key is missing or a service times out, the engine gracefully captures the error and calculates risk using surviving data. |
| **Risk Aggregation** | Anchors the final risk score to the highest critical vulnerability found, applying a compounding 10% penalty for secondary exposures (capped at 100). |
| **Persistence Layer** | Writes completed scan ledgers to a local `modernc.org/sqlite` database. |

---

## Core Features (v1.4.0)

* **DNS Security Auditing:** Native resolution of root SPF and DMARC enforcement policies to calculate Business Email Compromise (BEC) and spoofing risk.
* **Passive Subdomain Expansion:** Aggregation of passive DNS records to map internal cloud routing and hidden administrative surfaces.
* **Live CVE Parsing:** Integrates with Shodan to identify unpatched, actively exploitable vulnerabilities mapped directly to target infrastructure.
* **Identity Breach Mapping:** Cross-references target domains against known historical data dumps.
* **Automated CI/CD:** Native GitHub Actions pipeline for automated cross-compilation across macOS, Linux, and Windows upon every release tag.

---

## Installation

Ensure you have [Go](https://golang.org/dl/) 1.21 or later installed. Alternatively, download the pre-compiled binaries from the [Releases](../../releases) page.

```bash
# Clone the repository
git clone https://github.com/Atsukiiii01/AegisUnderwrite.git

cd AegisUnderwrite

# Compile the native binary
go build -o aegis main.go
```

---

## Configuration

AegisUnderwrite utilizes environment variables to authenticate with premium threat intelligence providers and the AI presentation layer. 

Export the following keys into your environment. *(Note: If an OSINT key is missing, the Go engine will safely bypass that specific module).*

```bash
# OSINT Providers
export SHODAN_API_KEY="your_shodan_key"
export AEGIS_PREMIUM_IDENTITY_KEY="your_premium_identity_key"

# AI Executive Wrapper
export GEMINI_API_KEY="your_gemini_api_key"
```

---

## Usage

Aegis operates via a strictly typed Command Line Interface.

### 1. Active Scanning
Initiate a concurrent analysis against a target. The engine will automatically detect the input type (IP, Email, Domain).

```bash
./aegis scan -target tesla.com
```

### 2. Pipeline / JQ Integration (Silent Mode)
Use the `-silent` flag to mute system logs and output pure, parsable JSON for integration with other security tools (e.g., Nuclei, JQ).

```bash
./aegis scan -target tesla.com -silent
```

### 3. Audit History
Retrieve the historical intelligence ledger for a specific target to analyze risk drift over time.

```bash
./aegis history -target tesla.com
```

---

## AI Executive Wrapper

To generate CISO-level text reports, use the Python wrapper. This script automatically runs the Go binary in silent mode, captures the JSON, and orchestrates the Gemini API to write a prioritized executive summary.

**Requirements:**
```bash
pip install google-genai
```

**Execution:**
```bash
python3 report.py tesla.com
```

**Expected Output:**
Generates a prioritized text report detailing *Overall Risk Posture*, *Primary Attack Vectors*, and *Actionable Remediation Steps*.

---

## Disclaimer

This tool is designed strictly for automated attack surface quantification, defensive research, and authorized underwriting assessments. Users are responsible for adhering to all applicable laws and terms of service for the integrated APIs. Ensure explicit authorization before actively scanning third-party infrastructure.