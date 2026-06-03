import os

class Settings:
    MODEL_PATH = os.getenv(
        "MODEL_PATH", 
        "/home/rames/personal-project/noosphere/services/ai-inference/models/Llama-3.2-3B-Instruct-Q4_K_M.gguf"
    )
    CONTEXT_WINDOW: int = 2048
    GPU_LAYERS: int = 0

settings = Settings()