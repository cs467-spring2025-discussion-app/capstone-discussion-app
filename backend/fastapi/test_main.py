# Sources:
# https://fastapi.tiangolo.com/tutorial/testing/
# https://fastapi.tiangolo.com/reference/testclient/
# https://www.python-httpx.org/api/
# https://fastapi.tiangolo.com/reference/response/
# https://docs.pydantic.dev/latest/concepts/validators/
# https://docs.pytest.org/en/7.1.x/how-to/fixtures.html

import uuid  # usage note: auth service uses v4
from typing import List
import json
import pytest
from fastapi.testclient import TestClient
from httpx import Response
from sqlalchemy.engine import Row
from sqlalchemy.orm import Session
from models import (
    MessageIdResponse,
    Reply,
    ReplyList,
    Topic,
)
from sqlalchemy import text
from tables import (
    DiscussionTopics,
    Replies,
)


class TestFastApi:
    """
    Test database setup sanity checks
    """

    def test_database_connection(self, db_session: Session) -> None:
        result = db_session.execute(text("SELECT 1")).scalar()
        assert result == 1

    def test_home(self, client: TestClient) -> None:
        response = client.get("/")
        assert response.status_code == 200
        assert response.json() == {"message": "Hello from PostgreSQL!"}


class TestDiscussionTopics:
    def test_submit_topic(
        self, db_session: Session, fixture_topics: List[MessageIdResponse]
    ) -> None:
        """
        Submitting a new topic successfully creates a topic in the database,
        returns 200 response with JSON.
        """
        topics = fixture_topics
        for topic_id, _ in enumerate(topics, start=1):
            topic = db_session.query(DiscussionTopics).filter_by(id=topic_id).first()
            assert topic is not None
            assert topic.id == topic_id
            assert topic.title == f"test topic {topic_id}"
            assert topic.prompt == f"test prompt {topic_id}"

    def test_get_topic(
        self, client: TestClient, fixture_topics: List[MessageIdResponse]
    ) -> None:
        """
        Requesting a topic with a valid id returns 200 response with JSON of
        topic data.
        """
        topics = fixture_topics
        for _, topic in enumerate(topics, start=1):
            response = helper_get_topic(client, topic.id)
            content = Topic.model_validate(response.json())
            assert response.status_code == 200
            assert content.id == topic.id
            assert content.title == f"test topic {topic.id}"
            assert content.prompt == f"test prompt {topic.id}"

    def test_update_topic(
        self, client: TestClient, fixture_topic: MessageIdResponse
    ) -> None:
        """
        Updating a topic with a valid id returns 200 response with JSON of
        topic data. Topic content should be updated in the database.
        """
        topic_id = 1
        new_title, new_prompt = f"update topic {topic_id}", f"update prompt {topic_id}"

        # Submit initial topic, check for matching ids
        topic = fixture_topic
        assert topic.id == topic_id

        # Submit topic updates
        response = helper_update_topic(
            client,
            topic_id,
            new_title,
            new_prompt,
        )
        assert response.status_code == 200

        # Check for updated content in db via endpoint
        response = helper_get_topic(client, topic_id)
        content = Topic.model_validate(response.json())
        assert response.status_code == 200
        assert content.id == topic_id
        assert content.title == new_title
        assert content.prompt == new_prompt

    def test_delete_topic(
        self, client: TestClient, fixture_topic: MessageIdResponse
    ) -> None:
        """
        Deleting a topic with a valid id returns 200 response with JSON of
        topic data. Topic should no longer be in the database.
        """
        topic_id = 1

        # Submit initial topic, check for matching ids
        topic = fixture_topic
        assert topic.id == topic_id

        # Delete topic
        response = helper_delete_topic(client, topic_id)
        assert response.status_code == 200
        # Check for deleted content in db via endpoint
        response = helper_get_topic(client, topic_id)
        assert response.status_code == 500

    def test_get_user_reply_count_by_topic(
        self, client: TestClient, db_session, fixture_topic
    ):
        topic = fixture_topic
        # Register test user
        user = helper_register_user(
            db_session, uuid.uuid4(), "test user", "test password"
        )

        total_replies = 10
        for reply_id in range(1, total_replies + 1):  # +1 for 10 total replies
            helper_submit_reply(client, str(user.id), topic.id, f"content {reply_id}")

        response = helper_get_user_reply_counts(
            client,
            str(user.id),
            topic.id,
        )
        assert response.status_code == 200
        content = json.loads(response.content)
        assert content["reply_count"] == total_replies


################################################################################
# Topic helpers
################################################################################


@pytest.fixture
def fixture_topic(client: TestClient) -> MessageIdResponse:
    """Create a topic for testing topic endpoints"""
    response = helper_submit_topic(client, "test topic", "test prompt")
    assert response.status_code == 200
    return MessageIdResponse.model_validate(response.json())


@pytest.fixture
def fixture_topics(client: TestClient) -> List[MessageIdResponse]:
    """Create a list of topics for testing topic endpoints"""
    topics = []
    for topic_id in range(1, 10):
        response = helper_submit_topic(
            client,
            f"test topic {topic_id}",
            f"test prompt {topic_id}",
        )
        assert response.status_code == 200
        content = MessageIdResponse.model_validate(response.json())
        assert content.id == topic_id
        topics.append(content)
    return topics


def helper_submit_topic(client: TestClient, title: str, prompt: str) -> Response:
    return client.post(
        "/topics/submit_topic",
        json={"title": title, "prompt": prompt},
    )


def helper_get_topic(client: TestClient, id: int) -> Response:
    return client.get(
        f"/topics/get_topic/{id}",
    )


def helper_update_topic(
    client: TestClient, id: int, title: str, prompt: str
) -> Response:
    return client.put(
        f"/topics/update_topic/{id}",
        json={
            "title": title,
            "prompt": prompt,
        },
    )


def helper_delete_topic(client: TestClient, id: int) -> Response:
    return client.delete(
        f"/topics/delete_topic/{id}",
    )


def helper_get_user_reply_counts(
    client: TestClient, user_id: str, topic_id: int
) -> Response:
    return client.post(
        f"/topics/user_reply_counts/{topic_id}",
        json={
            "user_id": user_id,
        },
    )


class TestReplies:
    def test_submit_reply(
        self, db_session: Session, fixture_replies: List[MessageIdResponse]
    ) -> None:
        replies = fixture_replies
        for reply_id, _ in enumerate(replies, start=1):
            reply = db_session.query(Replies).filter_by(id=reply_id).first()
            assert reply is not None
            assert reply.id == reply_id
            assert reply.content == f"test content {reply_id}"

    def test_get_reply(self, client, fixture_replies):
        replies = fixture_replies
        for reply_id, _ in enumerate(replies, start=1):
            response = helper_get_reply(client, reply_id)
            content = Reply.model_validate(response.json())
            assert response.status_code == 200
            assert content.id == reply_id
            assert content.content == f"test content {reply_id}"

    def test_get_all_replies(self, client, fixture_replies):
        replies = fixture_replies
        response = helper_get_all_replies(client, 1)
        content = ReplyList.model_validate(response.json())
        assert response.status_code == 200
        assert content.total_count == len(replies)
        for id, reply in enumerate(content.replies):
            assert reply.id == replies[id].id
            assert reply.content == f"test content {replies[id].id}"


################################################################################
# Reply helpers
################################################################################


@pytest.fixture
def fixture_replies(db_session: Session, client: TestClient) -> List[MessageIdResponse]:
    """Create a user and topic for testing reply end points"""
    user = helper_register_user(db_session, uuid.uuid4(), "test user", "test password")
    response = helper_submit_topic(client, "test topic", "test prompt")
    topic = MessageIdResponse.model_validate(response.json())
    replies = []
    for reply_id in range(1, 10):
        content = f"test content {reply_id}"
        response = helper_submit_reply(
            client,
            str(user.id),
            topic.id,
            content,
        )
        assert response.status_code == 200
        reply_response = MessageIdResponse.model_validate(response.json())
        assert reply_response.id == reply_id
        replies.append(reply_response)
    return replies


def helper_submit_reply(
    client: TestClient, user_id: str, topic_id: int, content: str
) -> Response:
    return client.post(
        f"/replies/submit_reply/{topic_id}",
        json={"user_id": user_id, "content": content},
    )


def helper_get_reply(client: TestClient, reply_id: int) -> Response:
    return client.get(
        f"/replies/get_reply/{reply_id}",
    )


def helper_get_all_replies(client: TestClient, discussion_id: int) -> Response:
    return client.get(f"/replies/get_all_replies/{discussion_id}")


def helper_register_user(
    db_session: Session, id: uuid.UUID, email: str, password: str
) -> Row:
    """Register a user in order to get a valid user ID for making replies and
    votes"""
    data = {"id": id, "email": email, "password": password}
    statement = text("""
        insert into users
        (id, email, password)
        values
        (:id, :email, :password)
        """)
    db_session.execute(statement, data)
    db_session.commit()
    statement = text("""
        select id from users where id = :id
    """)
    result = db_session.execute(statement, data).first()
    assert result is not None
    user_id = result.id
    assert isinstance(user_id, uuid.UUID)
    return result
