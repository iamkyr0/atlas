"""Structured Logging"""

import logging
import sys
import json
from enum import Enum
from typing import Optional, Dict, Any
from datetime import datetime


class LogLevel(Enum):
    DEBUG = logging.DEBUG
    INFO = logging.INFO
    WARNING = logging.WARNING
    ERROR = logging.ERROR
    CRITICAL = logging.CRITICAL


class AtlasLogger:
    def __init__(self, name: str, level: LogLevel = LogLevel.INFO):
        self.logger = logging.getLogger(name)
        self.logger.setLevel(level.value)
        self.name = name
        self._context: Dict[str, Any] = {}
    
    def set_context(self, **kwargs):
        self._context.update(kwargs)
    
    def clear_context(self):
        self._context.clear()
    
    def debug(self, message: str, **kwargs):
        self._log(logging.DEBUG, message, **kwargs)
    
    def info(self, message: str, **kwargs):
        self._log(logging.INFO, message, **kwargs)
    
    def warning(self, message: str, **kwargs):
        self._log(logging.WARNING, message, **kwargs)
    
    def error(self, message: str, **kwargs):
        self._log(logging.ERROR, message, **kwargs)
    
    def critical(self, message: str, **kwargs):
        self._log(logging.CRITICAL, message, **kwargs)
    
    def _log(self, level: int, message: str, **kwargs):
        extra = {**self._context, **kwargs}
        self.logger.log(level, message, extra=extra)


_loggers: Dict[str, AtlasLogger] = {}


def get_logger(name: str, level: Optional[LogLevel] = None) -> AtlasLogger:
    if name not in _loggers:
        _loggers[name] = AtlasLogger(name, level or LogLevel.INFO)
    return _loggers[name]


def setup_logging(
    level: LogLevel = LogLevel.INFO,
    format_type: str = "json",
    output: Optional[str] = None,
):
    root_logger = logging.getLogger()
    root_logger.setLevel(level.value)
    
    handler = logging.StreamHandler(sys.stdout) if output is None else logging.FileHandler(output)
    
    if format_type == "json":
        from .formatter import JSONFormatter
        handler.setFormatter(JSONFormatter())
    else:
        from .formatter import TextFormatter
        handler.setFormatter(TextFormatter())
    
    root_logger.addHandler(handler)

