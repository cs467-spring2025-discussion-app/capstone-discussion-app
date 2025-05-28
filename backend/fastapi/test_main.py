# Sources:
# https://fastapi.tiangolo.com/tutorial/testing/
# https://fastapi.tiangolo.com/reference/testclient/
# https://www.python-httpx.org/api/
# https://fastapi.tiangolo.com/reference/response/
# https://docs.pydantic.dev/latest/concepts/validators/

from sqlalchemy import text
from models import (
    MessageIdResponse,
    Topic,
)
from tables import (
    DiscussionTopics,
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
    def test_submit_topic(self, client, db_session):
        """
        Submitting a new topic successfully creates a topic in the database,
        returns 200 response with JSON.
        """

        for topic_id in range(1, 10):
            # Submit initial topics, check for matching ids
            title, prompt = f"submit topic {topic_id}", f"submit prompt {topic_id}"
            response = helper_submit_topic(client, title, prompt)
            content = MessageIdResponse.model_validate_json(response.content)
            assert content.id == topic_id

            # Expect topic row with id i in the database
            topic = db_session.query(DiscussionTopics).filter_by(id=topic_id).first()
            assert topic.title == title
            assert topic.prompt == prompt

    def test_get_topic(self, client):
        """
        Requesting a topic with a valid id returns 200 response with JSON of
        topic data.
        """
        for topic_id in range(1, 10):
            title, prompt = f"get topic {topic_id}", f"get prompt {topic_id}"

            # Submit initial topics, check for matching ids
            response = helper_submit_topic(client, title, prompt)
            assert response.status_code == 200
            content = MessageIdResponse.model_validate_json(response.content)
            assert content.id == topic_id

            # Get new topics
            response = helper_get_topic(client, topic_id)
            content = Topic.model_validate_json(response.content)
            assert response.status_code == 200
            assert content.id == topic_id
            assert content.title == title
            assert content.prompt == prompt

    def test_update_topic(self, client):
        """
        Updating a topic with a valid id returns 200 response with JSON of
        topic data. Topic content should be updated in the database.
        """
        topic_id = 1
        title, prompt = f"update topic {topic_id}", f"update prompt {topic_id}"

        # Submit initial topic, check for matching ids
        response = helper_submit_topic(client, title, prompt)
        assert response.status_code == 200
        content = MessageIdResponse.model_validate_json(response.content)
        assert content.id == topic_id

        # Submit topic updates
        response = helper_update_topic(
            client, topic_id, f"updated {title}", f"updated {prompt}"
        )
        assert response.status_code == 200

        # Check for updated content in db via endpoint
        response = helper_get_topic(client, topic_id)
        content = Topic.model_validate_json(response.content)
        assert response.status_code == 200
        assert content.id == topic_id
        assert content.title == f"updated {title}"
        assert content.prompt == f"updated {prompt}"

    def test_delete_topic(self, client):
        """
        Deleting a topic with a valid id returns 200 response with JSON of
        topic data. Topic should no longer be in the database.
        """
        topic_id = 1
        title, prompt = f"delete topic {topic_id}", f"delete prompt {topic_id}"

        # Submit initial topic, check for matching ids
        response = helper_submit_topic(client, title, prompt)
        assert response.status_code == 200
        content = MessageIdResponse.model_validate_json(response.content)
        assert content.id == topic_id

        # Delete topic
        response = helper_delete_topic(client, topic_id)
        assert response.status_code == 200
        # Check for deleted content in db via endpoint
        response = helper_get_topic(client, topic_id)
        assert response.status_code == 500


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
