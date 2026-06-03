import os

class Settings:
    MODEL_PATH: str = os.getenv("MODEL_PATH", "/app/models/Llama-3.2-3B-Instruct-Q4_K_M.gguf")
    CONTEXT_WINDOW: int = 2048
    GPU_LAYERS: int = 0

settings = Settings()