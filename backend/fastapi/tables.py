from sqlalchemy import Column, String, Boolean, Integer, TIMESTAMP, text, PrimaryKeyConstraint, ForeignKeyConstraint, \
    Index, CheckConstraint
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import declarative_base
import uuid

Base = declarative_base()


class Users(Base):
    __tablename__ = "users"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    email = Column(String(255), unique=True, nullable=False, index=True)
    password = Column(String, nullable=False)
    last_login = Column(TIMESTAMP, nullable=True)
    failed_login_attempts = Column(Integer, default=0)
    account_locked = Column(Boolean, default=False)
    account_locked_until = Column(TIMESTAMP, nullable=True)

    __table_args__ = (
        Index("idx_users_email", "email"),
    )


class Sessions(Base):
    __tablename__ = "sessions"

    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    user_id = Column(UUID(as_uuid=True), nullable=False, index=True)
    created_at = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)
    expires_at = Column(TIMESTAMP, nullable=False)

    __table_args__ = (
        ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        Index("idx_sessions_user_id", "user_id"),
    )


class Tags(Base):
    __tablename__ = "tags"

    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String(50), unique=True, nullable=False)
    description = Column(String(50), nullable=False)  # What does this tag indicate?
    is_custom_tag = Column(Boolean, nullable=False, default=False)  # true for user-created tags
    created_at = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)


# Junction table
class UserTags(Base):
    __tablename__ = "user_tags"

    user_id = Column(UUID(as_uuid=True), nullable=False)
    tag_id = Column(Integer, nullable=False)
    __table_args__ = (
        PrimaryKeyConstraint("user_id", "tag_id"),
        ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        ForeignKeyConstraint(["tag_id"], ["tags.id"]),
    )


class Badges(Base):
    __tablename__ = "badges"

    id = Column(Integer, primary_key=True, autoincrement=True)
    name = Column(String(50), unique=True, nullable=False)
    description = Column(String(50), nullable=False)  # What does this badge indicate?
    created_at = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)


# Junction table
class UserBadges(Base):
    __tablename__ = "user_badges"

    user_id = Column(UUID(as_uuid=True), nullable=False)
    badge_id = Column(Integer, nullable=False)
    __table_args__ = (
        PrimaryKeyConstraint("user_id", "badge_id"),
        ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        ForeignKeyConstraint(["badge_id"], ["badges.id"]),
    )


class UserReputation(Base):
    __tablename__ = "user_reputation"

    user_id = Column(UUID(as_uuid=True), primary_key=True, nullable=False)
    upvotes_received = Column(Integer, nullable=False, default=0)
    downvotes_received = Column(Integer, nullable=False, default=0)
    last_updated = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)
    __table_args__ = (
        ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
    )


class DiscussionTopics(Base):
    __tablename__ = "discussion_topics"

    id = Column(Integer, primary_key=True, autoincrement=True)
    title = Column(String(255), nullable=False)
    prompt = Column(String, nullable=False)
    created_at = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)
    last_shown = Column(TIMESTAMP, nullable=True)
    is_active = Column(Boolean, default=False)
    show_count = Column(Integer, default=0)
    activated_at = Column(TIMESTAMP, nullable=True)

    __table_args__ = (
        Index("idx_discussion_topics_created_at", "created_at"),
        Index("idx_discussion_topics_is_active", "is_active"),
    )


class Replies(Base):
    __tablename__ = "replies"

    id = Column(Integer, primary_key=True, autoincrement=True)
    discussion_id = Column(Integer, nullable=False, index=True)
    user_id = Column(UUID(as_uuid=True), nullable=False, index=True)
    content = Column(String, nullable=False)
    created_at = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)
    is_deleted = Column(Boolean, default=False)
    deleted_at = Column(TIMESTAMP, nullable=True)
    is_hidden = Column(Boolean, default=False)
    hidden_at = Column(TIMESTAMP, nullable=True)

    __table_args__ = (
        ForeignKeyConstraint(["discussion_id"], ["discussion_topics.id"], ondelete="CASCADE"),
        ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        Index("idx_replies_discussion_id", "discussion_id"),
        Index("idx_replies_user_id", "user_id"),
    )


class Votes(Base):
    __tablename__ = "votes"

    user_id = Column(UUID(as_uuid=True), nullable=False)
    reply_id = Column(Integer, nullable=False)
    # Positive int for upvote (e.g. 1), negative for downvote (e.g. -1), 0 if a user cancels their up/downvote
    vote_type = Column(Integer, nullable=False)
    created_at = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)
    updated_at = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)

    __table_args__ = (
        PrimaryKeyConstraint("user_id", "reply_id"),
        ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        ForeignKeyConstraint(["reply_id"], ["replies.id"], ondelete="CASCADE"),
        CheckConstraint("vote_type in (-1, 0, 1)", name="check_vote_type"),
        # prevent duplicate votes on the same reply
        Index("idx_votes_reply_id", "reply_id"),
        Index("idx_votes_user_id", "user_id"),
    )


class FlaggedReplies(Base):
    __tablename__ = "flagged_replies"

    id = Column(Integer, primary_key=True, autoincrement=True)
    reply_id = Column(Integer, nullable=False)
    user_id = Column(UUID(as_uuid=True), nullable=False)
    reason = Column(String, nullable=True)
    flagged_at = Column(TIMESTAMP, server_default=text("CURRENT_TIMESTAMP"), nullable=False)

    __tableargs__ = (
        ForeignKeyConstraint(["reply_id"], ["replies.id"], ondelete="CASCADE"),
        ForeignKeyConstraint(["user_id"], ["users.id"], ondelete="CASCADE"),
        Index("idx_flagged_replies_reply_id", "reply_id"),
        Index("idx_flagged_replies_user_id", "user_id"),
    )
