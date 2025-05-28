# Sources:
# https://fastapi.tiangolo.com/tutorial/testing/
# https://www.python-httpx.org/api/

from sqlalchemy import text


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
