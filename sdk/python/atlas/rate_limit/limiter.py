"""Rate Limiting Implementations"""

import time
from typing import Dict, Optional
from collections import defaultdict
from dataclasses import dataclass


@dataclass
class RateLimitConfig:
    requests_per_second: float
    burst_size: Optional[int] = None
    window_size: Optional[int] = None


class TokenBucket:
    def __init__(self, capacity: int, refill_rate: float):
        self.capacity = capacity
        self.refill_rate = refill_rate
        self.tokens = float(capacity)
        self.last_refill = time.time()
    
    def consume(self, tokens: int = 1) -> bool:
        self._refill()
        if self.tokens >= tokens:
            self.tokens -= tokens
            return True
        return False
    
    def _refill(self):
        now = time.time()
        elapsed = now - self.last_refill
        tokens_to_add = elapsed * self.refill_rate
        self.tokens = min(self.capacity, self.tokens + tokens_to_add)
        self.last_refill = now
    
    def available(self) -> float:
        self._refill()
        return self.tokens


class SlidingWindow:
    def __init__(self, max_requests: int, window_seconds: int):
        self.max_requests = max_requests
        self.window_seconds = window_seconds
        self.requests: Dict[str, list] = defaultdict(list)
    
    def is_allowed(self, key: str) -> bool:
        now = time.time()
        window_start = now - self.window_seconds
        
        if key not in self.requests:
            self.requests[key] = []
        
        requests = self.requests[key]
        requests[:] = [req_time for req_time in requests if req_time > window_start]
        
        if len(requests) < self.max_requests:
            requests.append(now)
            return True
        
        return False
    
    def remaining(self, key: str) -> int:
        now = time.time()
        window_start = now - self.window_seconds
        
        if key not in self.requests:
            return self.max_requests
        
        requests = self.requests[key]
        requests[:] = [req_time for req_time in requests if req_time > window_start]
        
        return max(0, self.max_requests - len(requests))


class RateLimiter:
    def __init__(self, config: RateLimitConfig):
        self.config = config
        
        if config.burst_size:
            self.limiter = TokenBucket(
                capacity=config.burst_size,
                refill_rate=config.requests_per_second,
            )
        elif config.window_size:
            self.limiter = SlidingWindow(
                max_requests=int(config.requests_per_second * config.window_size),
                window_seconds=config.window_size,
            )
        else:
            self.limiter = SlidingWindow(
                max_requests=int(config.requests_per_second),
                window_seconds=1,
            )
    
    def is_allowed(self, key: str = "default") -> bool:
        if isinstance(self.limiter, TokenBucket):
            return self.limiter.consume()
        elif isinstance(self.limiter, SlidingWindow):
            return self.limiter.is_allowed(key)
        return True
    
    def remaining(self, key: str = "default") -> int:
        if isinstance(self.limiter, TokenBucket):
            return int(self.limiter.available())
        elif isinstance(self.limiter, SlidingWindow):
            return self.limiter.remaining(key)
        return 0


_global_limiters: Dict[str, RateLimiter] = {}


def get_rate_limiter(name: str, config: RateLimitConfig) -> RateLimiter:
    if name not in _global_limiters:
        _global_limiters[name] = RateLimiter(config)
    return _global_limiters[name]

