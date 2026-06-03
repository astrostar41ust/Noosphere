import uvicorn
from contextlib import asynccontextmanager
from fastapi import FastAPI, HTTPException
from config import settings
from schemas import PromptRequest, PromptResponse
from engine import ai_engine

@asynccontextmanager
async def lifespan(app: FastAPI):
    try:
        ai_engine.load_model()
    except Exception as e:
        print(f"Startup Failed: {str(e)}")
    yield

app = FastAPI(title="Local AI Inference Sidecar", lifespan=lifespan)

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

if __name__ == "__main__":
    uvicorn.run("app:app", host="0.0.0.0", port=8000, reload=True)