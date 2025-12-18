"""Integration test for authentication and authorization"""

import pytest
from atlas.auth.api_key import APIKeyManager, validate_api_key
from atlas.auth.jwt import JWTManager, create_token, verify_token, initialize_jwt
from atlas.auth.rbac import RBAC, Permission, get_rbac


def test_api_key_generation():
    """Test API key generation and validation"""
    manager = APIKeyManager()
    
    key = manager.generate_key(
        user_id="test-user-1",
        permissions={"submit_job", "view_job"},
        expires_in_days=30,
        rate_limit=100,
        rate_limit_window=3600,
    )
    
    assert key is not None
    assert key.startswith("atlas_")
    
    api_key = manager.validate_key(key)
    assert api_key is not None
    assert api_key.user_id == "test-user-1"
    assert "submit_job" in api_key.permissions


def test_api_key_validation():
    """Test API key validation"""
    manager = APIKeyManager()
    
    key = manager.generate_key(user_id="test-user-2")
    
    valid_key = validate_api_key(key)
    assert valid_key is not None
    
    invalid_key = validate_api_key("invalid_key")
    assert invalid_key is None


def test_api_key_rate_limiting():
    """Test API key rate limiting"""
    manager = APIKeyManager()
    
    key = manager.generate_key(
        user_id="test-user-3",
        rate_limit=5,
        rate_limit_window=60,
    )
    
    for i in range(5):
        api_key = manager.validate_key(key)
        assert api_key is not None
    
    api_key = manager.validate_key(key)
    assert api_key is None


def test_jwt_token_creation():
    """Test JWT token creation and verification"""
    initialize_jwt("test-secret-key")
    
    token = create_token(
        user_id="test-user-4",
        permissions=["submit_job", "view_job"],
        expires_in_seconds=3600,
    )
    
    assert token is not None
    
    payload = verify_token(token)
    assert payload is not None
    assert payload.get("user_id") == "test-user-4"
    assert "submit_job" in payload.get("permissions", [])


def test_jwt_token_expiration():
    """Test JWT token expiration"""
    initialize_jwt("test-secret-key")
    
    token = create_token(
        user_id="test-user-5",
        expires_in_seconds=1,
    )
    
    payload = verify_token(token)
    assert payload is not None
    
    import time
    time.sleep(2)
    
    payload = verify_token(token)
    assert payload is None


def test_rbac_permissions():
    """Test RBAC permission checking"""
    rbac = get_rbac()
    
    rbac.assign_role("test-user-6", "user")
    
    assert rbac.has_permission("test-user-6", Permission.SUBMIT_JOB)
    assert rbac.has_permission("test-user-6", Permission.VIEW_JOB)
    assert not rbac.has_permission("test-user-6", Permission.ADMIN)
    
    rbac.assign_role("test-user-6", "admin")
    assert rbac.has_permission("test-user-6", Permission.ADMIN)


def test_rbac_role_assignment():
    """Test RBAC role assignment"""
    rbac = get_rbac()
    
    assert rbac.assign_role("test-user-7", "viewer")
    assert rbac.has_permission("test-user-7", Permission.VIEW_JOB)
    assert not rbac.has_permission("test-user-7", Permission.SUBMIT_JOB)
    
    rbac.revoke_role("test-user-7", "viewer")
    assert not rbac.has_permission("test-user-7", Permission.VIEW_JOB)


def test_custom_role():
    """Test creating and using custom roles"""
    rbac = get_rbac()
    
    custom_role = rbac.create_role(
        name="custom",
        permissions={Permission.VIEW_JOB, Permission.LIST_JOBS},
        description="Custom role for testing",
    )
    
    assert custom_role.name == "custom"
    assert Permission.VIEW_JOB in custom_role.permissions
    
    rbac.assign_role("test-user-8", "custom")
    assert rbac.has_permission("test-user-8", Permission.VIEW_JOB)
    assert not rbac.has_permission("test-user-8", Permission.SUBMIT_JOB)

