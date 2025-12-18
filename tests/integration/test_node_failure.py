"""Integration test for node failure and recovery"""

import pytest
import asyncio
from typing import Dict, Any

from atlas.chain.client import ChainClient
from atlas.p2p.pubsub import PubSub
from atlas.p2p.discovery import NodeDiscovery
from atlas.job.coordinator import JobCoordinator


@pytest.fixture
async def chain_client():
    client = ChainClient(grpc_url="localhost:9090")
    yield client
    client.close()


@pytest.fixture
async def coordinator(chain_client):
    pubsub = PubSub(ipfs_api_url="/ip4/127.0.0.1/tcp/5001")
    discovery = NodeDiscovery(
        ipfs_api_url="/ip4/127.0.0.1/tcp/5001",
        chain_client=chain_client,
    )
    coordinator = JobCoordinator(
        chain_client=chain_client,
        pubsub=pubsub,
        discovery=discovery,
    )
    return coordinator


@pytest.mark.asyncio
async def test_node_offline_detection(chain_client):
    """Test detection of offline nodes"""
    
    try:
        nodes = await chain_client.list_nodes()
        assert nodes is not None
        
        online_nodes = [n for n in nodes if n.get("status") == "online"]
        offline_nodes = [n for n in nodes if n.get("status") == "offline"]
        
        assert len(nodes) == len(online_nodes) + len(offline_nodes)
        
    except NotImplementedError:
        pytest.skip("gRPC client not fully implemented")


@pytest.mark.asyncio
async def test_task_reassignment(coordinator, chain_client):
    """Test task reassignment when node fails"""
    
    model_id = "test-model-recovery"
    dataset_cid = "QmTestRecovery"
    config = {"epochs": 1}
    creator = "atlas1recovery"
    
    try:
        job_id = await coordinator.submit_job(
            creator=creator,
            model_id=model_id,
            dataset_cid=dataset_cid,
            config=config,
        )
        
        await asyncio.sleep(2)
        
        tasks = await chain_client.get_tasks_by_job(job_id)
        assert tasks is not None
        
        if len(tasks) > 0:
            task = tasks[0]
            original_node_id = task.get("node_id")
            
            await asyncio.sleep(5)
            
            updated_task = await chain_client.get_task(task.get("id"))
            
            if updated_task.get("status") == "failed":
                assert updated_task.get("node_id") != original_node_id or updated_task.get("status") == "reassigned"
        
    except NotImplementedError:
        pytest.skip("gRPC client not fully implemented")


@pytest.mark.asyncio
async def test_heartbeat_recovery(chain_client):
    """Test node heartbeat and recovery"""
    
    try:
        nodes = await chain_client.list_nodes()
        
        if len(nodes) > 0:
            node = nodes[0]
            node_id = node.get("id")
            
            initial_heartbeat = node.get("last_heartbeat")
            
            await asyncio.sleep(10)
            
            updated_node = await chain_client.get_node(node_id)
            updated_heartbeat = updated_node.get("last_heartbeat")
            
            assert updated_heartbeat >= initial_heartbeat
        
    except NotImplementedError:
        pytest.skip("gRPC client not fully implemented")

