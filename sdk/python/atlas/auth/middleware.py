"""Authentication Middleware"""

from typing import Optional, Callable, Any
from aiohttp import web
from functools import wraps

from .api_key import validate_api_key
from .jwt import verify_token
from .rbac import Permission, get_rbac


def get_auth_header(request: web.Request) -> Optional[str]:
    auth_header = request.headers.get("Authorization", "")
    if auth_header.startswith("Bearer "):
        return auth_header[7:]
    elif auth_header.startswith("ApiKey "):
        return auth_header[7:]
    return None


def require_auth(permission: Optional[Permission] = None):
    def decorator(func: Callable) -> Callable:
        @wraps(func)
        async def wrapper(request: web.Request, *args, **kwargs) -> web.Response:
            auth_token = get_auth_header(request)
            
            if not auth_token:
                return web.json_response(
                    {"error": "Authentication required"},
                    status=401
                )
            
            user_id = None
            permissions = []
            
            jwt_payload = verify_token(auth_token)
            if jwt_payload:
                user_id = jwt_payload.get("user_id")
                permissions = jwt_payload.get("permissions", [])
            else:
                api_key = validate_api_key(auth_token)
                if api_key:
                    user_id = api_key.user_id
                    permissions = list(api_key.permissions)
            
            if not user_id:
                return web.json_response(
                    {"error": "Invalid authentication token"},
                    status=401
                )
            
            if permission:
                rbac = get_rbac()
                has_permission = False
                
                if permission.value in permissions:
                    has_permission = True
                elif rbac.has_permission(user_id, permission):
                    has_permission = True
                
                if not has_permission:
                    return web.json_response(
                        {"error": f"Permission denied: {permission.value} required"},
                        status=403
                    )
            
            request["user_id"] = user_id
            request["permissions"] = permissions
            
            return await func(request, *args, **kwargs)
        
        return wrapper
    return decorator


def optional_auth(func: Callable) -> Callable:
    @wraps(func)
    async def wrapper(request: web.Request, *args, **kwargs) -> web.Response:
        auth_token = get_auth_header(request)
        
        if auth_token:
            user_id = None
            permissions = []
            
            jwt_payload = verify_token(auth_token)
            if jwt_payload:
                user_id = jwt_payload.get("user_id")
                permissions = jwt_payload.get("permissions", [])
            else:
                api_key = validate_api_key(auth_token)
                if api_key:
                    user_id = api_key.user_id
                    permissions = list(api_key.permissions)
            
            if user_id:
                request["user_id"] = user_id
                request["permissions"] = permissions
        
        return await func(request, *args, **kwargs)
    
    return wrapper

