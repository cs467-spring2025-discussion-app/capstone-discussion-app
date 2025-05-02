from datetime import datetime
from pydantic import BaseModel
from typing import Optional
from uuid import UUID


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
    last_shown: Optional[datetime] = None
    is_active: bool
    show_count: int


class TopicList(BaseModel):
    topics: list[Topic]


################################################################################
# Discussion Reply Models
################################################################################

class NewReply(BaseModel):
    user_id: UUID | None
    content: str


class Reply(BaseModel):
    id: int
    discussion_id: int
    user_id: UUID | None
    content: str
    created_at: datetime
    is_deleted: bool
    deleted_at: datetime | None
    is_hidden: bool
    hidden_at: datetime | None


class ReplyList(BaseModel):
    replies: list[Reply]
    total_count: int

