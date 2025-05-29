# Sources:
# https://fastapi.tiangolo.com/tutorial/testing/
# https://fastapi.tiangolo.com/reference/testclient/
# https://www.python-httpx.org/api/
# https://fastapi.tiangolo.com/reference/response/
# https://docs.pydantic.dev/latest/concepts/validators/
# https://docs.pytest.org/en/7.1.x/how-to/fixtures.html

import uuid  # usage note: auth service uses v4

import pytest
from models import (
    MessageIdResponse,
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

    def test_database_connection(self, db_session):
        result = db_session.execute(text("SELECT 1")).scalar()
        assert result == 1

    def test_home(self, client):
        response = client.get("/")
        assert response.status_code == 200
        assert response.json() == {"message": "Hello from PostgreSQL!"}


class TestDiscussionTopics:
    def test_submit_topic(self, db_session, fixture_topics):
        """
        Submitting a new topic successfully creates a topic in the database,
        returns 200 response with JSON.
        """
        topics = fixture_topics
        for topic_id, _ in enumerate(topics, start=1):
            # Check that topic actually exists in database
            topic = db_session.query(DiscussionTopics).filter_by(id=topic_id).first()
            assert topic.id == topic_id
            assert topic.title == f"test topic {topic_id}"
            assert topic.prompt == f"test prompt {topic_id}"

    def test_get_topic(self, client, fixture_topics):
        """
        Requesting a topic with a valid id returns 200 response with JSON of
        topic data.
        """
        topics = fixture_topics
        for _, topic in enumerate(topics, start=1):
            response = helper_get_topic(client, topic.id)
            content = Topic.model_validate_json(response.content)
            assert response.status_code == 200
            assert content.id == topic.id
            assert content.title == f"test topic {topic.id}"
            assert content.prompt == f"test prompt {topic.id}"

    def test_update_topic(self, client, fixture_topic):
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
            client, topic_id, new_title, new_prompt,
        )
        assert response.status_code == 200

        # Check for updated content in db via endpoint
        response = helper_get_topic(client, topic_id)
        content = Topic.model_validate_json(response.content)
        assert response.status_code == 200
        assert content.id == topic_id
        assert content.title == new_title
        assert content.prompt == new_prompt

    def test_delete_topic(self, client, fixture_topic):
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


################################################################################
# Topic helpers
################################################################################


@pytest.fixture
def fixture_topic(client):
    """Create a user for testing user endpoints"""
    response = helper_submit_topic(client, "test topic", "test prompt")
    assert response.status_code == 200
    return MessageIdResponse.model_validate_json(response.content)


@pytest.fixture
def fixture_topics(db_session, client):
    """Create a list of users for testing user endpoints"""
    topics = []
    for topic_id in range(1, 10):
        response = helper_submit_topic(
            client,
            f"test topic {topic_id}",
            f"test prompt {topic_id}",
        )
        assert response.status_code == 200
        content = MessageIdResponse.model_validate_json(response.content)
        assert content.id == topic_id
        topics.append(content)
    return topics


def helper_submit_topic(client, title, prompt):
    return client.post(
        "/topics/submit_topic",
        json={"title": title, "prompt": prompt},
    )


def helper_get_topic(client, id):
    return client.get(
        f"/topics/get_topic/{id}",
    )


def helper_update_topic(client, id, title, prompt):
    return client.put(
        f"/topics/update_topic/{id}",
        json={
            "title": title,
            "prompt": prompt,
        },
    )


def helper_delete_topic(client, id):
    return client.delete(
        f"/topics/delete_topic/{id}",
    )


class TestReplies:
    def test_submit_reply(self, db_session, fixture_replies):
        # Register user for a valid user id, submit topic for valid topic id
        # (note, can't use fixture in a loop)
        replies = fixture_replies
        for reply_id, _ in enumerate(replies, start=1):
            # Check that reply actually exists in database
            reply = db_session.query(Replies).filter_by(id=reply_id).first()
            assert reply.id == reply_id
            assert reply.content == f"test content {reply_id}"


################################################################################
# Reply helpers
################################################################################


@pytest.fixture
def fixture_replies(db_session, client):
    """Create a user and topic for testing reply end points"""
    user = helper_register_user(db_session, uuid.uuid4(), "test user", "test password")
    response = helper_submit_topic(client, "test topic", "test prompt")
    topic = MessageIdResponse.model_validate_json(response.content)
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
        content = MessageIdResponse.model_validate_json(response.content)
        assert content.id == reply_id
        replies.append(content)
    return replies


def helper_submit_reply(client, user_id, topic_id, content):
    return client.post(
        f"/replies/submit_reply/{topic_id}",
        json={"user_id": user_id, "content": content},
    )


def helper_register_user(db_session, id, email, password):
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
    user_id = result.id
    assert isinstance(user_id, uuid.UUID)
    return result
