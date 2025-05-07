import React, { useState } from 'react';
import CommentFlair from './CommentFlair';

export default function DiscussionComment({ content, fadeIn }) {
    const [votes, setVotes] = useState(0)
    const handleUpvote = () => {
        setVotes((prev) => prev + 1)
    };
    const handleDownvote = () => {
        setVotes((prev) => (prev > 0 ? prev - 1 : 0))
    };

    return (
        <div className={`
        pl-5 shadow-lg rounded-lg mb-5 pt-5 bg-white w-full max-w-5xl text-left
        ${fadeIn ? 'opacity-0 animate-fade-in' : ''}
        `}
        >
            <CommentFlair />
            <CommentFlair />
            <CommentFlair />
            <p className='mt-2'>{content}</p>
            <div className='w-full flex justify-between items-center text-sm text-gray-500 my-5 pr-5'>
                <div className=''>2h ago</div>
                <div className='flex items-center gap-5'>
                    <button 
                      className="p-2 rounded hover:bg-gray-100"
                      onClick={handleUpvote}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                        <path strokeLinecap="round" strokeLinejoin="round" d="m4.5 15.75 7.5-7.5 7.5 7.5" />
                        </svg>
                    </button>

                    <p>{votes}</p>

                    <button
                      className="p-2 rounded hover:bg-gray-100"
                      onClick={handleDownvote}
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                        <path strokeLinecap="round" strokeLinejoin="round" d="m19.5 8.25-7.5 7.5-7.5-7.5" />
                        </svg>
                    </button>

                    <button className="p-2 rounded hover:bg-gray-100">
                        <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="#c084fc" className="size-6">
                        <path fillRule="evenodd" d="M3 2.25a.75.75 0 0 1 .75.75v.54l1.838-.46a9.75 9.75 0 0 1 6.725.738l.108.054A8.25 8.25 0 0 0 18 4.524l3.11-.732a.75.75 0 0 1 .917.81 47.784 47.784 0 0 0 .005 10.337.75.75 0 0 1-.574.812l-3.114.733a9.75 9.75 0 0 1-6.594-.77l-.108-.054a8.25 8.25 0 0 0-5.69-.625l-2.202.55V21a.75.75 0 0 1-1.5 0V3A.75.75 0 0 1 3 2.25Z" clipRule="evenodd" />
                        </svg>
                    </button>
                    </div>
            </div>
        </div>
    )
}