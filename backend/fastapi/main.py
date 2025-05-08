from fastapi import FastAPI, HTTPException, Depends
from sqlalchemy import create_engine, select, text
from sqlalchemy.orm import Session, sessionmaker
from tables import Base, DiscussionTopics, Replies, Votes
from models import NewTopic, Topic, TopicList, NewReply, Reply, ReplyList, NewVote
from datetime import datetime
import os
from fastapi.middleware.cors import CORSMiddleware  # for the app.add_middleware function below

DATABASE_URL = os.getenv("DATABASE_URL")
if not DATABASE_URL:
    raise ValueError("DATABASE_URL environment variable is not set.")

engine = create_engine(DATABASE_URL, echo=True)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base.metadata.create_all(bind=engine)

app = FastAPI()

# This allows the browser and server to communicate locally
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:5173"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


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


@app.post("/topics/submit_topic")
# Submit a new topic
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


@app.get("/topics/get_daily_topics", response_model=TopicList)
# Get the five topics of the day
def get_daily_topics(db: Session = Depends(get_db)):
    try:
        stmt = select(DiscussionTopics).where(DiscussionTopics.is_active)
        daily_topics = db.execute(stmt).scalars().fetchmany(5)
        if not daily_topics:
            raise HTTPException(status_code=404, detail="No active topic found.")
        if len(daily_topics) < 5:
            raise HTTPException(status_code=404, detail="Less than 5 active topics found.")
        return TopicList(topics=daily_topics)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/topics/get_topic/{topic_id}", response_model=Topic)
# Get a topic by its ID
def get_topic(topic_id: int, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")
        return topic
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/topics/get_all_topics", response_model=TopicList)
# Get all topics
def get_all_topics(db: Session = Depends(get_db)):
    try:
        topics = db.query(DiscussionTopics).all()
        if not topics:
            raise HTTPException(status_code=404, detail="No topics found.")
        return TopicList(topics=topics)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/topics/update_topic/{topic_id}")
# Update a topic by its ID
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
# Activate a topic by its ID
def activate_topic(topic_id: int, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")
        if topic.is_active:
            raise HTTPException(status_code=400, detail="Topic is already active.")
        topic.is_active = True
        db.commit()
        return {"message": f"Topic with ID {topic_id} activated successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/topics/deactivate_topic/{topic_id}")
# Deactivate a topic by its ID
def deactivate_topic(topic_id: int, db: Session = Depends(get_db)):
    try:
        topic = db.query(DiscussionTopics).filter(DiscussionTopics.id == topic_id).first()
        if not topic:
            raise HTTPException(status_code=404, detail="Topic not found.")
        if not topic.is_active:
            raise HTTPException(status_code=400, detail="Topic is already inactive.")
        topic.is_active = False
        db.commit()
        return {"message": f"Topic with ID {topic_id} deactivated successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/topics/show_topic/{topic_id}")
# Mark a topic as shown by its ID
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
# Delete a topic by its ID
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
# Submit a new reply to a discussion
def submit_reply(new_reply: NewReply, discussion_id: int, db: Session = Depends(get_db)):
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
# Get all replies for a discussion
def get_all_replies(discussion_id: int, db: Session = Depends(get_db)):
    try:
        replies = db.query(Replies).filter(Replies.discussion_id == discussion_id).all()
        if not replies:
            raise HTTPException(status_code=404, detail="No replies found.")
        return ReplyList(replies=replies, total_count=len(replies))
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.get("/replies/get_reply/{reply_id}", response_model=Reply)
# Get a reply by its ID
def get_reply(reply_id: int, db: Session = Depends(get_db)):
    try:
        reply = db.query(Replies).filter(Replies.id == reply_id).first()
        if not reply:
            raise HTTPException(status_code=404, detail="Reply not found.")
        return reply
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.put("/replies/update_reply/{reply_id}")
# Update a reply by its ID
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
# Delete a reply by its ID
def delete_reply(reply_id: int, db: Session = Depends(get_db)):
    try:
        reply = db.query(Replies).filter(Replies.id == reply_id).first()
        if not reply:
            raise HTTPException(status_code=404, detail="Reply not found.")

        db.delete(reply)
        db.commit()
        return {"message": f"Reply with ID {reply_id} deleted successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


################################################################################
# Vote Endpoints
################################################################################

@app.post("/votes/submit_vote/{reply_id}")
# Submit a vote for a reply
def submit_vote(new_vote: NewVote, reply_id: int, db: Session = Depends(get_db)):
    try:
        # Check if reply exists
        reply = db.query(Replies).filter_by(id=reply_id).first()
        if not reply:
            raise HTTPException(status_code=404, detail="Reply not found.")

        # Check if vote already exists
        existing_vote = db.query(Votes).filter_by(
            user_id=new_vote.user_id,
            reply_id=reply_id
        ).first()

        if existing_vote:
            # Update existing vote
            existing_vote.vote_type = new_vote.vote_type
            existing_vote.updated_at = datetime.now()
        else:
            # Create new vote
            vote = Votes(
                user_id=new_vote.user_id,
                reply_id=reply_id,
                vote_type=new_vote.vote_type
            )
            db.add(vote)

        db.commit()
        return {"message": "Vote submitted successfully."}
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


#
if __name__ == "__main__":
    import uvicorn

    # use default port unless set in a variable (example: during GCP deployment)
    port = int(os.environ.get("PORT", 8080))
    uvicorn.run("main:app", host="0.0.0.0", port=port)
