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
