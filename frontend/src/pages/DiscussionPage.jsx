import { useEffect, useId, useState } from "react";
import { useParams } from "react-router-dom";
import { useUser } from "@/hooks/useUser";
import DiscussionComment from "@/components/DisccusionComment";
import HomeLink from "@/components/HomeLink";
import NotFound from "@/pages/NotFound";
import * as fastapiService from "@/api/fastapi";

export default function DiscussionPage() {
  const postTextAreaId = useId();
  const [text, setText] = useState("");
  const [comments, setComments] = useState([]);
  const [lastCommentId, setLastCommentId] = useState(null);

  const { user } = useUser();
  const { discussionId } = useParams();
  const [discussion, setDiscussion] = useState(null);
  const [loading, setLoading] = useState(true);

  // Get discussion topic data
  useEffect(() => {
    setLoading(true);
    async function getDiscussionData() {
      try {
        const data = await fastapiService.GetDiscussion(discussionId);
        setDiscussion(data);
      } catch (error) {
        console.error("Error getting discussion: ", error);
      } finally {
        setLoading(false);
      }
    }
    getDiscussionData();
  }, [discussionId]);

  // Get topic replies
  useEffect(() => {
    async function getReplies() {
      setLoading(true);
      try {
        const data = await fastapiService.GetReplies(discussionId);
        setComments(data.replies); // Need to access .replies array
      } catch (error) {
        console.error(`Error getting replies of ${discussionId}: `, error);
      } finally {
        setLoading(false);
      }
    }
    getReplies();
  }, [discussionId]);

  // Submit reply
  async function handleSubmit() {
    if (!text.trim()) return;

    try {
      const resp = await fastapiService.SubmitReply(
        discussionId,
        user.userID,
        text,
      );
      const newComment = await fastapiService.GetReply(resp.id);
      setLastCommentId(newComment.id);
      setComments([newComment, ...comments]);
      setText("");
    } catch (error) {
      console.error("Error submitting reply: ", error);
    }
  }

  // Wait to render while loading
  if (loading) {
    return <div>Loading...</div>;
    // Not found on invalid topic id
  } else if (!discussion) {
    return <NotFound />;
  } else {
    return (
      <div>
        <div className="flex flex-col justify-center">
          <HomeLink />
          <div className="flex flex-col items-center w-full max-w-4xl p-4 text-center">
            <div className="w-full max-w-5xl border-t-15 border-t-purple-400 h-35 p-5 rounded-lg shadow-lg bg-white text-left mb-8">
              <h3 className="mt-2 mb-5 font-bold text-2xl">
                {discussion.title}
              </h3>
              <p>{discussion.prompt}</p>
            </div>
            {user && (
              <div className="border-t-5 border-t-purple-400 p-5 rounded-lg shadow-lg w-full max-w-5xl bg-white text-left mb-8">
                <label htmlFor={postTextAreaId} className="font-semibold">
                  Join the discussion:
                  <textarea
                    name="postContent"
                    value={text}
                    id={postTextAreaId}
                    placeholder="Share your thoughts...(280 chars max.)"
                    onChange={(e) => setText(e.target.value)}
                    maxLength={280}
                    rows={4}
                    className="font-normal mt-5 mb-2 p-2.5 w-full text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500"
                  />
                </label>
                <div className="flex justify-between pb-5 border-b border-gray-300">
                  <div className="text-sm text-gray-700">{text.length}/280</div>
                  <div className="text-sm text-gray-700">
                    Responses limited to 10 per prompt
                  </div>
                </div>
                <div className="flex justify-end mt-2">
                  <button
                    className="cursor-pointer ml-2 flex bg-purple-600 rounded-full px-5 py-2 font-bold text-sm text-white"
                    onClick={handleSubmit}
                  >
                    Submit
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
                        d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5"
                      />
                    </svg>
                  </button>
                </div>
              </div>
            )}
            {/* Display discussion replies */}
            {comments.map((comment) => (
              <DiscussionComment
                key={comment.id}
                comment={comment}
                fadeIn={comment.id === lastCommentId}
              />
            ))}
          </div>
        </div>
      </div>
    );
  }
}
