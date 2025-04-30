from datetime import datetime
from pydantic import BaseModel


class NewPrompt(BaseModel):
    title: str
    prompt: str


class Prompt(BaseModel):
    id: int
    title: str
    prompt: str
    created_at: datetime
    last_shown: datetime | None


class PromptList(BaseModel):
    prompts: list[Prompt]
