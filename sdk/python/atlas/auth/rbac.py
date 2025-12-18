"""Role-Based Access Control (RBAC)"""

from enum import Enum
from typing import Set, Dict, Optional
from dataclasses import dataclass, field


class Permission(Enum):
    SUBMIT_JOB = "submit_job"
    VIEW_JOB = "view_job"
    LIST_JOBS = "list_jobs"
    REGISTER_MODEL = "register_model"
    VIEW_MODEL = "view_model"
    LIST_MODELS = "list_models"
    VIEW_NODE = "view_node"
    LIST_NODES = "list_nodes"
    INFERENCE = "inference"
    ADMIN = "admin"


@dataclass
class Role:
    name: str
    permissions: Set[Permission] = field(default_factory=set)
    description: str = ""


class RBAC:
    def __init__(self):
        self.roles: Dict[str, Role] = {}
        self.user_roles: Dict[str, Set[str]] = {}
        self._initialize_default_roles()
    
    def _initialize_default_roles(self):
        admin_role = Role(
            name="admin",
            permissions={p for p in Permission},
            description="Full access to all resources",
        )
        self.roles["admin"] = admin_role
        
        user_role = Role(
            name="user",
            permissions={
                Permission.SUBMIT_JOB,
                Permission.VIEW_JOB,
                Permission.LIST_JOBS,
                Permission.REGISTER_MODEL,
                Permission.VIEW_MODEL,
                Permission.LIST_MODELS,
                Permission.INFERENCE,
            },
            description="Standard user permissions",
        )
        self.roles["user"] = user_role
        
        viewer_role = Role(
            name="viewer",
            permissions={
                Permission.VIEW_JOB,
                Permission.LIST_JOBS,
                Permission.VIEW_MODEL,
                Permission.LIST_MODELS,
                Permission.VIEW_NODE,
                Permission.LIST_NODES,
            },
            description="Read-only access",
        )
        self.roles["viewer"] = viewer_role
    
    def create_role(self, name: str, permissions: Set[Permission], description: str = "") -> Role:
        role = Role(name=name, permissions=permissions, description=description)
        self.roles[name] = role
        return role
    
    def assign_role(self, user_id: str, role_name: str) -> bool:
        if role_name not in self.roles:
            return False
        
        if user_id not in self.user_roles:
            self.user_roles[user_id] = set()
        
        self.user_roles[user_id].add(role_name)
        return True
    
    def revoke_role(self, user_id: str, role_name: str) -> bool:
        if user_id not in self.user_roles:
            return False
        
        self.user_roles[user_id].discard(role_name)
        if not self.user_roles[user_id]:
            del self.user_roles[user_id]
        return True
    
    def has_permission(self, user_id: str, permission: Permission) -> bool:
        user_roles = self.user_roles.get(user_id, set())
        
        for role_name in user_roles:
            role = self.roles.get(role_name)
            if role and permission in role.permissions:
                return True
        
        return False
    
    def get_user_permissions(self, user_id: str) -> Set[Permission]:
        user_roles = self.user_roles.get(user_id, set())
        permissions = set()
        
        for role_name in user_roles:
            role = self.roles.get(role_name)
            if role:
                permissions.update(role.permissions)
        
        return permissions
    
    def check_permission(self, user_id: str, permission: Permission) -> bool:
        return self.has_permission(user_id, permission)


_global_rbac = RBAC()


def get_rbac() -> RBAC:
    return _global_rbac

