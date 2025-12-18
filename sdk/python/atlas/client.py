"""
Atlas Client - Main SDK client for interacting with Atlas platform (P2P)
"""

import asyncio
from typing import Optional, List, Dict, Any
from ipfshttpclient import Client as IPFSClient

from .models import Job, Model, Task
from .chain.client import ChainClient
from .p2p.pubsub import PubSub
from .p2p.discovery import NodeDiscovery
from .job.coordinator import JobCoordinator


class AtlasClient:
    """Main client for interacting with Atlas platform (P2P, no API Gateway)"""
    
    def __init__(
        self,
        ipfs_api_url: str = "/ip4/127.0.0.1/tcp/5001",
        chain_grpc_url: str = "localhost:9090",
        creator: str = "",
    ):
        """
        Initialize Atlas client (P2P mode)
        
        Args:
            ipfs_api_url: IPFS API URL
            chain_grpc_url: Chain gRPC URL (host:port)
            creator: Creator address for transactions
        """
        self.ipfs_client = IPFSClient(ipfs_api_url)
        self.chain_client = ChainClient(grpc_url=chain_grpc_url)
        self.pubsub = PubSub(ipfs_api_url=ipfs_api_url)
        self.discovery = NodeDiscovery(
            ipfs_api_url=ipfs_api_url,
            chain_client=self.chain_client,
        )
        self.coordinator = JobCoordinator(
            chain_client=self.chain_client,
            pubsub=self.pubsub,
            discovery=self.discovery,
        )
        self.creator = creator
    
    async def __aenter__(self):
        return self
    
    async def __aexit__(self, exc_type, exc_val, exc_tb):
        self.chain_client.close()
    
    async def submit_job(
        self,
        model_id: str,
        dataset_cid: str,
        config: Dict[str, Any],
    ) -> str:
        """
        Submit a training job (P2P mode)
        
        Args:
            model_id: Model identifier
            dataset_cid: IPFS CID of the dataset
            config: Training configuration
            
        Returns:
            Job ID
        """
        job_id = await self.coordinator.submit_job(
            creator=self.creator,
            model_id=model_id,
            dataset_cid=dataset_cid,
            config=config,
        )
        return job_id
    
    async def get_job(self, job_id: str) -> Dict[str, Any]:
        """Get job status (P2P mode)"""
        return await self.chain_client.get_job(job_id)
    
    async def list_jobs(self) -> List[Dict[str, Any]]:
        """List all jobs (P2P mode)"""
        return await self.chain_client.list_jobs()
    
    async def upload_dataset(self, file_path: str) -> str:
        """
        Upload dataset to IPFS
        
        Args:
            file_path: Path to dataset file
            
        Returns:
            IPFS CID
        """
        result = self.ipfs_client.add(file_path)
        return result["Hash"]
    
    async def download_model(self, model_cid: str, output_path: str):
        """
        Download model from IPFS
        
        Args:
            model_cid: IPFS CID of the model
            output_path: Path to save the model
        """
        self.ipfs_client.get(model_cid, output_path)
    
    async def subscribe_to_job_updates(self, job_id: str):
        """
        Subscribe to job updates via IPFS pub/sub
        
        Args:
            job_id: Job ID to subscribe to
            
        Yields:
            Update events
        """
        from .job.monitor import JobMonitor
        monitor = JobMonitor(self.pubsub, self.chain_client)
        async for update in monitor.monitor(job_id):
            yield update

