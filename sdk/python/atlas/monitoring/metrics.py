"""Metrics Collection"""

import time
from typing import Dict, List, Optional
from dataclasses import dataclass, field
from collections import defaultdict
from threading import Lock


@dataclass
class MetricValue:
    value: float
    timestamp: float
    labels: Dict[str, str] = field(default_factory=dict)


class Counter:
    def __init__(self, name: str, description: str = ""):
        self.name = name
        self.description = description
        self.values: Dict[str, float] = defaultdict(float)
        self._lock = Lock()
    
    def inc(self, value: float = 1.0, labels: Optional[Dict[str, str]] = None):
        key = self._get_key(labels)
        with self._lock:
            self.values[key] += value
    
    def get(self, labels: Optional[Dict[str, str]] = None) -> float:
        key = self._get_key(labels)
        return self.values.get(key, 0.0)
    
    def _get_key(self, labels: Optional[Dict[str, str]]) -> str:
        if labels:
            return ",".join(f"{k}={v}" for k, v in sorted(labels.items()))
        return ""


class Gauge:
    def __init__(self, name: str, description: str = ""):
        self.name = name
        self.description = description
        self.values: Dict[str, float] = defaultdict(float)
        self._lock = Lock()
    
    def set(self, value: float, labels: Optional[Dict[str, str]] = None):
        key = self._get_key(labels)
        with self._lock:
            self.values[key] = value
    
    def inc(self, value: float = 1.0, labels: Optional[Dict[str, str]] = None):
        key = self._get_key(labels)
        with self._lock:
            self.values[key] += value
    
    def dec(self, value: float = 1.0, labels: Optional[Dict[str, str]] = None):
        key = self._get_key(labels)
        with self._lock:
            self.values[key] -= value
    
    def get(self, labels: Optional[Dict[str, str]] = None) -> float:
        key = self._get_key(labels)
        return self.values.get(key, 0.0)
    
    def _get_key(self, labels: Optional[Dict[str, str]]) -> str:
        if labels:
            return ",".join(f"{k}={v}" for k, v in sorted(labels.items()))
        return ""


class Histogram:
    def __init__(self, name: str, description: str = "", buckets: Optional[List[float]] = None):
        self.name = name
        self.description = description
        self.buckets = buckets or [0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0, 2.5, 5.0, 10.0]
        self.values: Dict[str, List[float]] = defaultdict(list)
        self._lock = Lock()
    
    def observe(self, value: float, labels: Optional[Dict[str, str]] = None):
        key = self._get_key(labels)
        with self._lock:
            self.values[key].append(value)
    
    def get(self, labels: Optional[Dict[str, str]] = None) -> Dict[str, float]:
        key = self._get_key(labels)
        values = self.values.get(key, [])
        
        if not values:
            return {}
        
        sorted_values = sorted(values)
        count = len(sorted_values)
        
        result = {
            "count": count,
            "sum": sum(sorted_values),
            "min": min(sorted_values),
            "max": max(sorted_values),
        }
        
        if count > 0:
            result["avg"] = result["sum"] / count
            result["p50"] = sorted_values[int(count * 0.5)]
            result["p95"] = sorted_values[int(count * 0.95)]
            result["p99"] = sorted_values[int(count * 0.99)]
        
        return result
    
    def _get_key(self, labels: Optional[Dict[str, str]]) -> str:
        if labels:
            return ",".join(f"{k}={v}" for k, v in sorted(labels.items()))
        return ""


class Timer:
    def __init__(self, histogram: Histogram, labels: Optional[Dict[str, str]] = None):
        self.histogram = histogram
        self.labels = labels
        self.start_time: Optional[float] = None
    
    def __enter__(self):
        self.start_time = time.time()
        return self
    
    def __exit__(self, exc_type, exc_val, exc_tb):
        if self.start_time:
            elapsed = time.time() - self.start_time
            self.histogram.observe(elapsed, self.labels)


class MetricsCollector:
    def __init__(self):
        self.counters: Dict[str, Counter] = {}
        self.gauges: Dict[str, Gauge] = {}
        self.histograms: Dict[str, Histogram] = {}
    
    def counter(self, name: str, description: str = "") -> Counter:
        if name not in self.counters:
            self.counters[name] = Counter(name, description)
        return self.counters[name]
    
    def gauge(self, name: str, description: str = "") -> Gauge:
        if name not in self.gauges:
            self.gauges[name] = Gauge(name, description)
        return self.gauges[name]
    
    def histogram(self, name: str, description: str = "", buckets: Optional[List[float]] = None) -> Histogram:
        if name not in self.histograms:
            self.histograms[name] = Histogram(name, description, buckets)
        return self.histograms[name]
    
    def timer(self, name: str, description: str = "", labels: Optional[Dict[str, str]] = None) -> Timer:
        histogram = self.histogram(name, description)
        return Timer(histogram, labels)


_global_collector = MetricsCollector()


def get_metrics() -> MetricsCollector:
    return _global_collector

