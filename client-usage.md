# Client Usage Guide

Guide for using Atlas as a client to fine-tune and serve AI models.

## Installation

```bash
cd sdk/python
pip install -r requirements.txt
```

## Quick Start

### 1. Initialize Client

```python
from atlas import AtlasClient

client = AtlasClient(
    ipfs_api_url="/ip4/127.0.0.1/tcp/5001",
    chain_grpc_url="localhost:9090",
    creator="your_wallet_address"
)
```

### 2. Upload Dataset

```python
async with client:
    dataset_cid = await client.upload_dataset("path/to/dataset.zip")
    print(f"Dataset CID: {dataset_cid}")
```

### 3. Submit Training Job

```python
job = await client.submit_job(
    model_id="model-123",
    dataset_cid=dataset_cid,
    config={
        "epochs": 10,
        "batch_size": 32,
        "learning_rate": 0.001,
        "lora_rank": 8
    }
)
print(f"Job ID: {job.id}")
```

### 4. Monitor Job

```python
async for update in client.subscribe_to_job_updates(job.id):
    print(f"Status: {update.get('status')}, Progress: {update.get('progress')}")
    if update.get("status") == "completed":
        break
```

### 5. Download Model

```python
job = await client.get_job(job.id)
if job.status == "completed":
    await client.download_model(
        model_cid=job.model_cid,
        output_path="./trained_model.pt"
    )
```

## CLI Usage

### Submit Job

```bash
atlas submit-job model-123 QmXXX... --config config.json
```

### List Jobs

```bash
atlas list-jobs
```

### Get Job Status

```bash
atlas get-job job-123
```

### Upload Dataset

```bash
atlas upload-dataset dataset.zip --encrypt
```

### Download Model

```bash
atlas download-model model-123 --output ./models/
```

### Serve Model

```bash
atlas serve-model model-123 --port 8000
```

## Daemon Mode

Start daemon for REST API integration:

```bash
atlas daemon --port 8080
```

Daemon provides REST API at `http://localhost:8080/api/v1`:
- `POST /jobs`: Submit job
- `GET /jobs`: List jobs
- `GET /jobs/{id}`: Get job status
- `POST /models/register`: Register model
- `POST /models/{id}/predict`: Model inference

## Advanced Features

### Private Datasets

```python
dataset_cid = await client.upload_dataset(
    "private_dataset.zip",
    encrypt=True,
    private_network=True
)
```

### Custom Training Script

```python
script_cid = await client.upload_dataset("custom_train.py")
job = await client.submit_job(
    model_id="model-123",
    dataset_cid=dataset_cid,
    config={
        "training_script_cid": script_cid,
        "epochs": 10
    }
)
```

### Federated Learning

```python
job = await client.submit_job(
    model_id="model-123",
    dataset_cid=dataset_cid,
    config={
        "training_type": "federated",
        "min_clients": 5,
        "rounds": 10
    }
)
```

## Configuration

### Environment Variables

- `ATLAS_IPFS_API`: IPFS API URL
- `ATLAS_CHAIN_GRPC`: Chain gRPC URL
- `ATLAS_CREATOR`: Creator wallet address

### Config File

```json
{
  "epochs": 10,
  "batch_size": 32,
  "learning_rate": 0.001,
  "lora_rank": 8,
  "lora_alpha": 16
}
```

## Troubleshooting

### Connection Issues

- Check IPFS node is running: `ipfs id`
- Verify chain gRPC URL is correct
- Check network connectivity

### Job Failures

- Check job logs via blockchain query
- Verify dataset CID is accessible
- Ensure model ID exists

### Model Download

- Verify job status is "completed"
- Check model CID is valid
- Ensure sufficient disk space

