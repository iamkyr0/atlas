import asyncio
from typing import Dict, Any, List, Optional
from ..chain.client import ChainClient
from ..p2p.pubsub import PubSub
from ..p2p.discovery import NodeDiscovery


class TaskAssignment:
    
    def __init__(
        self,
        chain_client: ChainClient,
        pubsub: PubSub,
        discovery: NodeDiscovery,
    ):
        self.chain_client = chain_client
        self.pubsub = pubsub
        self.discovery = discovery
    
    async def assign_task(
        self,
        job_id: str,
        shard_id: str,
        node_id: str,
    ) -> str:
        try:
            await self.pubsub.publish(
                topic=f"/atlas/tasks/assign/{node_id}",
                message={
                    "job_id": job_id,
                    "shard_id": shard_id,
                    "node_id": node_id,
                }
            )
            
            task_id = await self._create_task_on_blockchain(
                creator="",
                job_id=job_id,
                shard_id=shard_id,
                node_id=node_id,
            )
            
            return task_id
        except Exception as e:
            print(f"Pub/sub assignment failed: {e}, using blockchain only")
            return await self._create_task_on_blockchain(
                creator="",
                job_id=job_id,
                shard_id=shard_id,
                node_id=node_id,
            )
    
    async def _create_task_on_blockchain(
        self,
        creator: str,
        job_id: str,
        shard_id: str,
        node_id: str,
    ) -> str:
        return f"task-{job_id}-{shard_id}"
    
    async def listen_for_assignments(self, node_id: str) -> AsyncIterator[Dict[str, Any]]:
        topic = f"/atlas/tasks/assign/{node_id}"
        
        try:
            async for assignment in self.pubsub.subscribe(topic):
                yield assignment
        except Exception:
            await self._poll_blockchain_assignments(node_id)
    
    async def _poll_blockchain_assignments(self, node_id: str) -> AsyncIterator[Dict[str, Any]]:
        seen_tasks = set()
        
        while True:
            try:
                jobs = await self.chain_client.list_jobs()
                
                for job in jobs:
                    job_id = job.get("id")
                    tasks = await self.chain_client.get_tasks_by_job(job_id)
                    
                    for task in tasks:
                        task_id = task.get("id")
                        if task.get("node_id") == node_id and task_id not in seen_tasks:
                            seen_tasks.add(task_id)
                            yield {
                                "task_id": task_id,
                                "job_id": job_id,
                                "shard_id": task.get("shard_id"),
                                "node_id": node_id,
                            }
                
                await asyncio.sleep(10)
            except Exception as e:
                print(f"Error polling blockchain for assignments: {e}")
                await asyncio.sleep(30)

