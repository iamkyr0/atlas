import asyncio
from aiohttp import web
import json
from typing import Dict, Any

from ..chain.client import ChainClient
from ..p2p.pubsub import PubSub
from ..p2p.discovery import NodeDiscovery
from ..job.coordinator import JobCoordinator
from ..serving.server import ModelServer
from ..auth.middleware import require_auth, optional_auth
from ..auth.rbac import Permission


class DaemonServer:
    
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
        self.model_server = ModelServer(chain_client)
        self.app = web.Application()
        self._setup_routes()
    
    def _setup_routes(self):
        self.app.router.add_post("/api/v1/jobs", self.handle_submit_job)
        self.app.router.add_get("/api/v1/jobs", self.handle_list_jobs)
        self.app.router.add_get("/api/v1/jobs/{job_id}", self.handle_get_job)
        self.app.router.add_post("/api/v1/models", self.handle_register_model)
        self.app.router.add_get("/api/v1/models/{model_id}", self.handle_get_model)
        self.app.router.add_post("/api/v1/predict", self.handle_predict)
        self.app.router.add_get("/health", self.handle_health)
    
    @require_auth(permission=Permission.SUBMIT_JOB)
    async def handle_submit_job(self, request: web.Request) -> web.Response:
        try:
            data = await request.json()
            creator = data.get("creator", "")
            model_id = data.get("model_id")
            dataset_cid = data.get("dataset_cid")
            config = data.get("config", {})
            
            if not model_id or not dataset_cid:
                return web.json_response(
                    {"error": "model_id and dataset_cid required"},
                    status=400
                )
            
            job_id = await self.coordinator.submit_job(
                creator=creator,
                model_id=model_id,
                dataset_cid=dataset_cid,
                config=config,
            )
            
            return web.json_response({"job_id": job_id})
        except Exception as e:
            return web.json_response({"error": str(e)}, status=500)
    
    @require_auth(permission=Permission.LIST_JOBS)
    async def handle_list_jobs(self, request: web.Request) -> web.Response:
        try:
            jobs = await self.chain_client.list_jobs()
            return web.json_response({"jobs": jobs})
        except Exception as e:
            return web.json_response({"error": str(e)}, status=500)
    
    @require_auth(permission=Permission.VIEW_JOB)
    async def handle_get_job(self, request: web.Request) -> web.Response:
        try:
            job_id = request.match_info["job_id"]
            job = await self.chain_client.get_job(job_id)
            return web.json_response(job)
        except Exception as e:
            return web.json_response({"error": str(e)}, status=500)
    
    @require_auth(permission=Permission.REGISTER_MODEL)
    async def handle_register_model(self, request: web.Request) -> web.Response:
        try:
            data = await request.json()
            return web.json_response({"error": "Not implemented"}, status=501)
        except Exception as e:
            return web.json_response({"error": str(e)}, status=500)
    
    @require_auth(permission=Permission.VIEW_MODEL)
    async def handle_get_model(self, request: web.Request) -> web.Response:
        try:
            model_id = request.match_info["model_id"]
            model = await self.chain_client.get_model(model_id)
            return web.json_response(model)
        except Exception as e:
            return web.json_response({"error": str(e)}, status=500)
    
    @optional_auth
    async def handle_predict(self, request: web.Request) -> web.Response:
        return await self.model_server.handle_predict(request)
    
    async def handle_health(self, request: web.Request) -> web.Response:
        return web.json_response({"status": "ok"})
    
    async def start(self, host: str = "127.0.0.1", port: int = 8080):
        print(f"Starting Atlas daemon on {host}:{port}")
        runner = web.AppRunner(self.app)
        await runner.setup()
        site = web.TCPSite(runner, host, port)
        await site.start()
        
        print(f"Daemon running. API available at http://{host}:{port}/api/v1")
        
        try:
            await asyncio.Event().wait()
        except KeyboardInterrupt:
            print("Shutting down daemon...")
            await runner.cleanup()

