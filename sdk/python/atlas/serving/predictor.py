from typing import Any, Dict, List, Optional
from pathlib import Path


class Predictor:
    
    def __init__(self, model_path: str):
        self.model_path = model_path
        self.model: Optional[Any] = None
        self._load_model()
    
    def _load_model(self):
        path = Path(self.model_path)
        extension = path.suffix.lower()
        
        if extension in [".pth", ".pt"]:
            try:
                import torch
                self.model = torch.load(self.model_path, map_location="cpu")
            except ImportError:
                raise ImportError("PyTorch not installed. Install with: pip install torch")
        elif extension == ".h5":
            try:
                import tensorflow as tf
                self.model = tf.keras.models.load_model(self.model_path)
            except ImportError:
                raise ImportError("TensorFlow not installed. Install with: pip install tensorflow")
        elif extension == ".onnx":
            try:
                import onnxruntime as ort
                self.model = ort.InferenceSession(self.model_path)
            except ImportError:
                raise ImportError("ONNX Runtime not installed. Install with: pip install onnxruntime")
        else:
            raise ValueError(f"Unsupported model format: {extension}")
    
    def predict(self, input_data: Any, model_type: str = "auto", options: dict = None) -> Any:
        if self.model is None:
            raise RuntimeError("Model not loaded")
        
        if options is None:
            options = {}
        
        if model_type == "llm" and hasattr(self.model, "generate"):
            return self._predict_llm(input_data, options)
        elif model_type == "vision" or (model_type == "auto" and self._is_vision_model()):
            return self._predict_vision(input_data, options)
        elif model_type == "speech" or (model_type == "auto" and self._is_speech_model()):
            return self._predict_speech(input_data, options)
        elif model_type == "embedding" or (model_type == "auto" and self._is_embedding_model()):
            return self._predict_embedding(input_data, options)
        else:
            return self._predict_generic(input_data)
    
    def _predict_generic(self, input_data: Any) -> Any:
        if hasattr(self.model, "predict"):
            return self.model.predict(input_data)
        elif hasattr(self.model, "__call__"):
            return self.model(input_data)
        else:
            raise RuntimeError("Model does not support prediction")
    
    def _predict_llm(self, input_data: Any, options: dict) -> Any:
        try:
            import torch
        except ImportError:
            return self._predict_generic(input_data)
        
        if isinstance(self.model, torch.nn.Module):
            if isinstance(input_data, str):
                if hasattr(self.model, "tokenizer"):
                    tokens = self.model.tokenizer.encode(input_data, return_tensors="pt")
                else:
                    return self._predict_generic(input_data)
            else:
                tokens = input_data
            
            self.model.eval()
            with torch.no_grad():
                if hasattr(self.model, "generate"):
                    output = self.model.generate(
                        tokens,
                        max_length=options.get("max_length", 100),
                        temperature=options.get("temperature", 1.0),
                        do_sample=options.get("do_sample", True),
                    )
                else:
                    output = self.model(tokens)
            
            if hasattr(self.model, "tokenizer"):
                return self.model.tokenizer.decode(output[0])
            return output.tolist()
        return self._predict_generic(input_data)
    
    def _predict_vision(self, input_data: Any, options: dict) -> Any:
        try:
            import torch
            import numpy as np
        except ImportError:
            return self._predict_generic(input_data)
        
        try:
            from PIL import Image
        except ImportError:
            pass
        
        if isinstance(input_data, str):
            try:
                image = Image.open(input_data)
                if hasattr(self.model, "preprocess"):
                    input_tensor = self.model.preprocess(image)
                else:
                    try:
                        import torchvision.transforms as transforms
                        transform = transforms.Compose([
                            transforms.Resize((224, 224)),
                            transforms.ToTensor(),
                        ])
                        input_tensor = transform(image).unsqueeze(0)
                    except ImportError:
                        return self._predict_generic(input_data)
            except Exception:
                return self._predict_generic(input_data)
        elif isinstance(input_data, np.ndarray):
            input_tensor = torch.from_numpy(input_data).float()
        else:
            input_tensor = input_data
        
        if isinstance(self.model, torch.nn.Module):
            self.model.eval()
            with torch.no_grad():
                output = self.model(input_tensor)
            return output.tolist()
        return self._predict_generic(input_data)
    
    def _predict_speech(self, input_data: Any, options: dict) -> Any:
        return self._predict_generic(input_data)
    
    def _predict_embedding(self, input_data: Any, options: dict) -> Any:
        result = self._predict_generic(input_data)
        if isinstance(result, (list, tuple)) and len(result) > 0:
            if isinstance(result[0], (list, tuple)):
                return result[0]
        return result
    
    def _is_vision_model(self) -> bool:
        return "vision" in str(type(self.model)).lower() or "resnet" in str(type(self.model)).lower()
    
    def _is_speech_model(self) -> bool:
        return "whisper" in str(type(self.model)).lower() or "speech" in str(type(self.model)).lower()
    
    def _is_embedding_model(self) -> bool:
        return "embedding" in str(type(self.model)).lower() or "encoder" in str(type(self.model)).lower()

