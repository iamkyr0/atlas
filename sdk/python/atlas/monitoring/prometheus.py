"""Prometheus Exporter"""

from typing import Dict
from aiohttp import web
from .metrics import MetricsCollector, get_metrics


class PrometheusExporter:
    def __init__(self, metrics: MetricsCollector = None):
        self.metrics = metrics or get_metrics()
    
    def export(self) -> str:
        lines = []
        
        for counter in self.metrics.counters.values():
            lines.append(f"# HELP {counter.name} {counter.description}")
            lines.append(f"# TYPE {counter.name} counter")
            for key, value in counter.values.items():
                labels = self._parse_key(key)
                label_str = self._format_labels(labels)
                lines.append(f"{counter.name}{{{label_str}}} {value}")
        
        for gauge in self.metrics.gauges.values():
            lines.append(f"# HELP {gauge.name} {gauge.description}")
            lines.append(f"# TYPE {gauge.name} gauge")
            for key, value in gauge.values.items():
                labels = self._parse_key(key)
                label_str = self._format_labels(labels)
                lines.append(f"{gauge.name}{{{label_str}}} {value}")
        
        for histogram in self.metrics.histograms.values():
            lines.append(f"# HELP {histogram.name} {histogram.description}")
            lines.append(f"# TYPE {histogram.name} histogram")
            for key, values in histogram.values.items():
                labels = self._parse_key(key)
                label_str = self._format_labels(labels)
                stats = histogram.get(labels)
                if stats:
                    lines.append(f"{histogram.name}_count{{{label_str}}} {stats['count']}")
                    lines.append(f"{histogram.name}_sum{{{label_str}}} {stats['sum']}")
                    lines.append(f"{histogram.name}_avg{{{label_str}}} {stats.get('avg', 0)}")
        
        return "\n".join(lines)
    
    def _parse_key(self, key: str) -> Dict[str, str]:
        if not key:
            return {}
        labels = {}
        for part in key.split(","):
            if "=" in part:
                k, v = part.split("=", 1)
                labels[k] = v
        return labels
    
    def _format_labels(self, labels: Dict[str, str]) -> str:
        if not labels:
            return ""
        return ",".join(f'{k}="{v}"' for k, v in sorted(labels.items()))
    
    async def metrics_handler(self, request: web.Request) -> web.Response:
        return web.Response(
            text=self.export(),
            content_type="text/plain; version=0.0.4",
        )

