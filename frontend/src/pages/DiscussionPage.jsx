import React from 'react'
import Topic from "../components/Topic"
import TabToggle from '../components/TabToggle'

export default function DiscussionPage() {
    const postTextAreaId = React.useId();
    return (
        <div className="flex justify-center">
            <div className='flex flex-col items-center w-full max-w-xl p-4 text-center'>
                <div className='border-t-15 border-t-purple-400 h-35 p-5 rounded-lg shadow-lg w-2xl bg-white text-left mb-8'>
                    <h3 className='mt-2 mb-5 font-bold text-2xl'>Topic</h3>
                    <p>Question</p>
                </div>
                <div className='border-t-5 border-t-purple-400 p-5 rounded-lg shadow-lg w-2xl bg-white text-left mb-8'>
                    <label htmlFor={postTextAreaId}>
                        Join the discussion:
                        <textarea
                            name="postContent"
                            id={postTextAreaId}
                            placeholder="Share your thoughts...(280 chars max.)"
                            rows={4}
                            class='mt-5 mb-2 p-2.5 w-full text-gray-900 bg-gray-50 rounded-lg border border-gray-300 focus:ring-blue-500 focus:border-blue-500'
                        />
                    </label>
                    <div className='flex justify-between pb-5 border-b border-gray-300'>
                        <div className='text-sm text-gray-700'>
                            Char Count/280
                        </div>
                        <div className='text-sm text-gray-700'>
                            Responses limited to 10 per prompt
                        </div>
                    </div>
                    <div className='flex justify-end mt-2'>
                        <button className='ml-2 flex bg-purple-600 rounded-full px-5 py-2 font-bold text-sm text-white'>
                            Submit
                            <svg 
                            xmlns="http://www.w3.org/2000/svg"
                            fill="none"
                            viewBox="0 0 24 24"
                            strokeWidth={1.5}
                            stroke="currentColor"
                            className="size-6">
                            <path 
                                strokeLinecap="round"
                                strokeLinejoin="round"
                                d="M6 12 3.269 3.125A59.769 59.769 0 0 1 21.485 12 59.768 59.768 0 0 1 3.27 20.875L5.999 12Zm0 0h7.5" />
                        </svg>

                        </button>
                    </div>
                </div>
            </div>
        </div>
    )
}