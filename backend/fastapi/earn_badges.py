from sqlalchemy import Column, func, select
from sqlalchemy.orm import Session
from tables import Badges, Replies, UserBadges, UserReputation, Votes

TOP_CONTRIBUTOR_BADGE_NAME = "Top Contributor"  # Must match existing DB badge
INSIGHTFUL_BADGE_NAME = "Insightful"  # Must match existing DB badge
TOP_CONTRIBUTOR_MIN = 5
INSIGHTFUL_MIN_RATIO = 5
INSIGHTFUL_MIN_UPVOTES = 5


def check_insightful(db: Session, user_id: Column):
    # Get "Insightful" badge id
    badge_id = db.execute(
        select(Badges.id)
        .select_from(Badges)
        .where(Badges.name == INSIGHTFUL_BADGE_NAME)
    ).scalar()

    if badge_id is None:
        raise ValueError(f"'{INSIGHTFUL_BADGE_NAME}' is not a valid badge name")

    userBadge = (
        db.query(UserBadges).filter_by(user_id=user_id, badge_id=badge_id).first()
    )

    # Calculate the user's up/downvote ratio from their reputation
    # If their ratio is above the minimum, award the Insightful badge
    reputation = db.query(UserReputation).filter_by(user_id=user_id).first()

    # Abort check if reputation hasn't been initialized yet
    if reputation is None:
        return

    # Remove badge if user no longer meets conditions
    if userBadge is not None and (
        reputation.upvotes_received < INSIGHTFUL_MIN_UPVOTES
        or (
            reputation.downvotes_received > 0
            and reputation.upvotes_received / reputation.downvotes_received
            < INSIGHTFUL_MIN_RATIO
        )
    ):
        db.delete(userBadge)
        db.commit()
    # Grant badge if user meets conditions
    elif reputation.upvotes_received >= INSIGHTFUL_MIN_UPVOTES and (
        reputation.downvotes_received == 0
        or reputation.upvotes_received / reputation.downvotes_received
        >= INSIGHTFUL_MIN_RATIO
    ):
        if userBadge is None:
            userBadge = UserBadges(user_id=user_id, badge_id=badge_id)
        db.add(userBadge)
        db.commit()
        db.refresh(userBadge)


def check_top_contributor(db: Session, user_id: Column):
    """
    Assess a user's earned votes and deterime if they've earned the "Top
    Contributor" badge. Because we do not track reply history, this is a badge
    that can be earned or lost daily, encouraging frequent contributions.
    """

    # Get "Top Contributor" badge id
    badge_id = db.execute(
        select(Badges.id)
        .select_from(Badges)
        .where(Badges.name == TOP_CONTRIBUTOR_BADGE_NAME)
    ).scalar()

    if badge_id is None:
        raise ValueError(f"'{TOP_CONTRIBUTOR_BADGE_NAME}' is not a valid badge name")

    # Try to find existing user badge
    userBadge = (
        db.query(UserBadges).filter_by(user_id=user_id, badge_id=badge_id).first()
    )

    # Count the user's total unique upvotes across vote history
    # User must have one upvote on the minimum number of replies
    # Own votes do not count
    statement = (
        select(func.count("*"))
        .select_from(Votes)
        .join(Replies, Votes.reply_id == Replies.id)
        .where(Votes.user_id != user_id)
        .where(Votes.vote_type == 1)
    )
    count = db.execute(statement).scalar()

    # Award 'Top Contributor' badge if user does not already have it and user
    # has enough replies with at least one upvote
    if userBadge is None and count is not None and count >= TOP_CONTRIBUTOR_MIN:
        user_badge = UserBadges(user_id=user_id, badge_id=badge_id)
        db.add(user_badge)
        db.commit()
        db.refresh(user_badge)
    # If user does have badge, and they are below reply threshold, remove badge
    # from user (effective daily reset)
    elif userBadge is not None and count is not None and count < TOP_CONTRIBUTOR_MIN:
        db.delete(userBadge)
        db.commit()
