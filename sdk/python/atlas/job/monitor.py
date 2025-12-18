"""
Job Monitor - Monitor jobs via IPFS pub/sub
"""

import asyncio
from typing import Dict, Any, AsyncIterator
from ..p2p.pubsub import PubSub
from ..chain.client import ChainClient


class JobMonitor:
    """Monitor job progress via IPFS pub/sub"""
    
    def __init__(self, pubsub: PubSub, chain_client: ChainClient):
        """
        Initialize job monitor
        
        Args:
            pubsub: Pub/Sub instance
            chain_client: Blockchain client for fallback
        """
        self.pubsub = pubsub
        self.chain_client = chain_client
    
    async def monitor(self, job_id: str) -> AsyncIterator[Dict[str, Any]]:
        """
        Monitor job updates
        
        Args:
            job_id: Job ID to monitor
            
        Yields:
            Job update events
        """
        topic = f"/atlas/jobs/{job_id}/updates"
        
        # Try pub/sub first
        try:
            async for update in self.pubsub.subscribe(topic):
                yield update
        except Exception:
            # Fallback to blockchain polling
            async for update in self._poll_blockchain(job_id):
                yield update
    
    async def _poll_blockchain(self, job_id: str) -> AsyncIterator[Dict[str, Any]]:
        """
        Poll blockchain for job updates (fallback)
        
        Args:
            job_id: Job ID
            
        Yields:
            Job updates
        """
        last_progress = None
        
        while True:
            try:
                job = await self.chain_client.get_job(job_id)
                current_progress = job.get("progress", 0.0)
                
                if current_progress != last_progress:
                    yield {
                        "job_id": job_id,
                        "status": job.get("status"),
                        "progress": current_progress,
                        "tasks": job.get("tasks", []),
                    }
                    last_progress = current_progress
                
                # Stop if job is completed or failed
                if job.get("status") in ["completed", "failed"]:
                    break
                
                await asyncio.sleep(5)
            except Exception as e:
                print(f"Error polling blockchain: {e}")
                await asyncio.sleep(10)

