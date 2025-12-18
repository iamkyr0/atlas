"""Retry Logic with Exponential Backoff"""

import asyncio
import time
from typing import Callable, Type, Tuple, Optional
from functools import wraps
from dataclasses import dataclass


@dataclass
class RetryConfig:
    max_attempts: int = 3
    initial_delay: float = 1.0
    max_delay: float = 60.0
    exponential_base: float = 2.0
    jitter: bool = True
    retryable_exceptions: Tuple[Type[Exception], ...] = (Exception,)


def exponential_backoff(
    attempt: int,
    initial_delay: float = 1.0,
    max_delay: float = 60.0,
    exponential_base: float = 2.0,
    jitter: bool = True,
) -> float:
    delay = min(initial_delay * (exponential_base ** attempt), max_delay)
    if jitter:
        import random
        delay = delay * (0.5 + random.random() * 0.5)
    return delay


def retry(config: Optional[RetryConfig] = None):
    if config is None:
        config = RetryConfig()
    
    def decorator(func: Callable) -> Callable:
        if asyncio.iscoroutinefunction(func):
            @wraps(func)
            async def async_wrapper(*args, **kwargs):
                last_exception = None
                
                for attempt in range(config.max_attempts):
                    try:
                        return await func(*args, **kwargs)
                    except config.retryable_exceptions as e:
                        last_exception = e
                        
                        if attempt < config.max_attempts - 1:
                            delay = exponential_backoff(
                                attempt,
                                config.initial_delay,
                                config.max_delay,
                                config.exponential_base,
                                config.jitter,
                            )
                            await asyncio.sleep(delay)
                        else:
                            raise
                
                if last_exception:
                    raise last_exception
            
            return async_wrapper
        else:
            @wraps(func)
            def sync_wrapper(*args, **kwargs):
                last_exception = None
                
                for attempt in range(config.max_attempts):
                    try:
                        return func(*args, **kwargs)
                    except config.retryable_exceptions as e:
                        last_exception = e
                        
                        if attempt < config.max_attempts - 1:
                            delay = exponential_backoff(
                                attempt,
                                config.initial_delay,
                                config.max_delay,
                                config.exponential_base,
                                config.jitter,
                            )
                            time.sleep(delay)
                        else:
                            raise
                
                if last_exception:
                    raise last_exception
            
            return sync_wrapper
    
    return decorator

