"""
Atlas CLI - Command-line interface for Atlas platform
"""

import sys
import argparse
import asyncio
from typing import Optional
import os
from pathlib import Path

from ..chain.client import ChainClient
from ..p2p.pubsub import PubSub
from ..p2p.discovery import NodeDiscovery
from ..job.coordinator import JobCoordinator
from .commands import (
    SubmitJobCommand,
    ListJobsCommand,
    GetJobCommand,
    UploadDatasetCommand,
    RegisterModelCommand,
    DownloadModelCommand,
    ServeModelCommand,
    DaemonCommand,
)


class AtlasCLI:
    """Main CLI application"""
    
    def __init__(self):
        """Initialize CLI"""
        self.parser = argparse.ArgumentParser(
            prog="atlas",
            description="Atlas Decentralized AI Platform CLI",
        )
        self.parser.add_argument(
            "--ipfs-api",
            default=os.getenv("ATLAS_IPFS_API", "/ip4/127.0.0.1/tcp/5001"),
            help="IPFS API URL",
        )
        self.parser.add_argument(
            "--chain-grpc",
            default=os.getenv("ATLAS_CHAIN_GRPC", "localhost:9090"),
            help="Chain gRPC URL",
        )
        self.parser.add_argument(
            "--creator",
            default=os.getenv("ATLAS_CREATOR", ""),
            help="Creator address for transactions",
        )
        
        # Subcommands
        subparsers = self.parser.add_subparsers(dest="command", help="Commands")
        
        # Submit job
        submit_parser = subparsers.add_parser("submit-job", help="Submit a training job")
        submit_parser.add_argument("model_id", help="Model ID")
        submit_parser.add_argument("dataset_cid", help="Dataset IPFS CID")
        submit_parser.add_argument("--config", help="Config JSON file path")
        
        # List jobs
        subparsers.add_parser("list-jobs", help="List all jobs")
        
        # Get job
        get_job_parser = subparsers.add_parser("get-job", help="Get job status")
        get_job_parser.add_argument("job_id", help="Job ID")
        
        # Upload dataset
        upload_parser = subparsers.add_parser("upload-dataset", help="Upload dataset to IPFS")
        upload_parser.add_argument("path", help="Path to dataset file")
        upload_parser.add_argument("--encrypt", action="store_true", help="Encrypt dataset")
        
        # Register model
        register_parser = subparsers.add_parser("register-model", help="Register a model")
        register_parser.add_argument("path", help="Path to model file")
        register_parser.add_argument("--name", required=True, help="Model name")
        register_parser.add_argument("--version", required=True, help="Model version")
        
        # Download model
        download_parser = subparsers.add_parser("download-model", help="Download model")
        download_parser.add_argument("model_id", help="Model ID")
        download_parser.add_argument("--output", help="Output directory")
        
        # Serve model
        serve_parser = subparsers.add_parser("serve-model", help="Serve model for inference")
        serve_parser.add_argument("model_id", help="Model ID")
        serve_parser.add_argument("--port", type=int, default=8000, help="Server port")
        serve_parser.add_argument("--host", default="0.0.0.0", help="Server host")
        
        # Daemon
        daemon_parser = subparsers.add_parser("daemon", help="Start daemon mode")
        daemon_parser.add_argument("--port", type=int, default=8080, help="Daemon port")
        daemon_parser.add_argument("--host", default="127.0.0.1", help="Daemon host")
    
    def run(self, args: Optional[list] = None):
        """Run CLI"""
        parsed_args = self.parser.parse_args(args)
        
        if not parsed_args.command:
            self.parser.print_help()
            return
        
        # Initialize clients
        chain_client = ChainClient(grpc_url=parsed_args.chain_grpc)
        pubsub = PubSub(ipfs_api_url=parsed_args.ipfs_api)
        discovery = NodeDiscovery(
            ipfs_api_url=parsed_args.ipfs_api,
            chain_client=chain_client,
        )
        coordinator = JobCoordinator(
            chain_client=chain_client,
            pubsub=pubsub,
            discovery=discovery,
        )
        
        # Execute command
        try:
            if parsed_args.command == "submit-job":
                cmd = SubmitJobCommand(coordinator, parsed_args.creator)
                asyncio.run(cmd.execute(
                    model_id=parsed_args.model_id,
                    dataset_cid=parsed_args.dataset_cid,
                    config_path=parsed_args.config,
                ))
            elif parsed_args.command == "list-jobs":
                cmd = ListJobsCommand(chain_client)
                asyncio.run(cmd.execute())
            elif parsed_args.command == "get-job":
                cmd = GetJobCommand(chain_client)
                asyncio.run(cmd.execute(parsed_args.job_id))
            elif parsed_args.command == "upload-dataset":
                cmd = UploadDatasetCommand(pubsub, parsed_args.ipfs_api)
                asyncio.run(cmd.execute(
                    path=parsed_args.path,
                    encrypt=parsed_args.encrypt,
                ))
            elif parsed_args.command == "register-model":
                cmd = RegisterModelCommand(chain_client, pubsub, parsed_args.ipfs_api)
                asyncio.run(cmd.execute(
                    path=parsed_args.path,
                    name=parsed_args.name,
                    version=parsed_args.version,
                    creator=parsed_args.creator,
                ))
            elif parsed_args.command == "download-model":
                cmd = DownloadModelCommand(chain_client, parsed_args.ipfs_api)
                asyncio.run(cmd.execute(
                    model_id=parsed_args.model_id,
                    output=parsed_args.output,
                ))
            elif parsed_args.command == "serve-model":
                cmd = ServeModelCommand(chain_client, parsed_args.ipfs_api)
                asyncio.run(cmd.execute(
                    model_id=parsed_args.model_id,
                    host=parsed_args.host,
                    port=parsed_args.port,
                ))
            elif parsed_args.command == "daemon":
                cmd = DaemonCommand(
                    chain_client,
                    pubsub,
                    discovery,
                    coordinator,
                )
                asyncio.run(cmd.execute(
                    host=parsed_args.host,
                    port=parsed_args.port,
                ))
        finally:
            chain_client.close()


def main():
    """Main entry point"""
    cli = AtlasCLI()
    cli.run()


if __name__ == "__main__":
    main()

