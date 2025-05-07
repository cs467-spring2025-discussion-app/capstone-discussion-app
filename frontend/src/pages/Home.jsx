import React, { useEffect, useState } from 'react';
import Topic from "@/components/Topic";
import TabToggle from '@/components/TabToggle';

export default function Home() {

    // create variables for topics coming from backend
    const [topics, setTopics] = useState([]);

    // create variables for error message in case we don't get topics back
    const [errorMessage, setErrorMessage] = useState("");

    // get discussion topics from backend 
    useEffect(() => {
        fetch("http://localhost:8000/topics/get_all_topics")
            // convert response to JSON
            .then(res => res.json())
            // store response
            .then(data => {
                setTopics(data.topics);
                if (data.topics.len < 5) {
                    setErrorMessage("Fewer than 5 discussion topics available. Please check back later.");
                } else {
                    setErrorMessage("");
                }
            })
            .catch(err => console.error("Failed to get topics:", err));
            setErrorMessage("Error getting topics. Please check back later.");
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

                {errorMessage && (
                    <div className="text-red-600 bg-red-100 rounded-lg p-2 mb-4 max-w-md">
                        {errorMessage}
                    </div>
                )}

                <div className='md:flex md:flex-wrap md:flex-row md:justify-center md:gap-2 md:max-w-screen'>
                    {/* Present the list of cards that are sent from the backend */}
                    {topics.map(topic => (
                        <Topic key={topic.id} topic={topic} />
                    ))}
                </div>
            </div>
        </div>
    )
}
