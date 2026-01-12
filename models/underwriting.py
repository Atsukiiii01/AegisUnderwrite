from sqlalchemy import Column, Integer, String, DateTime
from datetime import datetime

from core.database import Base

class UnderwritingRun(Base):
    __tablename__ = "underwriting_runs"

    id = Column(Integer, primary_key=True, index=True)
    username = Column(String, index=True)

    industry = Column(String)
    employee_count = Column(Integer)
    data_sensitivity = Column(String)
    previous_incidents = Column(Integer)

    risk_score = Column(Integer)
    security_score = Column(Integer)
    status = Column(String)
    risk_tier = Column(String)
    premium = Column(Integer)
    coverage = Column(Integer)

    created_at = Column(DateTime, default=datetime.utcnow)
