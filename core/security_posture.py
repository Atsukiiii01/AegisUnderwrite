import socket

class PostureScanner:
    def __init__(self, target_host):
        self.target = target_host
        # Critical ports that represent insurance-level liability
        self.critical_ports = {
            21: "FTP",
            22: "SSH",
            445: "SMB (WannaCry Risk)",
            3389: "RDP (Ransomware Entry)",
            5432: "PostgreSQL",
            8080: "HTTP-Alt"
        }

    def scan(self):
        findings = []
        print(f"[*] Auditing technical posture for: {self.target}")
        
        for port, service in self.critical_ports.items():
            try:
                with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
                    s.settimeout(1.5)
                    result = s.connect_ex((self.target, port))
                    if result == 0:
                        findings.append({"port": port, "service": service})
            except Exception as e:
                continue
        return findings

    def check_ssl(self):
        try:
            with socket.create_connection((self.target, 443), timeout=2):
                return True
        except:
            return False