"""
CLI Commands - Individual command implementations
"""

import json
import asyncio
from typing import Dict, Any, Optional
from pathlib import Path

from ..chain.client import ChainClient
from ..p2p.pubsub import PubSub
from ..p2p.discovery import NodeDiscovery
from ..job.coordinator import JobCoordinator
from ..serving.server import ModelServer
from ipfshttpclient import Client as IPFSClient


class SubmitJobCommand:
    """Submit job command"""
    
    def __init__(self, coordinator: JobCoordinator, creator: str):
        self.coordinator = coordinator
        self.creator = creator
    
    async def execute(
        self,
        model_id: str,
        dataset_cid: str,
        config_path: Optional[str] = None,
    ):
        """Execute submit job command"""
        config = {}
        if config_path:
            with open(config_path) as f:
                config = json.load(f)
        
        job_id = await self.coordinator.submit_job(
            creator=self.creator,
            model_id=model_id,
            dataset_cid=dataset_cid,
            config=config,
        )
        
        print(f"Job submitted: {job_id}")
        print(f"Monitor with: atlas get-job {job_id}")


class ListJobsCommand:
    """List jobs command"""
    
    def __init__(self, chain_client: ChainClient):
        self.chain_client = chain_client
    
    async def execute(self):
        """Execute list jobs command"""
        jobs = await self.chain_client.list_jobs()
        
        print(f"Found {len(jobs)} jobs:")
        for job in jobs:
            print(f"  {job.get('id')}: {job.get('status')} - {job.get('progress', 0.0)*100:.1f}%")


class GetJobCommand:
    """Get job command"""
    
    def __init__(self, chain_client: ChainClient):
        self.chain_client = chain_client
    
    async def execute(self, job_id: str):
        """Execute get job command"""
        job = await self.chain_client.get_job(job_id)
        
        print(f"Job ID: {job.get('id')}")
        print(f"Status: {job.get('status')}")
        print(f"Progress: {job.get('progress', 0.0)*100:.1f}%")
        print(f"Model ID: {job.get('model_id')}")
        print(f"Dataset CID: {job.get('dataset_cid')}")
        print(f"Tasks: {len(job.get('tasks', []))}")


class UploadDatasetCommand:
    """Upload dataset command"""
    
    def __init__(self, pubsub: PubSub, ipfs_api_url: str):
        self.pubsub = pubsub
        self.ipfs_client = IPFSClient(ipfs_api_url)
    
    async def execute(self, path: str, encrypt: bool = False):
        """Execute upload dataset command"""
        file_path = Path(path)
        if not file_path.exists():
            print(f"Error: File not found: {path}")
            return
        
        # Encrypt if requested
        if encrypt:
            from ..encryption.encrypt import encrypt_file
            encrypted_path = await encrypt_file(str(file_path))
            file_path = Path(encrypted_path)
        
        # Upload to IPFS
        print(f"Uploading {file_path} to IPFS...")
        result = self.ipfs_client.add(str(file_path))
        cid = result["Hash"]
        
        print(f"Dataset uploaded! CID: {cid}")
        print(f"Use this CID in submit-job command")


class RegisterModelCommand:
    """Register model command"""
    
    def __init__(self, chain_client: ChainClient, pubsub: PubSub, ipfs_api_url: str):
        self.chain_client = chain_client
        self.pubsub = pubsub
        self.ipfs_client = IPFSClient(ipfs_api_url)
    
    async def execute(
        self,
        path: str,
        name: str,
        version: str,
        creator: str,
    ):
        """Execute register model command"""
        file_path = Path(path)
        if not file_path.exists():
            print(f"Error: File not found: {path}")
            return
        
        # Upload to IPFS
        print(f"Uploading model to IPFS...")
        result = self.ipfs_client.add(str(file_path))
        cid = result["Hash"]
        
        # Register on blockchain
        model_id = await self.chain_client.register_model(
            creator=creator,
            name=name,
            version=version,
            cid=cid,
            metadata={},
        )
        
        print(f"Model registered! ID: {model_id}")
        print(f"CID: {cid}")


class DownloadModelCommand:
    """Download model command"""
    
    def __init__(self, chain_client: ChainClient, ipfs_api_url: str):
        self.chain_client = chain_client
        self.ipfs_client = IPFSClient(ipfs_api_url)
    
    async def execute(self, model_id: str, output: Optional[str] = None):
        """Execute download model command"""
        # Get model info from blockchain
        model = await self.chain_client.get_model(model_id)
        cid = model.get("cid")
        
        if not cid:
            print(f"Error: Model {model_id} not found")
            return
        
        # Determine output path
        if not output:
            output = f"./models/{model_id}"
        
        output_path = Path(output)
        output_path.mkdir(parents=True, exist_ok=True)
        
        # Download from IPFS
        print(f"Downloading model {model_id} from IPFS...")
        self.ipfs_client.get(cid, str(output_path))
        
        print(f"Model downloaded to: {output_path}")


class ServeModelCommand:
    """Serve model command"""
    
    def __init__(self, chain_client: ChainClient, ipfs_api_url: str):
        self.chain_client = chain_client
        self.ipfs_api_url = ipfs_api_url
    
    async def execute(self, model_id: str, host: str, port: int):
        """Execute serve model command"""
        server = ModelServer(
            chain_client=self.chain_client,
            ipfs_api_url=self.ipfs_api_url,
        )
        
        print(f"Serving model {model_id} on {host}:{port}")
        await server.serve(model_id=model_id, host=host, port=port)


class DaemonCommand:
    """Daemon command"""
    
    def __init__(
        self,
        chain_client: ChainClient,
        pubsub: PubSub,
        discovery: NodeDiscovery,
        coordinator: JobCoordinator,
    ):
        self.chain_client = chain_client
        self.pubsub = pubsub
        self.discovery = discovery
        self.coordinator = coordinator
    
    async def execute(self, host: str, port: int):
        """Execute daemon command"""
        from .daemon import DaemonServer
        
        server = DaemonServer(
            chain_client=self.chain_client,
            pubsub=self.pubsub,
            discovery=self.discovery,
            coordinator=self.coordinator,
        )
        
        print(f"Starting Atlas daemon on {host}:{port}")
        await server.start(host=host, port=port)

