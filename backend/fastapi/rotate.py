# This file houses code related to the rotation of topics

from sqlalchemy.orm import Session
from sqlalchemy.sql import func
from datetime import datetime, timezone
from tables import DiscussionTopics


def rotate_topics(db: Session, count: int = 5):
    """
    This function rotates the daily topics by deactivating all topics and activating 5
    new ones that have never been used
    """

    # deactivate all topics
    db.query(DiscussionTopics).update({DiscussionTopics.is_active: False})

    # Get topics that haven't been activated and randomly select 5
    new_topics = (
        db.query(DiscussionTopics)
        .filter(DiscussionTopics.activated_at == None)
        .order_by(func.random())
        .limit(count)
        .all()
    )

    # if there aren't at least 5 new topics, then we can recycle
    if len(new_topics) < count:
        shortfall = count - len(new_topics)
        # get ids we already have so we don't duplicate
        exclude_these_ids = [i.id for i in new_topics]

        # query for additional topics
        add_topics = (
            db.query(DiscussionTopics)
            .filter(DiscussionTopics.id.notin_(exclude_these_ids))
            .order_by(func.random())
            .limit(shortfall)
            .all()
        )

        new_topics += add_topics

    # deliver an error if there are not 5 topics
    if len(new_topics) < count:
        raise ValueError(f"Not enough topics available.")

    # activate the topics in the new_topic set
    for topic in new_topics:
        topic.is_active = True
        if topic.activated_at is None:
            topic.activated_at = datetime.now(timezone.utc)
    db.commit()
