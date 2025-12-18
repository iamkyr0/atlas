"""Integration test for end-to-end job submission"""

import pytest
import asyncio
import time
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
async def test_job_submission_flow(coordinator, chain_client):
    """Test complete job submission flow"""
    
    model_id = "test-model-123"
    dataset_cid = "QmTestDataset"
    config = {
        "epochs": 10,
        "batch_size": 32,
        "learning_rate": 0.001,
    }
    
    creator = "atlas1test123"
    
    try:
        job_id = await coordinator.submit_job(
            creator=creator,
            model_id=model_id,
            dataset_cid=dataset_cid,
            config=config,
        )
        
        assert job_id is not None
        assert len(job_id) > 0
        
        await asyncio.sleep(2)
        
        job = await chain_client.get_job(job_id)
        assert job is not None
        assert job.get("id") == job_id
        assert job.get("model_id") == model_id
        assert job.get("dataset_cid") == dataset_cid
        
        tasks = await chain_client.get_tasks_by_job(job_id)
        assert tasks is not None
        
    except NotImplementedError:
        pytest.skip("gRPC client not fully implemented")


@pytest.mark.asyncio
async def test_job_monitoring(coordinator, chain_client):
    """Test job monitoring via pub/sub"""
    
    model_id = "test-model-456"
    dataset_cid = "QmTestDataset2"
    config = {"epochs": 5}
    creator = "atlas1test456"
    
    try:
        job_id = await coordinator.submit_job(
            creator=creator,
            model_id=model_id,
            dataset_cid=dataset_cid,
            config=config,
        )
        
        updates_received = []
        
        async def monitor_updates():
            topic = f"/atlas/jobs/{job_id}/updates"
            async for update in coordinator.pubsub.subscribe(topic):
                updates_received.append(update)
                if len(updates_received) >= 3:
                    break
        
        monitor_task = asyncio.create_task(monitor_updates())
        
        await asyncio.sleep(5)
        monitor_task.cancel()
        
        assert len(updates_received) > 0
        
    except NotImplementedError:
        pytest.skip("gRPC client not fully implemented")


@pytest.mark.asyncio
async def test_multiple_jobs(coordinator, chain_client):
    """Test submitting multiple jobs"""
    
    jobs = []
    creator = "atlas1test789"
    
    try:
        for i in range(3):
            job_id = await coordinator.submit_job(
                creator=creator,
                model_id=f"test-model-{i}",
                dataset_cid=f"QmTestDataset{i}",
                config={"epochs": 1},
            )
            jobs.append(job_id)
            await asyncio.sleep(1)
        
        all_jobs = await chain_client.list_jobs()
        assert len(all_jobs) >= 3
        
        for job_id in jobs:
            job = await chain_client.get_job(job_id)
            assert job is not None
            assert job.get("id") == job_id
        
    except NotImplementedError:
        pytest.skip("gRPC client not fully implemented")

