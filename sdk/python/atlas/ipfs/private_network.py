import os
from typing import List, Optional
from ipfshttpclient import Client as IPFSClient


class PrivateNetwork:
    
    def __init__(self, ipfs_api_url: str = "/ip4/127.0.0.1/tcp/5001"):
        self.ipfs_client = IPFSClient(ipfs_api_url)
    
    def configure_private_network(self, swarm_key: str, bootstrap_peers: List[str]):
        for peer in bootstrap_peers:
            try:
                self.ipfs_client.bootstrap.add(peer)
            except Exception as e:
                print(f"Failed to add bootstrap peer {peer}: {e}")
    
    def remove_public_peers(self):
        try:
            peers = self.ipfs_client.bootstrap.list()
            default_peers = [
                "/dnsaddr/bootstrap.libp2p.io/p2p/QmNnooDu7bfjPFoTZYxMNLWUQJyrVwtbZg5gBMjTezGAJN",
                "/dnsaddr/bootstrap.libp2p.io/p2p/QmbLHAnMoJPWSCR5Zp4Hz9HqHxNVSVX3qXqTfzJz3rP7xT",
            ]
            
            for peer in default_peers:
                try:
                    self.ipfs_client.bootstrap.rm(peer)
                except Exception:
                    pass
        except Exception as e:
            print(f"Failed to remove public peers: {e}")
    
    def add_peer(self, peer_multiaddr: str):
        try:
            self.ipfs_client.bootstrap.add(peer_multiaddr)
        except Exception as e:
            print(f"Failed to add peer: {e}")
    
    def list_peers(self) -> List[str]:
        try:
            peers = self.ipfs_client.swarm.peers()
            return [peer.get("addr", "") for peer in peers]
        except Exception as e:
            print(f"Failed to list peers: {e}")
            return []


class BootstrapNode:
    
    def __init__(self, ipfs_api_url: str = "/ip4/127.0.0.1/tcp/5001"):
        self.ipfs_client = IPFSClient(ipfs_api_url)
    
    def get_bootstrap_address(self) -> str:
        try:
            node_id = self.ipfs_client.id()["ID"]
            addresses = self.ipfs_client.id()["Addresses"]
            if addresses:
                return addresses[0]
            return f"/ip4/127.0.0.1/p2p/{node_id}"
        except Exception as e:
            print(f"Failed to get bootstrap address: {e}")
            return ""

