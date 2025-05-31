import React, { useState, useEffect } from "react";
import CommentFlair from "@/components/CommentFlair";
import { useUser } from "@/hooks/useUser";
import { SubmitVote, GetVoteCounts, FlagReply, GetUserTags } from "@/api/fastapi";

/**
 * Calculate the time elapsed since a given time
 * @param {string} time - time in format e.g. "2025-05-15T04:03:51.502206"
 * @returns {string} time elapsed in days, hours, or minutes
 */
function calculateTimeElapsed(time) {
  const date = new Date(time);
  const now = new Date();

  // Calculate time difference in milliseconds with offset
  const offset = now.getTimezoneOffset()
  const diff = now.getTime() - date.getTime() + (offset * 60 * 1000);

  const daysDiff = Math.floor(diff / (1000 * 60 * 60 * 24));
  const hoursDiff = Math.floor(diff / (1000 * 60 * 60));
  const minutesDiff = Math.floor(diff / (1000 * 60));

  if (daysDiff >= 1) {
    return `${daysDiff}d`;
  } else if (hoursDiff >= 1) {
    return `${hoursDiff}h`;
  } else {
    return `${minutesDiff}m`;
  }
}

export default function DiscussionComment({ comment, fadeIn }) {
  const [votes, setVotes] = useState(0);
  const { user } = useUser();
  const [showFlagDialog, setShowFlagDialog] = useState(false);
  const [flagReason, setFlagReason] = useState("");
  const [userTags, setUserTags] = useState([]);


  useEffect(() => {
    async function fetchVotes() {
      try {
        const data = await GetVoteCounts(comment.id);
        setVotes(data.total_votes);
      } catch (error) {
        console.error("Failed to fetch vote count:", error);
      }
    }
    fetchVotes();
  }, [comment.id]);

  useEffect(() => {
    async function fetchUserTags() {
      try {
        const tags = await GetUserTags(comment.user_id);
        setUserTags(tags);
      } catch (error) {
        console.error("Failed to fetch user tags:", error);
      }
    }
  fetchUserTags();
}, [comment.user_id]);

  const handleUpvote = async () => {
    if (!user) return;
    try {
      await SubmitVote(comment.id, user.userID, 1);
      const data = await GetVoteCounts(comment.id);
      setVotes(data.total_votes);
    } catch (error) {
      console.error("Failed to upvote:", error);
    }
  };

  const handleDownvote = async () => {
    if (!user) return;
    try {
      await SubmitVote(comment.id, user.userID, -1);
      const data = await GetVoteCounts(comment.id);
      setVotes(data.total_votes);
    } catch (error) {
      console.error("Failed to downvote:", error);
    }
  };

  const handleFlagSubmit = async () => {
    if(!user || !flagReason.trim()) return;
    try {
      await FlagReply(comment.id, user.userID, flagReason);
      alert("Reply flagged for review");
      setShowFlagDialog(false);
      setFlagReason("");
    } catch (error) {
      console.error("Failed to flag reply:", error);
      alert("There was an error. Please try again.");
    }
  };

  const commentTimeElapsed = calculateTimeElapsed(comment.created_at);

  return (
    <>
      <div
        className={`
          pl-5 shadow-lg rounded-lg mb-5 pt-5 bg-white w-full max-w-5xl text-left
          ${fadeIn ? "opacity-0 animate-fade-in" : ""}
          `}
      >
        <div className="flex flex-wrap gap-2">
          {userTags.map((tag) => (
            <CommentFlair key={tag.id} tagName={tag.name} />
          ))}
        </div>

        <p className="mt-2">{comment.content}</p>
        <div className="w-full flex justify-between items-center text-sm text-gray-500 my-5 pr-5">
          <div className="">{commentTimeElapsed}</div>
          <div className="flex items-center gap-5">
            <button
              className="p-2 rounded hover:bg-gray-100"
              onClick={handleUpvote}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="m4.5 15.75 7.5-7.5 7.5 7.5"
                />
              </svg>
            </button>

            <p>{votes}</p>

            <button
              className="p-2 rounded hover:bg-gray-100"
              onClick={handleDownvote}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                strokeWidth={1.5}
                stroke="currentColor"
                className="size-6"
              >
                <path
                  strokeLinecap="round"
                  strokeLinejoin="round"
                  d="m19.5 8.25-7.5 7.5-7.5-7.5"
                />
              </svg>
            </button>

            <button 
              className="p-2 rounded hover:bg-gray-100"
              onClick={() => setShowFlagDialog(true)}
            >
              <svg
                xmlns="http://www.w3.org/2000/svg"
                viewBox="0 0 24 24"
                fill="#c084fc"
                className="size-6"
              >
                <path
                  fillRule="evenodd"
                  d="M3 2.25a.75.75 0 0 1 .75.75v.54l1.838-.46a9.75 9.75 0 0 1 6.725.738l.108.054A8.25 8.25 0 0 0 18 4.524l3.11-.732a.75.75 0 0 1 .917.81 47.784 47.784 0 0 0 .005 10.337.75.75 0 0 1-.574.812l-3.114.733a9.75 9.75 0 0 1-6.594-.77l-.108-.054a8.25 8.25 0 0 0-5.69-.625l-2.202.55V21a.75.75 0 0 1-1.5 0V3A.75.75 0 0 1 3 2.25Z"
                  clipRule="evenodd"
                />
              </svg>
            </button>
          </div>
        </div>
      </div>

      {showFlagDialog && (
        <div className="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 z-50">
          <div className="bg-white p-6 rounded-lg w-96 shadow-lg">
            <h2 className="text-lg font-semibold mb-2">Report this Post</h2>
            <p className="mb-4 text-sm text-gray-600">
              Please explain why you are flagging this post:
            </p>
            <textarea
              className="w-full p-2 border rounded mb-4"
              rows={3}
              value={flagReason}
              onChange={(e) => setFlagReason(e.target.value)}
              placeholder="Reason for flagging..."
            />
            <div className="flex justify-end gap-4">
              <button
                onClick={() => setShowFlagDialog(false)}
                className="text-gray-600 hover:text-black"
              >
                Cancel
              </button>
              <button
                onClick={handleFlagSubmit}
                className="bg-red-500 text-white px-4 py-2 rounded hover:bg-red-600"
              >
                Submit
              </button>
            </div>
          </div>
        </div>
      )}
  </>
);
}

