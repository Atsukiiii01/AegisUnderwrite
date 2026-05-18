import subprocess
import sys
import os
from google import genai

def generate_report(target):
    print(f"[*] Executing AegisUnderwrite scan against {target}...")
    
    # 1. Execute the Go binary silently and capture the JSON output
    try:
        result = subprocess.run(
            ["./aegis", "scan", "-target", target, "-silent"],
            capture_output=True,
            text=True,
            check=True
        )
    except subprocess.CalledProcessError as e:
        print(f"[!] Engine failure: {e}")
        sys.exit(1)

    raw_json = result.stdout.strip()
    
    if not raw_json:
        print("[!] Engine returned empty output. Did it compile correctly?")
        sys.exit(1)

    print("[*] Engine completed. Generating CISO-level executive summary via Gemini...")

    # 2. Check for API key (Client automatically picks it up from the environment)
    if not os.environ.get("GEMINI_API_KEY"):
        print("[!] ERROR: GEMINI_API_KEY environment variable is not set.")
        sys.exit(1)
        
    # 3. Construct the prompt
    prompt = f"""
    You are a high-level Cybersecurity Advisor and Underwriter. 
    Review the following raw JSON output from the AegisUnderwrite engine for the target: {target}.

    Raw JSON:
    {raw_json}

    Based strictly on this data, write a brief, authoritative Executive Summary for a Chief Information Security Officer (CISO).
    Do not use generic fluff. Be direct.
    
    Include:
    1. Overall Risk Posture (Based on the overall_risk score).
    2. Primary Attack Vectors (Identify the most critical findings).
    3. Actionable Remediation Steps (Prioritized list to lower the risk score).
    """

    # 4. Generate and print the report using the new SDK architecture
    try:
        client = genai.Client()
        response = client.models.generate_content(
            model='gemini-2.5-flash',
            contents=prompt,
        )
        print("\n" + "="*50)
        print(" AEGIS UNDERWRITE - EXECUTIVE REPORT")
        print("="*50)
        print(response.text)
        print("="*50 + "\n")
    except Exception as e:
        print(f"[!] LLM Generation failed: {e}")

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python3 report.py <target>")
        sys.exit(1)
        
    target = sys.argv[1]
    generate_report(target)