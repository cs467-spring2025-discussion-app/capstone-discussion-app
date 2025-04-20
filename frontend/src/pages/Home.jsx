import React from 'react'
import Topic from "../components/Topic"
import TabToggle from '../components/TabToggle'

export default function Home() {
    return (
        <div className="flex justify-center">
            <div className='flex flex-col items-center w-full max-w-xl p-4 text-center'>
                <h1 className="text-4xl font-bold mb-5">Daily Discussion</h1>
                <h3 className="font-light mb-3">Join the conversation, anonymously.</h3>
                <div className="flex justify-around w-xs mb-5">
                    <div className="text-purple-600 font-bold bg-purple-100 rounded-lg px-3 py-1">
                        Anonymous
                    </div>
                    <div className="text-purple-600 font-bold bg-purple-100 rounded-lg px-3 py-1">
                        Resets at Midnight
                    </div>
                </div>
                <div>
                    <TabToggle />
                </div>
                    <p className='font-light mt-3 mb-5'>Select a prompt to join the discussion. All discussions reset at midnight.</p>
                    <Topic />
                    <Topic />
                    <Topic />
                    <Topic />
                    <Topic />
            </div>
        </div>
    )
}