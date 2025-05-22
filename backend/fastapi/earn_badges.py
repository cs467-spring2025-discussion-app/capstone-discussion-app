from sqlalchemy import Column, select
from sqlalchemy.orm import Session
from tables import Badges, UserBadges, Votes

TOP_CONTRIBUTOR_MIN = 5
TOP_CONTRIBUTOR_BADGE_NAME = "Top Contributor"  # Must match existing DB badge


def check_top_contributor(db: Session, user_id: Column):
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
    count = db.query(Votes).filter_by(user_id=user_id, vote_type=1).count()
    # Award 'Top Contributor' badge if user does not already have it and user
    # has enough replies with at least one upvote
    if userBadge is None and count is not None and count >= TOP_CONTRIBUTOR_MIN:
        user_badge = UserBadges(user_id=user_id, badge_id=badge_id)
        db.add(user_badge)
        db.commit()
        db.refresh(user_badge)
    # If user does have badge, and they are below reply threshold, remove badge
    # from user (effective daily reset)
    else:
        db.delete(userBadge)
        db.commit()
