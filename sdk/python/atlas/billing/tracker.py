"""Billing History Tracker"""

import time
from typing import List, Dict, Optional
from dataclasses import dataclass, field
from datetime import datetime
from collections import defaultdict

from .calculator import BillingCalculator, BillingConfig


@dataclass
class RequestRecord:
    request_id: str
    user_id: str
    model_id: str
    model_type: str
    cost: float
    latency_ms: float
    tokens: Optional[int] = None
    timestamp: float = field(default_factory=time.time)
    metadata: Dict = field(default_factory=dict)


class BillingTracker:
    def __init__(self, config: Optional[BillingConfig] = None):
        self.calculator = BillingCalculator(config)
        self.records: List[RequestRecord] = []
        self.user_totals: Dict[str, float] = defaultdict(float)
        self.model_totals: Dict[str, float] = defaultdict(float)
    
    def record_request(
        self,
        request_id: str,
        user_id: str,
        model_id: str,
        model_type: str,
        latency_ms: float,
        tokens: Optional[int] = None,
        input_size: Optional[int] = None,
        metadata: Optional[Dict] = None,
    ) -> RequestRecord:
        cost = self.calculator.calculate_request_cost(
            model_type=model_type,
            latency_ms=latency_ms,
            tokens=tokens,
            input_size=input_size,
        )
        
        record = RequestRecord(
            request_id=request_id,
            user_id=user_id,
            model_id=model_id,
            model_type=model_type,
            cost=cost,
            latency_ms=latency_ms,
            tokens=tokens,
            metadata=metadata or {},
        )
        
        self.records.append(record)
        self.user_totals[user_id] += cost
        self.model_totals[model_id] += cost
        
        return record
    
    def get_user_total(self, user_id: str) -> float:
        return self.user_totals.get(user_id, 0.0)
    
    def get_model_total(self, model_id: str) -> float:
        return self.model_totals.get(model_id, 0.0)
    
    def get_user_records(
        self,
        user_id: str,
        limit: Optional[int] = None,
    ) -> List[RequestRecord]:
        records = [r for r in self.records if r.user_id == user_id]
        if limit:
            records = records[-limit:]
        return records
    
    def get_model_records(
        self,
        model_id: str,
        limit: Optional[int] = None,
    ) -> List[RequestRecord]:
        records = [r for r in self.records if r.model_id == model_id]
        if limit:
            records = records[-limit:]
        return records
    
    def get_total_cost(self) -> float:
        return sum(r.cost for r in self.records)
    
    def get_stats(self) -> Dict:
        return {
            "total_requests": len(self.records),
            "total_cost": self.get_total_cost(),
            "unique_users": len(self.user_totals),
            "unique_models": len(self.model_totals),
            "avg_cost_per_request": (
                self.get_total_cost() / len(self.records) if self.records else 0
            ),
        }


_global_tracker = BillingTracker()


def get_billing_tracker() -> BillingTracker:
    return _global_tracker

