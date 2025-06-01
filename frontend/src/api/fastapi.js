/**
 * Get a discussion topic by id
 * @param {number} id - id of the discussion topic to get
 * @returns {Promise<Object>} Response data
 * @throws {Error} On get failure
 */
export async function GetDiscussion(id) {
  const resp = await fetch(`/fastapi/topics/get_topic/${id}`, {
    method: "GET",
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(errorData.error || `get_topic/${id} failed`);
  }

  return await resp.json();
}

/**
 * Get all replies to a discussion topic
 * @param {number} discussionId - id of the discussion topic to reply to
 * @returns {Promise<Object>} Response data
 * @throws {Error} On get failure
 */
export async function GetReplies(discussionId) {
  const resp = await fetch(`/fastapi/replies/get_all_replies/${discussionId}`, {
    method: "GET",
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(
      errorData.error || `replies/get_all_replies/${discussionId} failed`,
    );
  }

  return await resp.json();
}

/**
 * Get a reply by its id
 * @param {number} id - id of the reply to get
 * @returns {Promise<Object>} Response data
 * @throws {Error} On get failure
 */
export async function GetReply(id) {
  const resp = await fetch(`/fastapi/replies/get_reply/${id}`, {
    method: "GET",
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(
      errorData.error || `replies/get_reply/${id} failed`,
    );
  }

  return await resp.json();
}

/**
 * Submit a reply
 * @param {number} discussionId - id of the discussion topic to reply to
 * @param {number} userId - id of the user submitting the reply
 * @param {string} content - content of the reply
 * @returns {Promise<Object>} Response data
 * @throws {Error} On get failure
 */
export async function SubmitReply(discussionId, userId, content) {
  const resp = await fetch(`/fastapi/replies/submit_reply/${discussionId}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ user_id: userId, content: content }),
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(
      errorData.error || `replies/submit_reply/${discussionId} failed`,
    );
  }

  return await resp.json();
}

/**
 * Get the count of a user's replies to a discussion
 * @param {number} discussionId - id of the discussion topic to reply to
 * @param {number} userId - id of the user submitting the reply
 * @returns {Promise<Object>} Response data
 * @throws {Error} On get failure
 */
export async function GetUserReplyCount(discussionId, userId) {
  const resp = await fetch(`/fastapi/topics/user_reply_counts/${discussionId}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ user_id: userId}),
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(
      errorData.error || `replies/user_counts/ failed`,
    );
  }

  return await resp.json();
}

/**
 * Submit a vote for a reply
 * @param {number} replyId - ID of the reply to vote on
 * @param {string} userId - UUID of the user submitting the vote
 * @param {number} voteType - 1 for upvote, -1 for downvote, 0 to remove vote
 * @returns {Promise<Object>} Response data
 * @throws {Error} On vote submission failure
 */
export async function SubmitVote(replyId, userId, voteType) {
  const resp = await fetch(`/fastapi/votes/submit_vote/${replyId}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ user_id: userId, vote_type: voteType }),
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(
      errorData.error || `votes/submit_vote/${replyId} failed`,
    );
  }

  return await resp.json();
}

/**
 * Get vote counts for a reply
 * @param {number} replyId - ID of the reply
 * @returns {Promise<Object>} { total_votes, upvotes, downvotes }
 */
export async function GetVoteCounts(replyId) {
  const resp = await fetch(`/fastapi/votes/get_vote_counts/${replyId}`, {
    method: "GET",
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(errorData.error || `votes/get_vote_counts/${replyId} failed`);
  }

  return await resp.json();
}

/**
 * Flag a post
 * @param {number} replyId - id of the post being flagged
 * @param {number} userId - id of the user submitting the flag
 * @param {string} reason - reason for the flagging
 * @returns {Promise<Object>} Response data
 * @throws {Error} On submission failure
 */
export async function FlagReply(replyId, userId, reason) {
  const resp = await fetch (`/fastapi/replies/flag_reply/${replyId}`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ user_id: userId, reason }),
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(
      errorData.error || `flag_reply/${replyId} failed`,
    );
  }

  return await resp.json();
}

/**
 * Get tags for a user
 * @param {number} userId - id of the user whose tags to fetch
 * @returns {Promise<Object>} Response data
 * @throws {Error} On get failure
 */
export async function GetUserTags(userId) {
  const resp = await fetch(`/fastapi/user_tags/user`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ user_id: userId }),
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(errorData.error || `tags/user failed`);
  }

  return await resp.json();
}

/**
 * Get badges for a user
 * @param {number} userId - ID of the user
 * @returns {Promise<Object>} Response data
 * @throws {Error} On fetch failure
 */
export async function GetUserBadges(userId) {
  const resp = await fetch(`/fastapi/badges/user`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify({ user_id: userId }),
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(errorData.detail || `badges/user failed`);
  }

  return await resp.json();
}

export async function ToggleUserTag(userTag) {
  const resp = await fetch(`/fastapi/user_tags/toggle_tag`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userTag),
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(errorData.detail || "Failed to toggle user tag");
  }

  return await resp.json();
}

export async function GetAllTags() {
  const resp = await fetch(`/tags/get_all_tags`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  if (!resp.ok) {
    const errorData = await resp.json();
    throw new Error(errorData.detail || "Failed to fetch all tags");
  }

  return await resp.json();  // returns { tags: [...], total_count: number }
}
