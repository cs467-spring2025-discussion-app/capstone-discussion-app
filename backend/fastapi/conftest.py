# Source:
# https://docs.pytest.org/en/7.1.x/reference/fixtures.html

import pytest
from fastapi.testclient import TestClient
from main import TEST_DB_URL, app
from sqlalchemy import create_engine, text
from sqlalchemy.orm import sessionmaker
from tables import Base


@pytest.fixture(scope="session", autouse=True)
def setup_test_db():
    """Setup and teardown test database schema"""
    engine = create_engine(TEST_DB_URL)

    # Setup (reset tables)
    with engine.begin() as conn:
        conn.execute(text("DROP SCHEMA public CASCADE"))
        conn.execute(text("CREATE SCHEMA public"))
        conn.execute(text("GRANT ALL ON SCHEMA public TO public"))

    Base.metadata.create_all(bind=engine)

    yield

    # Teardown
    with engine.begin() as conn:
        conn.execute(text("DROP SCHEMA public CASCADE"))
        conn.execute(text("CREATE SCHEMA public"))
        conn.execute(text("GRANT ALL ON SCHEMA public TO public"))

    engine.dispose()


@pytest.fixture(scope="function")
def db_session():
    """Set up test session"""
    engine = create_engine(TEST_DB_URL)
    SessionLocal = sessionmaker(bind=engine)
    session = SessionLocal()

    try:
        yield session
    finally:
        session.close()


@pytest.fixture(scope="function", autouse=True)
def clean_db():
    """Empty tables and reset ids between tests"""
    engine = create_engine(TEST_DB_URL)

    yield

    with engine.begin() as conn:
        conn.execute(
            text("""
            truncate table
                votes,
                user_badges,
                user_reputation,
                replies,
                discussion_topics
            restart identity cascade
        """)
        )

    engine.dispose()


@pytest.fixture(scope="function")
def client():
    with TestClient(app) as test_client:
        yield test_client
