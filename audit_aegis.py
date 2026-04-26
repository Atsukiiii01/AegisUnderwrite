import os
import re
from datetime import datetime

def audit_codebase(directory="."):
    print(f"--- Initiating Brutal Audit for: {os.path.abspath(directory)} ---\n")
    
    target_extensions = {".py", ".js", ".html", ".css", ".sql", ".sh"}
    ignore_dirs = {".git", "__pycache__", "venv", "env", "node_modules"}
    
    todos = []
    empty_functions = []
    stale_files = []
    
    todo_pattern = re.compile(r"#\s*(TODO|FIXME|HACK):\s*(.*)", re.IGNORECASE)
    empty_func_pattern = re.compile(r"def\s+\w+\(.*?\):\s*(?:\"\"\".*?\"\"\"\s*)?pass", re.DOTALL)

    total_files = 0
    total_lines = 0

    for root, dirs, files in os.walk(directory):
        dirs[:] = [d for d in dirs if d not in ignore_dirs]
        
        for file in files:
            ext = os.path.splitext(file)[1]
            if ext not in target_extensions:
                continue
                
            total_files += 1
            filepath = os.path.join(root, file)
            
            try:
                with open(filepath, "r", encoding="utf-8") as f:
                    content = f.read()
                    lines = content.splitlines()
                    total_lines += len(lines)
                    
                    # Check for TODOs/FIXMEs
                    for i, line in enumerate(lines):
                        match = todo_pattern.search(line)
                        if match:
                            todos.append(f"[{file}:{i+1}] {match.group(1).upper()}: {match.group(2)}")
                    
                    # Check for empty/pass functions (Python specific laziness)
                    if ext == ".py" and empty_func_pattern.search(content):
                        empty_functions.append(f"[{file}] Contains abandoned 'pass' blocks or empty functions.")
                        
            except Exception as e:
                print(f"[!] Error reading {file}: {e}")

    # Output Results
    print(f"Scanned {total_files} files and {total_lines} lines of code.\n")
    
    if todos:
        print("--- OUTSTANDING DEBT (TODOs / FIXMEs) ---")
        for t in todos: print(t)
        print()
    else:
        print("--- OUTSTANDING DEBT ---\nNone found. (Did you even document what you were building?)\n")
        
    if empty_functions:
        print("--- ABANDONED LOGIC ---")
        for e in empty_functions: print(e)
        print()

    print("--- AUDIT COMPLETE ---")
    print("If the above is empty, you don't have a code problem, you have a blank page problem.")

if __name__ == "__main__":
    # Assumes you run this from the AegisUnderwrite directory
    audit_codebase()