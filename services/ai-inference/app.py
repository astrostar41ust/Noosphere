from fastapi import FastAPI, HTTPException
from .config import settings
from .schemas import PromptRequest, PromptResponse
from .engine import ai_engine

app = FastAPI(title="Local AI Inference Sidecar")

@app.on_event("startup")
def startup_event():
    try:
        ai_engine.load_model()
    except Exception as e:
        print(f"Startup Failed: {str(e)}")

@app.post("/v1/chat", response_model=PromptResponse)
async def chat_endpoint(request: PromptRequest):
    try:
        generated_text = ai_engine.generate(
            prompt=request.prompt,
            max_tokens=request.max_tokens,
            temperature=request.temperature
        )
        return PromptResponse(response=generated_text)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))