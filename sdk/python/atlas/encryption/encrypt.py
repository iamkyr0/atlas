"""
Encryption - Encrypt datasets before upload
"""

import os
import tempfile
from pathlib import Path
from cryptography.fernet import Fernet
from typing import Optional


class DatasetEncryption:
    """Encrypt datasets for private training"""
    
    def __init__(self, key: Optional[bytes] = None):
        """
        Initialize encryption
        
        Args:
            key: Encryption key (if None, generates new key)
        """
        if key is None:
            key = Fernet.generate_key()
        self.key = key
        self.cipher = Fernet(key)
    
    def get_key(self) -> bytes:
        """Get encryption key"""
        return self.key
    
    async def encrypt_file(self, file_path: str, output_path: Optional[str] = None) -> str:
        """
        Encrypt a file
        
        Args:
            file_path: Path to file to encrypt
            output_path: Output path (if None, creates temp file)
            
        Returns:
            Path to encrypted file
        """
        file_path_obj = Path(file_path)
        
        if not file_path_obj.exists():
            raise FileNotFoundError(f"File not found: {file_path}")
        
        # Read file
        with open(file_path, "rb") as f:
            data = f.read()
        
        # Encrypt
        encrypted_data = self.cipher.encrypt(data)
        
        # Write encrypted file
        if output_path is None:
            output_path = str(file_path_obj.parent / f"{file_path_obj.stem}.encrypted")
        
        with open(output_path, "wb") as f:
            f.write(encrypted_data)
        
        return output_path
    
    async def decrypt_file(self, encrypted_path: str, output_path: Optional[str] = None) -> str:
        """
        Decrypt a file
        
        Args:
            encrypted_path: Path to encrypted file
            output_path: Output path (if None, creates temp file)
            
        Returns:
            Path to decrypted file
        """
        encrypted_path_obj = Path(encrypted_path)
        
        if not encrypted_path_obj.exists():
            raise FileNotFoundError(f"File not found: {encrypted_path}")
        
        # Read encrypted file
        with open(encrypted_path, "rb") as f:
            encrypted_data = f.read()
        
        # Decrypt
        data = self.cipher.decrypt(encrypted_data)
        
        # Write decrypted file
        if output_path is None:
            output_path = str(encrypted_path_obj.parent / f"{encrypted_path_obj.stem}.decrypted")
        
        with open(output_path, "wb") as f:
            f.write(data)
        
        return output_path


# Convenience function
async def encrypt_file(file_path: str, key: Optional[bytes] = None) -> str:
    """
    Encrypt a file
    
    Args:
        file_path: Path to file
        key: Encryption key (if None, generates new key)
        
    Returns:
        Path to encrypted file
    """
    enc = DatasetEncryption(key=key)
    return await enc.encrypt_file(file_path)

