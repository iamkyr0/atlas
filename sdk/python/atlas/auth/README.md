# Atlas Authentication & Authorization

Authentication and authorization module for Atlas platform.

## Features

- **API Key Management**: Generate, validate, and manage API keys with rate limiting
- **JWT Tokens**: Create and verify JWT tokens for stateless authentication
- **RBAC**: Role-based access control with predefined and custom roles
- **Middleware**: Easy-to-use decorators for protecting endpoints

## API Key Management

```python
from atlas.auth.api_key import APIKeyManager

manager = APIKeyManager()

# Generate API key
key = manager.generate_key(
    user_id="user-123",
    permissions={"submit_job", "view_job"},
    expires_in_days=30,
    rate_limit=100,
    rate_limit_window=3600,
)

# Validate API key
api_key = manager.validate_key(key)
if api_key:
    print(f"User: {api_key.user_id}")
    print(f"Permissions: {api_key.permissions}")
```

## JWT Tokens

```python
from atlas.auth.jwt import initialize_jwt, create_token, verify_token

# Initialize JWT (once at startup)
initialize_jwt("your-secret-key")

# Create token
token = create_token(
    user_id="user-123",
    permissions=["submit_job", "view_job"],
    expires_in_seconds=3600,
)

# Verify token
payload = verify_token(token)
if payload:
    print(f"User: {payload['user_id']}")
    print(f"Permissions: {payload['permissions']}")
```

## RBAC

```python
from atlas.auth.rbac import get_rbac, Permission

rbac = get_rbac()

# Assign role
rbac.assign_role("user-123", "user")

# Check permission
if rbac.has_permission("user-123", Permission.SUBMIT_JOB):
    print("User can submit jobs")

# Create custom role
custom_role = rbac.create_role(
    name="custom",
    permissions={Permission.VIEW_JOB, Permission.LIST_JOBS},
    description="Custom role for testing",
)
rbac.assign_role("user-456", "custom")
```

## Middleware

```python
from aiohttp import web
from atlas.auth.middleware import require_auth, optional_auth
from atlas.auth.rbac import Permission

# Require authentication and permission
@require_auth(permission=Permission.SUBMIT_JOB)
async def submit_job_handler(request: web.Request):
    user_id = request["user_id"]
    permissions = request["permissions"]
    # Handle request...

# Optional authentication
@optional_auth
async def public_handler(request: web.Request):
    user_id = request.get("user_id")  # May be None
    # Handle request...
```

## Default Roles

- **admin**: Full access to all resources
- **user**: Standard user permissions (submit jobs, view models, inference)
- **viewer**: Read-only access

## Permissions

- `SUBMIT_JOB`: Submit training jobs
- `VIEW_JOB`: View job details
- `LIST_JOBS`: List all jobs
- `REGISTER_MODEL`: Register new models
- `VIEW_MODEL`: View model details
- `LIST_MODELS`: List all models
- `VIEW_NODE`: View node details
- `LIST_NODES`: List all nodes
- `INFERENCE`: Run model inference
- `ADMIN`: Full administrative access

## Usage in Requests

### API Key
```
Authorization: ApiKey atlas_xxxxxxxxxxxxx
```

### JWT Token
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

