from pydantic import BaseModel


class NewPrompt(BaseModel):
    title: str
    prompt: str
