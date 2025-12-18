"""
P2P Messaging Layer - Direct peer-to-peer communication
"""

from typing import Dict, Any, Optional, AsyncIterator
from .pubsub import PubSub


class P2PMessaging:
    """P2P messaging layer for direct communication"""
    
    def __init__(self, pubsub: PubSub):
        """
        Initialize P2P messaging
        
        Args:
            pubsub: PubSub instance
        """
        self.pubsub = pubsub
    
    async def send_to_node(self, node_id: str, message: Dict[str, Any]) -> None:
        """
        Send message directly to a node
        
        Args:
            node_id: Target node ID
            message: Message dict
        """
        topic = f"/atlas/node/{node_id}"
        await self.pubsub.publish(topic, message)
    
    async def broadcast(self, message: Dict[str, Any], topic_suffix: str = "broadcast") -> None:
        """
        Broadcast message to all nodes
        
        Args:
            message: Message dict
            topic_suffix: Topic suffix for broadcast
        """
        topic = f"/atlas/{topic_suffix}"
        await self.pubsub.publish(topic, message)
    
    async def listen_for_node(self, node_id: str) -> AsyncIterator[Dict[str, Any]]:
        """
        Listen for messages directed to this node
        
        Args:
            node_id: This node's ID
            
        Yields:
            Received messages
        """
        topic = f"/atlas/node/{node_id}"
        async for message in self.pubsub.subscribe(topic):
            yield message

