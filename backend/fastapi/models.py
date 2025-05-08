from datetime import datetime
from pydantic import BaseModel, ConfigDict
from typing import Optional
from uuid import UUID


################################################################################
# Discussion Topic Models
################################################################################

class NewTopic(BaseModel):
    title: str
    prompt: str


class Topic(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    id: int
    title: str
    prompt: str
    created_at: datetime
    last_shown: Optional[datetime] = None
    is_active: bool
    show_count: int


class TopicList(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    topics: list[Topic]


################################################################################
# Discussion Reply Models
################################################################################

class NewReply(BaseModel):
    user_id: UUID | None
    content: str


class Reply(BaseModel):
    model_config = ConfigDict(from_attributes=True)

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
    model_config = ConfigDict(from_attributes=True)

    replies: list[Reply]
    total_count: int


################################################################################
# Vote Models
################################################################################

class NewVote(BaseModel):
    user_id: UUID
    vote_type: int  # 1 for upvote, -1 for downvote, 0 to remove vote


class Vote(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    user_id: UUID
    reply_id: int
    vote_type: int
    created_at: datetime
    updated_at: datetime


class VoteList(BaseModel):
    model_config = ConfigDict(from_attributes=True)

    votes: list[Vote]
    total_count: int
