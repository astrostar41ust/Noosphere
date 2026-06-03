import os
from llama_cpp import Llama
from .config import settings

class InferenceEngine:
    def __init__(self):
        self.llm = None

    def load_model(self):
        if not os.path.exists(settings.MODEL_PATH):
            raise FileNotFoundError(f"Model file missing at {settings.MODEL_PATH}")
        
        print(f"Loading model into memory: {settings.MODEL_PATH}")
        self.llm = Llama(
            model_path=settings.MODEL_PATH,
            n_ctx=settings.CONTEXT_WINDOW,
            n_gpu_layers=settings.GPU_LAYERS
        )
        print("Model loaded successfully!")

    def generate(self, prompt: str, max_tokens: int, temperature: float) -> str:
        if self.llm is None:
            raise RuntimeError("Model is not loaded.")
        
        output = self.llm(
            prompt,
            max_tokens=max_tokens,
            temperature=temperature,
            stop=["</s>", "User:", "Assistant:"]
        )
        return output["choices"][0]["text"]

# Instantiate a single global engine instance
ai_engine = InferenceEngine()