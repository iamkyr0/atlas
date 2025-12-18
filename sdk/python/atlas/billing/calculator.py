"""Billing Cost Calculator"""

from typing import Dict, Optional
from dataclasses import dataclass


@dataclass
class BillingConfig:
    base_cost_per_request: float = 0.001
    cost_per_token: float = 0.0001
    cost_per_second: float = 0.01
    model_multipliers: Dict[str, float] = None
    
    def __post_init__(self):
        if self.model_multipliers is None:
            self.model_multipliers = {
                "llm": 2.0,
                "vision": 1.5,
                "speech": 1.2,
                "embedding": 1.0,
                "generic": 1.0,
            }


class BillingCalculator:
    def __init__(self, config: Optional[BillingConfig] = None):
        self.config = config or BillingConfig()
    
    def calculate_request_cost(
        self,
        model_type: str = "generic",
        latency_ms: float = 0,
        tokens: Optional[int] = None,
        input_size: Optional[int] = None,
    ) -> float:
        cost = self.config.base_cost_per_request
        
        model_multiplier = self.config.model_multipliers.get(model_type, 1.0)
        cost *= model_multiplier
        
        if latency_ms > 0:
            latency_seconds = latency_ms / 1000.0
            cost += latency_seconds * self.config.cost_per_second
        
        if tokens:
            cost += tokens * self.config.cost_per_token
        
        if input_size:
            cost += (input_size / 1024.0) * 0.00001
        
        return round(cost, 6)
    
    def calculate_training_cost(
        self,
        epochs: int,
        batch_size: int,
        dataset_size: int,
        duration_hours: float,
    ) -> float:
        base_cost = 0.1
        epoch_cost = epochs * 0.01
        batch_cost = (dataset_size / batch_size) * 0.001
        time_cost = duration_hours * 0.5
        
        return round(base_cost + epoch_cost + batch_cost + time_cost, 6)


def calculate_request_cost(
    model_type: str = "generic",
    latency_ms: float = 0,
    tokens: Optional[int] = None,
    input_size: Optional[int] = None,
) -> float:
    calculator = BillingCalculator()
    return calculator.calculate_request_cost(model_type, latency_ms, tokens, input_size)

