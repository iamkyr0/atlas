"""
Job Coordinator - Coordinates job submission and task assignment
"""

import asyncio
from typing import Dict, Any, List, Optional
from ..chain.client import ChainClient
from ..p2p.pubsub import PubSub
from ..p2p.discovery import NodeDiscovery


class JobCoordinator:
    """Coordinates job submission and task assignment"""
    
    def __init__(
        self,
        chain_client: ChainClient,
        pubsub: PubSub,
        discovery: NodeDiscovery,
    ):
        """
        Initialize job coordinator
        
        Args:
            chain_client: Blockchain client
            pubsub: Pub/Sub instance
            discovery: Node discovery instance
        """
        self.chain_client = chain_client
        self.pubsub = pubsub
        self.discovery = discovery
    
    async def submit_job(
        self,
        creator: str,
        model_id: str,
        dataset_cid: str,
        config: Dict[str, Any],
    ) -> str:
        """
        Submit a training job
        
        Args:
            creator: Creator address
            model_id: Model ID
            dataset_cid: Dataset IPFS CID
            config: Training configuration
            
        Returns:
            Job ID
        """
        # Submit to blockchain
        job_id = await self.chain_client.submit_job(
            creator=creator,
            model_id=model_id,
            dataset_cid=dataset_cid,
            config=config,
        )
        
        # Publish job announcement via IPFS pub/sub
        await self.pubsub.publish(
            topic="/atlas/jobs/new",
            message={
                "job_id": job_id,
                "model_id": model_id,
                "dataset_cid": dataset_cid,
                "config": config,
            }
        )
        
        return job_id
    
    async def monitor_job(self, job_id: str) -> AsyncIterator[Dict[str, Any]]:
        """
        Monitor job updates via IPFS pub/sub with blockchain fallback
        
        Args:
            job_id: Job ID to monitor
            
        Yields:
            Job update events
        """
        # Subscribe to job updates via pub/sub
        topic = f"/atlas/jobs/{job_id}/updates"
        
        # Use pub/sub as primary
        try:
            async for update in self.pubsub.subscribe(topic):
                yield update
        except Exception:
            # Fallback to blockchain polling
            await self._poll_blockchain(job_id)
    
    async def _poll_blockchain(self, job_id: str) -> AsyncIterator[Dict[str, Any]]:
        """
        Poll blockchain for job updates (fallback)
        
        Args:
            job_id: Job ID
            
        Yields:
            Job updates
        """
        last_status = None
        
        while True:
            try:
                job = await self.chain_client.get_job(job_id)
                current_status = job.get("status")
                
                if current_status != last_status:
                    yield {
                        "job_id": job_id,
                        "status": current_status,
                        "progress": job.get("progress", 0.0),
                    }
                    last_status = current_status
                
                await asyncio.sleep(5)  # Poll every 5 seconds
            except Exception as e:
                print(f"Error polling blockchain: {e}")
                await asyncio.sleep(10)  # Wait longer on error

