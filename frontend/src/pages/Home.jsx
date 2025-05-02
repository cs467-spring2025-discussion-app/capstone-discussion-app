import React, { useEffect, useState } from 'react';
import Topic from "@/components/Topic";
import TabToggle from '@/components/TabToggle';

export default function Home() {

    // create variables for topics coming from backend
    const [topics, setTopics] = useState([]);

    useEffect(() => {
        fetch("http://localhost:8000/topics/get_all_topics")
            // convert response to JSON
            .then(res => res.json())
            // store response
            .then(data => setTopics(data.topics))
            .catch(err => console.error("Failed to get topics:", err));
    }, []);

    return (
        <div className="flex justify-center">
            <div className='flex flex-col items-center w-full max-w-screen p-4 text-center'>
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
                <div className='md:flex md:flex-wrap md:flex-row md:justify-center md:gap-2 md:max-w-screen'>
                    {topics.map(topic => (
                        <Topic key={topic.id} topic={topic} />
                    ))}
                </div>
            </div>
        </div>
    )
}
