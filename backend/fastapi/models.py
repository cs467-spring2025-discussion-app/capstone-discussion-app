from datetime import datetime
from pydantic import BaseModel


################################################################################
# Discussion Topic Models
################################################################################

class NewTopic(BaseModel):
    title: str
    prompt: str


class Topic(BaseModel):
    id: int
    title: str
    prompt: str
    created_at: datetime
    last_shown: datetime | None


class TopicList(BaseModel):
    prompts: list[Topic]


################################################################################
# Discussion Reply Models
################################################################################

class NewReply(BaseModel):
    discussion_id: int
    user_id: str | None
    content: str


class Reply(BaseModel):
    id: int
    discussion_id: int
    user_id: str | None
    content: str
    created_at: datetime
    is_deleted: bool
    deleted_at: datetime | None
    is_hidden: bool
    hidden_at: datetime | None


class ReplyList(BaseModel):
    replies: list[Reply]
    total_count: int

