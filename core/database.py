# core/database.py
import sqlite3
import json
from datetime import datetime

class AegisDB:
    def __init__(self, db_path="aegis.db"):
        self.db_path = db_path
        self._init_db()

    def _init_db(self):
        with sqlite3.connect(self.db_path) as conn:
            cursor = conn.cursor()
            cursor.execute('''
                CREATE TABLE IF NOT EXISTS audits (
                    id INTEGER PRIMARY KEY AUTOINCREMENT,
                    target_url TEXT NOT NULL,
                    industry TEXT,
                    risk_score REAL,
                    grade TEXT,
                    findings TEXT,
                    timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
                )
            ''')
            conn.commit()

    def save_audit(self, target, industry, score, grade, findings):
        with sqlite3.connect(self.db_path) as conn:
            cursor = conn.cursor()
            cursor.execute('''
                INSERT INTO audits (target_url, industry, risk_score, grade, findings, timestamp)
                VALUES (?, ?, ?, ?, ?, ?)
            ''', (target, industry, score, grade, json.dumps(findings), datetime.now()))
            conn.commit()

    def get_history(self, target):
        with sqlite3.connect(self.db_path) as conn:
            cursor = conn.cursor()
            cursor.execute('SELECT * FROM audits WHERE target_url = ? ORDER BY timestamp DESC', (target,))
            return cursor.fetchall()