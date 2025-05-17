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
 * @param {number} id - id of the discussion topic's replies to get
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
