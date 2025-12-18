"""
Atlas Blockchain Types - Type definitions for blockchain messages
"""

from typing import Dict, Any, Optional
from dataclasses import dataclass
from datetime import datetime


@dataclass
class JobMessage:
    """Job message for blockchain"""
    creator: str
    model_id: str
    dataset_cid: str
    config: Dict[str, Any]


@dataclass
class TaskMessage:
    """Task message for blockchain"""
    creator: str
    job_id: str
    shard_id: str
    node_id: str


@dataclass
class ModelMessage:
    """Model registration message"""
    creator: str
    name: str
    version: str
    cid: str
    metadata: Dict[str, str]


@dataclass
class NodeMessage:
    """Node registration message"""
    creator: str
    node_id: str
    address: str
    cpu_cores: int
    gpu_count: int
    memory_gb: int
    storage_gb: int

