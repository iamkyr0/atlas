"""Rate Limiting Middleware"""

from typing import Callable, Optional
from aiohttp import web
from functools import wraps

from .limiter import RateLimiter, get_rate_limiter, RateLimitConfig
from ..auth.middleware import get_auth_header


def rate_limit_middleware(
    requests_per_second: float,
    burst_size: Optional[int] = None,
    window_size: Optional[int] = None,
    key_func: Optional[Callable] = None,
):
    config = RateLimitConfig(
        requests_per_second=requests_per_second,
        burst_size=burst_size,
        window_size=window_size,
    )
    limiter = RateLimiter(config)
    
    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(request: web.Request, *args, **kwargs) -> web.Response:
            if key_func:
                key = key_func(request)
            else:
                auth_token = get_auth_header(request)
                key = auth_token or request.remote
            
            if not limiter.is_allowed(key):
                return web.json_response(
                    {
                        "error": "Rate limit exceeded",
                        "retry_after": limiter.remaining(key),
                    },
                    status=429,
                    headers={
                        "X-RateLimit-Limit": str(int(config.requests_per_second)),
                        "X-RateLimit-Remaining": str(limiter.remaining(key)),
                        "Retry-After": "1",
                    },
                )
            
            response = await func(request, *args, **kwargs)
            
            if isinstance(response, web.Response):
                response.headers["X-RateLimit-Limit"] = str(int(config.requests_per_second))
                response.headers["X-RateLimit-Remaining"] = str(limiter.remaining(key))
            
            return response
        
        return wrapper
    return decorator

