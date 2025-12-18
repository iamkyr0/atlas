"""Configuration Management"""

import os
import json
import yaml
from typing import Dict, Any, Optional
from pathlib import Path
from dataclasses import dataclass, field


@dataclass
class ConfigManager:
    config: Dict[str, Any] = field(default_factory=dict)
    env_prefix: str = "ATLAS_"
    
    def load_from_file(self, file_path: str):
        path = Path(file_path)
        
        if not path.exists():
            raise FileNotFoundError(f"Config file not found: {file_path}")
        
        if path.suffix in [".yaml", ".yml"]:
            with open(path, "r") as f:
                self.config.update(yaml.safe_load(f) or {})
        elif path.suffix == ".json":
            with open(path, "r") as f:
                self.config.update(json.load(f))
        else:
            raise ValueError(f"Unsupported config file format: {path.suffix}")
    
    def load_from_env(self, prefix: Optional[str] = None):
        prefix = prefix or self.env_prefix
        
        for key, value in os.environ.items():
            if key.startswith(prefix):
                config_key = key[len(prefix):].lower().replace("_", ".")
                self._set_nested(self.config, config_key, value)
    
    def _set_nested(self, d: Dict, key: str, value: Any):
        keys = key.split(".")
        for k in keys[:-1]:
            d = d.setdefault(k, {})
        d[keys[-1]] = self._convert_value(value)
    
    def _convert_value(self, value: str) -> Any:
        if value.lower() == "true":
            return True
        elif value.lower() == "false":
            return False
        elif value.isdigit():
            return int(value)
        elif value.replace(".", "", 1).isdigit():
            return float(value)
        return value
    
    def get(self, key: str, default: Any = None) -> Any:
        keys = key.split(".")
        value = self.config
        
        for k in keys:
            if isinstance(value, dict):
                value = value.get(k)
                if value is None:
                    return default
            else:
                return default
        
        return value
    
    def set(self, key: str, value: Any):
        self._set_nested(self.config, key, value)
    
    def validate(self, required_keys: list[str]) -> bool:
        for key in required_keys:
            if self.get(key) is None:
                return False
        return True


_global_config = ConfigManager()


def load_config(
    file_path: Optional[str] = None,
    env_prefix: str = "ATLAS_",
    load_env: bool = True,
) -> ConfigManager:
    if file_path:
        _global_config.load_from_file(file_path)
    
    if load_env:
        _global_config.load_from_env(env_prefix)
    
    return _global_config


def get_config() -> ConfigManager:
    return _global_config

