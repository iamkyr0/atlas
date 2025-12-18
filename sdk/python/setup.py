"""
Setup script for Atlas Python SDK
"""

from setuptools import setup, find_packages

with open("README.md", "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name="atlas-sdk",
    version="0.1.0",
    author="Atlas Team",
    description="Atlas Decentralized AI Platform SDK",
    long_description=long_description,
    long_description_content_type="text/markdown",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
    ],
    python_requires=">=3.10",
    install_requires=[
        "requests>=2.31.0",
        "ipfshttpclient>=0.8.0",
        "pydantic>=2.0.0",
        "python-dotenv>=1.0.0",
        "aiohttp>=3.9.0",
        "websockets>=12.0",
        "grpcio>=1.60.0",
        "grpcio-tools>=1.60.0",
        "cryptography>=41.0.0",
        "protobuf>=4.25.0",
    ],
    entry_points={
        "console_scripts": [
            "atlas=atlas.cli.cli:main",
        ],
    },
)

