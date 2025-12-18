"""
Model Loader - Load models from IPFS
"""

import os
import tempfile
from pathlib import Path
from typing import Optional, Any
from ipfshttpclient import Client as IPFSClient


class ModelLoader:
    """Load models from IPFS"""
    
    def __init__(self, ipfs_api_url: str = "/ip4/127.0.0.1/tcp/5001"):
        """
        Initialize model loader
        
        Args:
            ipfs_api_url: IPFS API URL
        """
        self.ipfs_client = IPFSClient(ipfs_api_url)
        self._cache: dict[str, str] = {}  # model_id -> local_path
    
    def load(self, model_cid: str, cache_dir: Optional[str] = None) -> str:
        """
        Load model from IPFS
        
        Args:
            model_cid: IPFS CID of the model
            cache_dir: Cache directory (default: temp directory)
            
        Returns:
            Local path to model
        """
        # Check cache
        if model_cid in self._cache:
            cached_path = self._cache[model_cid]
            if os.path.exists(cached_path):
                return cached_path
        
        # Determine cache directory
        if not cache_dir:
            cache_dir = os.path.join(tempfile.gettempdir(), "atlas_models")
        
        cache_path = Path(cache_dir)
        cache_path.mkdir(parents=True, exist_ok=True)
        
        # Download from IPFS
        model_dir = cache_path / model_cid
        if not model_dir.exists():
            self.ipfs_client.get(model_cid, str(cache_path))
        
        # Find model file (could be .pth, .pt, .h5, etc.)
        model_files = list(model_dir.rglob("*"))
        if not model_files:
            raise ValueError(f"No model files found in {model_dir}")
        
        # Return first file (or could implement smarter detection)
        model_path = str(model_files[0])
        self._cache[model_cid] = model_path
        
        return model_path
    
    def get_model_path(self, model_cid: str) -> Optional[str]:
        """
        Get cached model path
        
        Args:
            model_cid: Model CID
            
        Returns:
            Cached path or None
        """
        return self._cache.get(model_cid)

