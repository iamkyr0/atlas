"""
Atlas SDK Models - Data models for Atlas platform
"""

from typing import Optional, Dict, Any, List
from enum import Enum
from datetime import datetime
from pydantic import BaseModel, Field


class TaskStatus(str, Enum):
    """Task status enumeration"""
    PENDING = "pending"
    ASSIGNED = "assigned"
    IN_PROGRESS = "in_progress"
    PAUSED = "paused"
    ROLLBACK = "rollback"
    DELEGATED = "delegated"
    COMPLETED = "completed"
    FAILED = "failed"


class NodeStatus(str, Enum):
    """Node status enumeration"""
    ONLINE = "online"
    DEGRADED = "degraded"
    OFFLINE = "offline"
    RECOVERING = "recovering"


class Job(BaseModel):
    """Training job model"""
    id: str
    model_id: str
    dataset_cid: str
    config: Dict[str, Any]
    status: TaskStatus
    created_at: datetime
    updated_at: datetime
    progress: float = Field(default=0.0, ge=0.0, le=1.0)
    tasks: List["Task"] = Field(default_factory=list)


class Task(BaseModel):
    """Task model"""
    id: str
    job_id: str
    shard_id: Optional[str] = None
    node_id: Optional[str] = None
    status: TaskStatus
    created_at: datetime
    updated_at: datetime
    progress: float = Field(default=0.0, ge=0.0, le=1.0)
    checkpoint_cid: Optional[str] = None


class Model(BaseModel):
    """Model registry entry"""
    id: str
    name: str
    version: str
    cid: str
    created_at: datetime
    metadata: Dict[str, Any] = Field(default_factory=dict)


class Shard(BaseModel):
    """Shard model"""
    id: str
    job_id: str
    cid: str
    size: int
    hash: str
    node_id: Optional[str] = None
    status: TaskStatus
    created_at: datetime


class Node(BaseModel):
    """Node model"""
    id: str
    status: NodeStatus
    uptime_percentage: float = Field(default=0.0, ge=0.0, le=100.0)
    reputation: float = Field(default=0.0, ge=0.0, le=100.0)
    resources: Dict[str, Any] = Field(default_factory=dict)
    last_heartbeat: Optional[datetime] = None

