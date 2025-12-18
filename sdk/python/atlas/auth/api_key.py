"""API Key Management"""

import hashlib
import secrets
import time
from typing import Optional, Dict, Set
from dataclasses import dataclass
from datetime import datetime, timedelta


@dataclass
class APIKey:
    key: str
    key_hash: str
    user_id: str
    permissions: Set[str]
    created_at: datetime
    expires_at: Optional[datetime]
    last_used: Optional[datetime]
    rate_limit: Optional[int]
    rate_limit_window: Optional[int]


class APIKeyManager:
    def __init__(self):
        self.keys: Dict[str, APIKey] = {}
        self.key_to_hash: Dict[str, str] = {}
        self.rate_limit_store: Dict[str, Dict[str, int]] = {}
    
    def generate_key(
        self,
        user_id: str,
        permissions: Set[str] = None,
        expires_in_days: Optional[int] = None,
        rate_limit: Optional[int] = None,
        rate_limit_window: Optional[int] = None,
    ) -> str:
        key = f"atlas_{secrets.token_urlsafe(32)}"
        key_hash = hashlib.sha256(key.encode()).hexdigest()
        
        expires_at = None
        if expires_in_days:
            expires_at = datetime.now() + timedelta(days=expires_in_days)
        
        api_key = APIKey(
            key=key,
            key_hash=key_hash,
            user_id=user_id,
            permissions=permissions or set(),
            created_at=datetime.now(),
            expires_at=expires_at,
            last_used=None,
            rate_limit=rate_limit,
            rate_limit_window=rate_limit_window or 3600,
        )
        
        self.keys[key_hash] = api_key
        self.key_to_hash[key] = key_hash
        
        return key
    
    def validate_key(self, key: str) -> Optional[APIKey]:
        key_hash = hashlib.sha256(key.encode()).hexdigest()
        
        api_key = self.keys.get(key_hash)
        if not api_key:
            return None
        
        if api_key.expires_at and datetime.now() > api_key.expires_at:
            return None
        
        api_key.last_used = datetime.now()
        
        if api_key.rate_limit:
            if not self._check_rate_limit(key_hash, api_key):
                return None
        
        return api_key
    
    def _check_rate_limit(self, key_hash: str, api_key: APIKey) -> bool:
        now = int(time.time())
        window_start = now - api_key.rate_limit_window
        
        if key_hash not in self.rate_limit_store:
            self.rate_limit_store[key_hash] = {}
        
        store = self.rate_limit_store[key_hash]
        
        store = {k: v for k, v in store.items() if int(k) > window_start}
        self.rate_limit_store[key_hash] = store
        
        request_count = sum(store.values())
        
        if request_count >= api_key.rate_limit:
            return False
        
        store[str(now)] = store.get(str(now), 0) + 1
        return True
    
    def revoke_key(self, key: str) -> bool:
        key_hash = hashlib.sha256(key.encode()).hexdigest()
        if key_hash in self.keys:
            del self.keys[key_hash]
            if key in self.key_to_hash:
                del self.key_to_hash[key]
            if key_hash in self.rate_limit_store:
                del self.rate_limit_store[key_hash]
            return True
        return False
    
    def get_key_info(self, key: str) -> Optional[APIKey]:
        key_hash = hashlib.sha256(key.encode()).hexdigest()
        return self.keys.get(key_hash)
    
    def list_keys_for_user(self, user_id: str) -> list[APIKey]:
        return [key for key in self.keys.values() if key.user_id == user_id]


_global_manager = APIKeyManager()


def validate_api_key(key: str) -> Optional[APIKey]:
    return _global_manager.validate_key(key)


def get_api_key_manager() -> APIKeyManager:
    return _global_manager

