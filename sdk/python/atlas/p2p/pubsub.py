"""
IPFS Pub/Sub Wrapper - Real-time P2P communication
"""

import asyncio
from typing import AsyncIterator, Optional, Callable
from ipfshttpclient import Client as IPFSClient
import json


class PubSub:
    """IPFS Pub/Sub wrapper for real-time communication"""
    
    def __init__(self, ipfs_api_url: str = "/ip4/127.0.0.1/tcp/5001"):
        """
        Initialize IPFS pub/sub
        
        Args:
            ipfs_api_url: IPFS API URL
        """
        self.ipfs_client = IPFSClient(ipfs_api_url)
    
    async def publish(self, topic: str, message: dict) -> None:
        """
        Publish message to topic
        
        Args:
            topic: Topic name
            message: Message dict to publish
        """
        message_str = json.dumps(message)
        self.ipfs_client.pubsub.pub(topic, message_str)
    
    async def subscribe(self, topic: str) -> AsyncIterator[dict]:
        """
        Subscribe to topic and yield messages
        
        Args:
            topic: Topic name
            
        Yields:
            Message dicts
        """
        # IPFS pub/sub is synchronous, so we need to wrap it
        for raw_message in self.ipfs_client.pubsub.sub(topic):
            try:
                message = json.loads(raw_message["data"])
                yield message
            except (json.JSONDecodeError, KeyError) as e:
                # Skip invalid messages
                continue
    
    async def subscribe_with_handler(
        self,
        topic: str,
        handler: Callable[[dict], None],
    ) -> None:
        """
        Subscribe to topic with callback handler
        
        Args:
            topic: Topic name
            handler: Async callback function
        """
        async for message in self.subscribe(topic):
            await handler(message)
    
    def unsubscribe(self, topic: str) -> None:
        """
        Unsubscribe from topic
        
        Args:
            topic: Topic name
        """
        # IPFS pub/sub doesn't have explicit unsubscribe
        # The subscription will end when the generator is closed
        pass

