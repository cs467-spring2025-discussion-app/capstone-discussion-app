import React, { useState, useEffect } from "react";
import CommentFlair from "@/components/CommentFlair";
import { useUser } from "@/hooks/useUser";
import { SubmitVote, GetVoteCounts } from "@/api/fastapi";

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

  const commentTimeElapsed = calculateTimeElapsed(comment.created_at);

  return (
    <div
      className={`
        pl-5 shadow-lg rounded-lg mb-5 pt-5 bg-white w-full max-w-5xl text-left
        ${fadeIn ? "opacity-0 animate-fade-in" : ""}
        `}
    >
      <CommentFlair />
      <CommentFlair />
      <CommentFlair />
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

          <button className="p-2 rounded hover:bg-gray-100">
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
  );
}

