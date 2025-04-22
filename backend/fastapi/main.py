from fastapi import FastAPI
from sqlalchemy import create_engine, text
import os

DATABASE_URL = os.getenv("DATABASE_URL")
engine = create_engine(DATABASE_URL, echo=True)

app = FastAPI()


@app.get("/")
def read_root():
    try:
        with engine.connect() as conn:
            result = conn.execute(text("SELECT 'Hello from PostgreSQL!'"))
            return {"message": result.scalar()}
    except Exception as e:
        return {"error": str(e)}
