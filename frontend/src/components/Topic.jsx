import React from 'react'

export default function Topic({ topic }) {
    return (
        <button className='w-full md:w-xl border-solid border-gray-700 p-8 pb-3 text-left bg-white mb-7 rounded-lg shadow-md'>
            <div className='flex justify-between mb-3'>
                <h3 className='text-lg font-semibold'>{topic.title}</h3>
                <span className='flex text-sm text-purple-600 bg-purple-100 rounded-lg px-3 py-1'><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-5 pr-1">
<path strokeLinecap="round" strokeLinejoin="round" d="M2.25 18 9 11.25l4.306 4.306a11.95 11.95 0 0 1 5.814-5.518l2.74-1.22m0 0-5.94-2.281m5.94 2.28-2.28 5.941" />
</svg>
Trending</span>
            </div>
            <p className='mb-3 border-b border-gray-200'>{topic.description}</p>
            <div className='flex text-sm text-gray-400'><svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-5 pr-1">
<path strokeLinecap="round" strokeLinejoin="round" d="M2.25 12.76c0 1.6 1.123 2.994 2.707 3.227 1.068.157 2.148.279 3.238.364.466.037.893.281 1.153.671L12 21l2.652-3.978c.26-.39.687-.634 1.153-.67 1.09-.086 2.17-.208 3.238-.365 1.584-.233 2.707-1.626 2.707-3.228V6.741c0-1.602-1.123-2.995-2.707-3.228A48.394 48.394 0 0 0 12 3c-2.392 0-4.744.175-7.043.513C3.373 3.746 2.25 5.14 2.25 6.741v6.018Z" />
</svg>
Response Count</div>
        </button>
    )   
}

