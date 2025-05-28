# Sources:
# https://fastapi.tiangolo.com/tutorial/testing/
# https://fastapi.tiangolo.com/reference/testclient/
# https://www.python-httpx.org/api/
# https://fastapi.tiangolo.com/reference/response/
# https://docs.pydantic.dev/latest/concepts/validators/

from sqlalchemy import text
from models import Topic
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

        for id in range(1, 10):
            title, prompt = f"topic {id}", f"prompt {id}"
            response = helper_submit_topic(client, title, prompt)

            # Expect Status OK with properly formatted request
            assert response.status_code == 200
            assert response.json() == {
                "message": "Topic submitted successfully.",
                "id": id,
            }

            # Expect topic row with id i in the database
            topic = db_session.query(DiscussionTopics).filter_by(id=id).first()
            assert topic.title == title
            assert topic.prompt == prompt

    def test_get_topic(self, client):
        for id in range(1, 10):
            title, prompt = f"topic {id}", f"prompt {id}"
            helper_submit_topic(client, title, prompt)
            response = helper_get_topic(client, id)
            content = Topic.model_validate_json(response.content)
            assert response.status_code == 200
            assert content.id == id
            assert content.title == title
            assert content.prompt == prompt


def helper_submit_topic(client, title, prompt):
    return client.post(
        "/topics/submit_topic",
        json={"title": title, "prompt": prompt},
    )


def helper_get_topic(client, id):
    return client.get(
        f"/topics/get_topic/{id}",
    )
