from pydantic import BaseModel

class PromptRequest(BaseModel):
    prompt: str
    max_tokens: int = 512
    temperature: float = 0.7

class PromptResponse(BaseModel):
    response: str