from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy import create_engine, select, text
from sqlalchemy.orm import Session, sessionmaker
from tables import Base, DiscussionTopics, Replies
from models import NewTopic, Topic, TopicList, NewReply, Reply, ReplyList
from datetime import datetime
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


################################################################################
# Discussion Topics Endpoints
################################################################################


@app.post("/topics/submit_topic", response_model=Topic)
def submit_topic(new_topic: NewTopic, db: Session = Depends(get_db)):
    try:
        topic = DiscussionTopics(
            title=new_topic.title,
            prompt=new_topic.prompt,
            is_active=False  # All new topics are inactive by default
        )
        db.add(topic)
        db.commit()
        db.refresh(topic)
        return {"message": "Topic submitted successfully.", "id": topic.id}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/topics/get_topics", response_model=TopicList)
def get_topics(db: Session = Depends(get_db)):
    try:
        # Fetch five random active topics
        stmt = select(DiscussionTopics).where(DiscussionTopics.is_active == True)
        daily_topics = db.execute(stmt).scalars().fetchmany(5)
        if not daily_topics:
            raise HTTPException(status_code=404, detail="No active topic found.")
        if len(daily_topics) < 5:
            raise HTTPException(status_code=404, detail="Less than 5 active topics found.")
        return {
            "topics": [
                {
                    "id": topic.id,
                    "title": topic.title,
                    "prompt": topic.prompt,
                    "created_at": topic.created_at.isoformat(),
                }
                for topic in daily_topics
            ]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/topics/get_topic/{topic_id}", response_model=Topic)
def get_topic(topic_id: int, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")
        return {
            "id": topic.id,
            "title": topic.title,
            "prompt": topic.prompt,
            "created_at": topic.created_at.isoformat(),
            "last_shown": topic.last_shown.isoformat() if topic.last_shown else None,
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/topics/get_all_topics", response_model=TopicList)
def get_all_topics(db: Session = Depends(get_db)):
    try:
        topics = db.query(DiscussionTopics).all()
        if not topics:
            raise HTTPException(status_code=404, detail="No topics found.")
        return {
            "topics": [
                {
                    "id": topic.id,
                    "title": topic.title,
                    "prompt": topic.prompt,
                    "created_at": topic.created_at.isoformat(),
                    "last_shown": topic.last_shown.isoformat() if topic.last_shown else None,
                }
                for topic in topics
            ]
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/topics/update_topic/{topic_id}")
def update_topic(topic_id: int, updated_topic: NewTopic, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")

        topic.title = updated_topic.title
        topic.prompt = updated_topic.prompt
        db.commit()
        return {"message": f"Topic with ID {topic_id} updated successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/topics/activate_topic/{topic_id}")
def activate_topic(topic_id: int, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")
        topic.is_active = True
        db.commit()
        return {"message": f"Topic with ID {topic_id} activated successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/topics/deactivate_topic/{topic_id}")
def deactivate_topic(topic_id: int, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")
        topic.is_active = False
        db.commit()
        return {"message": f"Topic with ID {topic_id} deactivated successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/topics/show_topic/{topic_id}")
def show_topic(topic_id: int, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")
        topic.show_count += 1
        topic.last_shown = datetime.now()
        db.commit()
        return {"message": f"Topic with ID {topic_id} shown successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.delete("/topics/delete_topic/{topic_id}")
def delete_topic(topic_id: int, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")

        db.delete(topic)
        db.commit()
        return {"message": f"Topic with ID {topic_id} deleted successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


################################################################################
# Discussion Replies Endpoints
################################################################################


@app.post("/replies/submit_reply/{discussion_id}")
def submit_topic(new_reply: NewReply, discussion_id: int, db: Session = Depends(get_db)):
    try:
        reply = Replies(
            discussion_id=discussion_id,
            user_id=new_reply.user_id,
            content=new_reply.content
        )
        db.add(reply)
        db.commit()
        db.refresh(reply)
        return {"message": "Reply submitted successfully.", "id": reply.id}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/replies/get_all_replies/{discussion_id}", response_model=ReplyList)
def get_all_replies(discussion_id: int, db: Session = Depends(get_db)):
    try:
        replies = db.query(Replies).filter(Replies.discussion_id == discussion_id).all()
        if not replies:
            raise HTTPException(status_code=404, detail="No replies found.")
        return {
            "replies": [
                {
                    "id": reply.id,
                    "discussion_id": reply.discussion_id,
                    "user_id": reply.user_id,
                    "content": reply.content,
                    "created_at": reply.created_at.isoformat(),
                    "is_deleted": reply.is_deleted,
                    "deleted_at": reply.deleted_at.isoformat() if reply.deleted_at else None,
                    "is_hidden": reply.is_hidden,
                    "hidden_at": reply.hidden_at.isoformat() if reply.hidden_at else None,
                }
                for reply in replies
            ],
            "total_count": len(replies),
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/replies/get_reply/{reply_id}", response_model=Reply)
def get_reply(reply_id: int, db: Session = Depends(get_db)):
    try:
        reply = db.query(Replies).filter(Replies.id == reply_id).first()
        if not reply:
            raise HTTPException(status_code=404, detail="Reply not found.")
        return {
            "id": reply.id,
            "discussion_id": reply.discussion_id,
            "user_id": reply.user_id,
            "content": reply.content,
            "created_at": reply.created_at.isoformat(),
            "is_deleted": reply.is_deleted,
            "deleted_at": reply.deleted_at.isoformat() if reply.deleted_at else None,
            "is_hidden": reply.is_hidden,
            "hidden_at": reply.hidden_at.isoformat() if reply.hidden_at else None,
        }
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/replies/update_reply/{reply_id}")
def update_reply(reply_id: int, updated_reply: NewReply, db: Session = Depends(get_db)):
    try:
        reply = db.query(Replies).filter(Replies.id == reply_id).first()
        if not reply:
            raise HTTPException(status_code=404, detail="Reply not found.")

        reply.content = updated_reply.content
        db.commit()
        return {"message": f"Reply with ID {reply_id} updated successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.delete("/replies/delete_reply/{reply_id}")
def delete_reply(reply_id: int, db: Session = Depends(get_db)):
    try:
        reply = db.query(Replies).filter(Replies.id == reply_id).first()
        if not reply:
            raise HTTPException(status_code=404, detail="Reply not found.")

        reply.is_deleted = True
        db.commit()
        return {"message": f"Reply with ID {reply_id} deleted successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
