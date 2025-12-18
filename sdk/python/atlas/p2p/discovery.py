import asyncio
import time
from typing import List, Dict, Any, Optional
from ipfshttpclient import Client as IPFSClient
from ..chain.client import ChainClient


class NodeDiscovery:
    
    def __init__(
        self,
        ipfs_api_url: str = "/ip4/127.0.0.1/tcp/5001",
        chain_client: Optional[ChainClient] = None,
        cache_ttl: int = 300,
    ):
        self.ipfs_client = IPFSClient(ipfs_api_url)
        self.chain_client = chain_client
        self.cache_ttl = cache_ttl
        self._node_cache: Dict[str, Dict[str, Any]] = {}
        self._cache_timestamps: Dict[str, float] = {}
    
    async def discover_nodes(self, max_nodes: Optional[int] = None) -> List[Dict[str, Any]]:
        nodes = []
        
        try:
            nodes = await self._discover_via_ipfs_dht(max_nodes)
            if nodes:
                return nodes
        except Exception as e:
            print(f"IPFS DHT discovery failed: {e}, falling back to blockchain")
        
        if self.chain_client:
            try:
                nodes = await self._discover_via_blockchain(max_nodes)
            except Exception as e:
                print(f"Blockchain discovery failed: {e}")
        
        return nodes
    
    async def _discover_via_ipfs_dht(self, max_nodes: Optional[int] = None) -> List[Dict[str, Any]]:
        nodes = []
        
        cached_nodes = self._get_cached_nodes()
        if cached_nodes:
            return cached_nodes[:max_nodes] if max_nodes else cached_nodes
        
        return nodes
    
    async def _discover_via_blockchain(self, max_nodes: Optional[int] = None) -> List[Dict[str, Any]]:
        if not self.chain_client:
            return []
        
        try:
            nodes_data = await self.chain_client.list_nodes()
            
            nodes = []
            for node_data in nodes_data:
                node_dict = {
                    "id": node_data.get("id"),
                    "address": node_data.get("address"),
                    "status": node_data.get("status"),
                    "resources": node_data.get("resources", {}),
                    "reputation": node_data.get("reputation", 0.0),
                }
                nodes.append(node_dict)
            
            self._cache_nodes(nodes)
            
            return nodes[:max_nodes] if max_nodes else nodes
        except Exception as e:
            print(f"Error querying blockchain for nodes: {e}")
            return []
    
    def _get_cached_nodes(self) -> List[Dict[str, Any]]:
        now = time.time()
        valid_nodes = []
        
        for node_id, timestamp in self._cache_timestamps.items():
            if now - timestamp < self.cache_ttl:
                if node_id in self._node_cache:
                    valid_nodes.append(self._node_cache[node_id])
        
        return valid_nodes
    
    def _cache_nodes(self, nodes: List[Dict[str, Any]]) -> None:
        now = time.time()
        for node in nodes:
            node_id = node.get("id")
            if node_id:
                self._node_cache[node_id] = node
                self._cache_timestamps[node_id] = now
    
    async def health_check(self, node_id: str) -> bool:
        if node_id in self._node_cache:
            node = self._node_cache[node_id]
            return node.get("status") == "online"
        
        if self.chain_client:
            try:
                node_data = await self.chain_client.get_node(node_id)
                return node_data.get("status") == "online"
            except Exception:
                return False
        
        return False

