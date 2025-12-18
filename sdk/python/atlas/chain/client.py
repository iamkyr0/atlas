import grpc
from typing import Dict, Any, List, Optional
import json
from google.protobuf import struct_pb2, message
from google.protobuf.json_format import MessageToDict, ParseDict
import asyncio

try:
    from google.protobuf import message
except ImportError:
    message = None


class ChainClient:
    
    def __init__(self, grpc_url: str = "localhost:9090", use_tls: bool = False, api_key: Optional[str] = None):
        self.grpc_url = grpc_url
        self.use_tls = use_tls
        self.api_key = api_key
        self.channel: Optional[grpc.Channel] = None
        self._connect()
    
    def _connect(self):
        if self.grpc_url.startswith("https://"):
            self.grpc_url = self.grpc_url.replace("https://", "")
            self.use_tls = True
        elif self.grpc_url.startswith("http://"):
            self.grpc_url = self.grpc_url.replace("http://", "")
            self.use_tls = False
        
        options = []
        if self.api_key:
            metadata = [('x-api-key', self.api_key)]
            options.append(('grpc.primary_user_agent', 'atlas-python-sdk'))
        
        if self.use_tls:
            credentials = grpc.ssl_channel_credentials()
            self.channel = grpc.secure_channel(self.grpc_url, credentials, options=options)
        else:
            self.channel = grpc.insecure_channel(self.grpc_url, options=options)
    
    def close(self):
        if self.channel:
            self.channel.close()
    
    def _call_rpc(self, service: str, method: str, request_dict: Dict[str, Any]) -> Dict[str, Any]:
        try:
            from google.protobuf import reflection
            from google.protobuf import descriptor_pb2
            
            stub_class = self._get_stub_class(service)
            if stub_class is None:
                raise NotImplementedError(f"Service {service} not available. Using gRPC reflection or REST fallback.")
            
            stub = stub_class(self.channel)
            method_func = getattr(stub, method, None)
            if method_func is None:
                raise NotImplementedError(f"Method {method} not found in service {service}")
            
            request_msg = self._dict_to_protobuf(request_dict, service, method, is_request=True)
            
            metadata = []
            if self.api_key:
                metadata.append(('x-api-key', self.api_key))
            
            response = method_func(request_msg, metadata=metadata)
            return MessageToDict(response, including_default_value_fields=True)
        except Exception as e:
            raise NotImplementedError(
                f"gRPC call failed: {str(e)}. "
                "Please ensure protobuf stubs are generated or use REST API fallback."
            )
    
    def _get_stub_class(self, service: str):
        try:
            if service == "training.Query":
                from atlas.chain import training_query_pb2_grpc
                return training_query_pb2_grpc.QueryStub
            elif service == "training.Msg":
                from atlas.chain import training_tx_pb2_grpc
                return training_tx_pb2_grpc.MsgStub
            elif service == "compute.Query":
                from atlas.chain import compute_query_pb2_grpc
                return compute_query_pb2_grpc.QueryStub
            elif service == "model.Query":
                from atlas.chain import model_query_pb2_grpc
                return model_query_pb2_grpc.QueryStub
        except ImportError:
            pass
        return None
    
    def _dict_to_protobuf(self, data: Dict[str, Any], service: str, method: str, is_request: bool = True):
        try:
            if service == "training.Query":
                if method == "GetJob":
                    from atlas.chain import training_query_pb2
                    msg = training_query_pb2.QueryGetJobRequest()
                    msg.job_id = data.get("job_id", "")
                    return msg
                elif method == "ListJobs":
                    from atlas.chain import training_query_pb2
                    return training_query_pb2.QueryListJobsRequest()
                elif method == "GetTask":
                    from atlas.chain import training_query_pb2
                    msg = training_query_pb2.QueryGetTaskRequest()
                    msg.task_id = data.get("task_id", "")
                    return msg
                elif method == "GetTasksByJob":
                    from atlas.chain import training_query_pb2
                    msg = training_query_pb2.QueryGetTasksByJobRequest()
                    msg.job_id = data.get("job_id", "")
                    return msg
            elif service == "training.Msg":
                if method == "SubmitJob":
                    from atlas.chain import training_tx_pb2
                    msg = training_tx_pb2.MsgSubmitJob()
                    msg.creator = data.get("creator", "")
                    msg.model_id = data.get("model_id", "")
                    msg.dataset_cid = data.get("dataset_cid", "")
                    if "config" in data:
                        for k, v in data["config"].items():
                            msg.config[k] = str(v)
                    return msg
            elif service == "compute.Query":
                if method == "GetNode":
                    from atlas.chain import compute_query_pb2
                    msg = compute_query_pb2.QueryGetNodeRequest()
                    msg.node_id = data.get("node_id", "")
                    return msg
                elif method == "ListNodes":
                    from atlas.chain import compute_query_pb2
                    return compute_query_pb2.QueryListNodesRequest()
            elif service == "model.Query":
                if method == "GetModel":
                    from atlas.chain import model_query_pb2
                    msg = model_query_pb2.QueryGetModelRequest()
                    msg.model_id = data.get("model_id", "")
                    return msg
                elif method == "ListModels":
                    from atlas.chain import model_query_pb2
                    return model_query_pb2.QueryListModelsRequest()
        except ImportError:
            pass
        
        struct_msg = struct_pb2.Struct()
        struct_msg.update(data)
        return struct_msg
    
    async def submit_job(
        self,
        creator: str,
        model_id: str,
        dataset_cid: str,
        config: Dict[str, Any],
    ) -> str:
        try:
            request = {
                "creator": creator,
                "model_id": model_id,
                "dataset_cid": dataset_cid,
                "config": config,
            }
            response = await asyncio.to_thread(
                self._call_rpc, "training.Msg", "SubmitJob", request
            )
            return response.get("job_id", "")
        except NotImplementedError:
            raise NotImplementedError(
                "submit_job requires generated protobuf stubs. "
                "Please generate Python stubs from .proto files using: "
                "python -m grpc_tools.protoc -I. --python_out=. --grpc_python_out=. *.proto"
            )
    
    async def get_job(self, job_id: str) -> Dict[str, Any]:
        try:
            request = {"job_id": job_id}
            response = await asyncio.to_thread(
                self._call_rpc, "training.Query", "GetJob", request
            )
            return response.get("job", {})
        except NotImplementedError:
            raise NotImplementedError("Requires generated protobuf stubs")
    
    async def list_jobs(self) -> List[Dict[str, Any]]:
        try:
            request = {}
            response = await asyncio.to_thread(
                self._call_rpc, "training.Query", "ListJobs", request
            )
            return response.get("jobs", [])
        except NotImplementedError:
            raise NotImplementedError("Requires generated protobuf stubs")
    
    async def get_task(self, task_id: str) -> Dict[str, Any]:
        try:
            request = {"task_id": task_id}
            response = await asyncio.to_thread(
                self._call_rpc, "training.Query", "GetTask", request
            )
            return response.get("task", {})
        except NotImplementedError:
            raise NotImplementedError("Requires generated protobuf stubs")
    
    async def get_tasks_by_job(self, job_id: str) -> List[Dict[str, Any]]:
        try:
            request = {"job_id": job_id}
            response = await asyncio.to_thread(
                self._call_rpc, "training.Query", "GetTasksByJob", request
            )
            return response.get("tasks", [])
        except NotImplementedError:
            raise NotImplementedError("Requires generated protobuf stubs")
    
    async def register_model(
        self,
        creator: str,
        name: str,
        version: str,
        cid: str,
        metadata: Dict[str, str],
    ) -> str:
        raise NotImplementedError("Model registration via gRPC not yet implemented")
    
    async def get_model(self, model_id: str) -> Dict[str, Any]:
        try:
            request = {"model_id": model_id}
            response = await asyncio.to_thread(
                self._call_rpc, "model.Query", "GetModel", request
            )
            model = response.get("model", {})
            if model:
                return {
                    "id": model.get("id", ""),
                    "name": model.get("name", ""),
                    "version": model.get("version", ""),
                    "cid": model.get("cid", ""),
                    "created_at": model.get("created_at", 0),
                    "metadata": model.get("metadata", {}),
                }
            return {}
        except NotImplementedError:
            raise NotImplementedError("Requires generated protobuf stubs")
    
    async def list_models(self) -> List[Dict[str, Any]]:
        try:
            request = {}
            response = await asyncio.to_thread(
                self._call_rpc, "model.Query", "ListModels", request
            )
            models = response.get("models", [])
            return [
                {
                    "id": m.get("id", ""),
                    "name": m.get("name", ""),
                    "version": m.get("version", ""),
                    "cid": m.get("cid", ""),
                    "created_at": m.get("created_at", 0),
                    "metadata": m.get("metadata", {}),
                }
                for m in models
            ]
        except NotImplementedError:
            raise NotImplementedError("Requires generated protobuf stubs")
    
    async def get_node(self, node_id: str) -> Dict[str, Any]:
        try:
            request = {"node_id": node_id}
            response = await asyncio.to_thread(
                self._call_rpc, "compute.Query", "GetNode", request
            )
            return response.get("node", {})
        except NotImplementedError:
            raise NotImplementedError("Requires generated protobuf stubs")
    
    async def list_nodes(self) -> List[Dict[str, Any]]:
        try:
            request = {}
            response = await asyncio.to_thread(
                self._call_rpc, "compute.Query", "ListNodes", request
            )
            return response.get("nodes", [])
        except NotImplementedError:
            raise NotImplementedError("Requires generated protobuf stubs")
