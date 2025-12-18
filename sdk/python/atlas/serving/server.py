"""
Model Server - HTTP/gRPC server for model inference
"""

import asyncio
import time
from typing import Optional
from aiohttp import web
import json

from .loader import ModelLoader
from .predictor import Predictor
from ..chain.client import ChainClient
from ..auth.middleware import require_auth, optional_auth
from ..auth.rbac import Permission


class ModelServer:
    """Model inference server"""
    
    def __init__(
        self,
        chain_client: ChainClient,
        ipfs_api_url: str = "/ip4/127.0.0.1/tcp/5001",
    ):
        """
        Initialize model server
        
        Args:
            chain_client: Blockchain client
            ipfs_api_url: IPFS API URL
        """
        self.chain_client = chain_client
        self.loader = ModelLoader(ipfs_api_url)
        self.predictors: dict[str, Predictor] = {}
        self.app = web.Application()
        self._setup_routes()
    
    def _setup_routes(self):
        """Setup HTTP routes"""
        self.app.router.add_post("/predict", self.handle_predict)
        self.app.router.add_get("/health", self.handle_health)
    
    @optional_auth
    async def handle_predict(self, request: web.Request) -> web.Response:
        """Handle prediction request"""
        try:
            data = await request.json()
            model_id = data.get("model_id")
            input_data = data.get("input")
            model_type = data.get("model_type", "auto")
            options = data.get("options", {})
            
            if not model_id or input_data is None:
                return web.json_response(
                    {"error": "model_id and input required"},
                    status=400
                )
            
            # Get predictor (load if needed)
            if model_id not in self.predictors:
                # Get model CID from blockchain
                model = await self.chain_client.get_model(model_id)
                model_cid = model.get("cid")
                
                if not model_cid:
                    return web.json_response(
                        {"error": f"Model {model_id} not found"},
                        status=404
                    )
                
                # Load model
                model_path = self.loader.load(model_cid)
                self.predictors[model_id] = Predictor(model_path)
            
            # Run prediction
            predictor = self.predictors[model_id]
            start_time = time.time()
            result = predictor.predict(input_data, model_type=model_type, options=options)
            latency_ms = int((time.time() - start_time) * 1000)
            
            return web.json_response({
                "result": result,
                "latency_ms": latency_ms,
                "model_id": model_id
            })
        except Exception as e:
            return web.json_response(
                {"error": str(e)},
                status=500
            )
    
    async def handle_health(self, request: web.Request) -> web.Response:
        """Handle health check"""
        return web.json_response({"status": "ok"})
    
    async def serve(self, model_id: str, host: str = "0.0.0.0", port: int = 8000):
        """
        Start serving model
        
        Args:
            model_id: Model ID to serve
            host: Server host
            port: Server port
        """
        # Pre-load model
        model = await self.chain_client.get_model(model_id)
        model_cid = model.get("cid")
        if model_cid:
            model_path = self.loader.load(model_cid)
            self.predictors[model_id] = Predictor(model_path)
        
        print(f"Starting model server on {host}:{port}")
        runner = web.AppRunner(self.app)
        await runner.setup()
        site = web.TCPSite(runner, host, port)
        await site.start()
        
        print(f"Model server running. Send POST requests to http://{host}:{port}/predict")
        
        # Keep running
        try:
            await asyncio.Event().wait()
        except KeyboardInterrupt:
            print("Shutting down server...")
            await runner.cleanup()

