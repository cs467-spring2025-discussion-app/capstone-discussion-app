from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy import create_engine, select, text
from sqlalchemy.orm import Session, sessionmaker
from models import Base, DiscussionTopics
from schemas import NewPrompt
import os

DATABASE_URL = os.getenv("DATABASE_URL")
if not DATABASE_URL:
    raise ValueError("DATABASE_URL environment variable is not set.")

engine = create_engine(DATABASE_URL, echo=True)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base.metadata.create_all(bind=engine)

app = FastAPI()


def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()


@app.get("/")
def read_root():
    try:
        with engine.connect() as conn:
            result = conn.execute(text("SELECT 'Hello from PostgreSQL!'"))
            return {"message": result.scalar()}
    except Exception as e:
        return {"error": str(e)}


@app.put("/submit_prompt")
def submit_prompt(new_prompt: NewPrompt, db: Session = Depends(get_db)):
    try:
        prompt = DiscussionTopics(
            title=new_prompt.title,
            prompt=new_prompt.prompt,
            is_active=False  # All new prompts are inactive by default
        )
        db.add(prompt)
        db.commit()
        db.refresh(prompt)
        return {"message": "Prompt submitted successfully.", "id": prompt.id}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/activate_prompt/{prompt_id}")
def activate_prompt(prompt_id: int, db: Session = Depends(get_db)):
    try:
        prompt = db.query(DiscussionTopics).filter(DiscussionTopics.id == prompt_id).first()
        if not prompt:
            raise HTTPException(status_code=404, detail="Prompt not found.")
        prompt.is_active = True
        db.commit()
        return {"message": f"Prompt with ID {prompt_id} activated successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/get_daily_prompt")
def get_daily_prompt(db: Session = Depends(get_db)):
    try:
        # Fetch the first active prompt
        stmt = select(DiscussionTopics).where(DiscussionTopics.is_active == True)
        daily_prompt = db.execute(stmt).scalars().first()
        if not daily_prompt:
            raise HTTPException(status_code=404, detail="No active prompt found.")
        return {
            "id": daily_prompt.id,
            "title": daily_prompt.title,
            "prompt": daily_prompt.prompt,
            "created_at": daily_prompt.created_at,
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.delete("/delete_prompt/{prompt_id}")
def delete_prompt(prompt_id: int, db: Session = Depends(get_db)):
    try:
        prompt = db.query(DiscussionTopics).filter(DiscussionTopics.id == prompt_id).first()
        if not prompt:
            raise HTTPException(status_code=404, detail="Prompt not found.")

        db.delete(prompt)
        db.commit()
        return {"message": f"Prompt with ID {prompt_id} deleted successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
