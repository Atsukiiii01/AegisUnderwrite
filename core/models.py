from sqlalchemy import Column, Integer, String, Boolean, DateTime, ForeignKey
from sqlalchemy.orm import relationship
from datetime import datetime
from core.database import Base

class User(Base):
    __tablename__ = "users"

    id = Column(Integer, primary_key=True)
    username = Column(String, unique=True, nullable=False)
    password_hash = Column(String, nullable=False)
    role = Column(String, nullable=False)
    is_active = Column(Boolean, default=True)
    must_change_password = Column(Boolean, default=True)
    created_at = Column(DateTime, default=datetime.utcnow)

    history = relationship("UnderwritingHistory", back_populates="user")


class UnderwritingHistory(Base):
    __tablename__ = "underwriting_history"

    id = Column(Integer, primary_key=True)
    user_id = Column(Integer, ForeignKey("users.id"))

    industry = Column(String)
    employee_count = Column(Integer)
    data_sensitivity = Column(String)
    previous_incidents = Column(Integer)

    risk_score = Column(Integer)
    security_score = Column(Integer)
    status = Column(String)
    risk_tier = Column(String)
    annual_premium = Column(Integer)
    coverage_limit = Column(Integer)

    created_at = Column(DateTime, default=datetime.utcnow)

    user = relationship("User", back_populates="history")
