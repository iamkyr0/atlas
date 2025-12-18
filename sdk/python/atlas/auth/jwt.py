"""JWT Token Management"""

import jwt
import time
from typing import Optional, Dict, Any
from datetime import datetime, timedelta


class JWTManager:
    def __init__(self, secret_key: str, algorithm: str = "HS256"):
        self.secret_key = secret_key
        self.algorithm = algorithm
    
    def create_token(
        self,
        user_id: str,
        permissions: list[str] = None,
        expires_in_seconds: int = 3600,
        additional_claims: Dict[str, Any] = None,
    ) -> str:
        now = datetime.utcnow()
        exp = now + timedelta(seconds=expires_in_seconds)
        
        payload = {
            "user_id": user_id,
            "permissions": permissions or [],
            "iat": int(now.timestamp()),
            "exp": int(exp.timestamp()),
        }
        
        if additional_claims:
            payload.update(additional_claims)
        
        token = jwt.encode(payload, self.secret_key, algorithm=self.algorithm)
        return token
    
    def verify_token(self, token: str) -> Optional[Dict[str, Any]]:
        try:
            payload = jwt.decode(
                token, self.secret_key, algorithms=[self.algorithm]
            )
            return payload
        except jwt.ExpiredSignatureError:
            return None
        except jwt.InvalidTokenError:
            return None
    
    def refresh_token(self, token: str, expires_in_seconds: int = 3600) -> Optional[str]:
        payload = self.verify_token(token)
        if not payload:
            return None
        
        user_id = payload.get("user_id")
        permissions = payload.get("permissions", [])
        
        return self.create_token(
            user_id, permissions, expires_in_seconds
        )


_global_jwt_manager: Optional[JWTManager] = None


def initialize_jwt(secret_key: str, algorithm: str = "HS256"):
    global _global_jwt_manager
    _global_jwt_manager = JWTManager(secret_key, algorithm)


def create_token(
    user_id: str,
    permissions: list[str] = None,
    expires_in_seconds: int = 3600,
) -> str:
    if _global_jwt_manager is None:
        raise RuntimeError("JWT not initialized. Call initialize_jwt() first.")
    return _global_jwt_manager.create_token(user_id, permissions, expires_in_seconds)


def verify_token(token: str) -> Optional[Dict[str, Any]]:
    if _global_jwt_manager is None:
        raise RuntimeError("JWT not initialized. Call initialize_jwt() first.")
    return _global_jwt_manager.verify_token(token)

