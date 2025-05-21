import React, { useEffect, useState } from 'react';
import Topic from "@/components/Topic";
import TabToggle from '@/components/TabToggle';
import AttributePage from './AttributePage';

export default function Home() {
    const [activeTab, setActiveTab] = useState('daily');

    // create variables for topics coming from backend
    const [topics, setTopics] = useState([]);

    // create variables for error message in case we don't get topics back
    const [errorMessage, setErrorMessage] = useState("");

    // get discussion topics from backend 
    useEffect(() => {
        //fetch asynchronously
        Promise.all([
            fetch("/fastapi/topics/get_daily_topics").then(res => res.json()),
            fetch("/fastapi/replies/counts").then(res => res.json())
        ])
        .then(([topicsData, repliesData]) => {
            const replyMap = new Map();
            repliesData.replies_counts.forEach(r => {
                replyMap.set(r.discussion_id, r.reply_count);
            });
            const topicsWithReplies = topicsData.topics.map(topic => ({
                ...topic,
                reply_count: replyMap.get(topic.id) || 0
            }));
            setTopics(topicsWithReplies);
        
            if (topicsWithReplies.length < 5) {
                setErrorMessage("Fewer than 5 discussion topics available. Please check back later.");
            } else {
                setErrorMessage("");
            }
        })
        .catch(err => {
            console.error("Failed to get topics:", err);
            setErrorMessage("Error getting topics. Please check back later.");
        });
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
                    <TabToggle activeTab={activeTab} setActiveTab={setActiveTab}/>
                </div>
                <div>
                    {activeTab === 'daily' && (
                    <div>
                        <p className='font-light mt-3 mb-5'>
                        Select a prompt to join the discussion. All discussions reset at midnight.
                        </p>
                        {errorMessage && (
                        <div className="text-red-600 bg-red-100 rounded-lg p-2 mb-4 max-w-md">
                            {errorMessage}
                        </div>
                        )}
                        <div className='md:flex md:flex-wrap md:flex-row md:justify-center md:gap-2 md:max-w-screen'>
                        {topics.map(topic => (
                            <Topic key={topic.id} topic={topic} />
                        ))}
                        </div>
                    </div>
                    )}
                    {activeTab === 'attributes' && <AttributePage />}
                </div>
             </div>
        </div>
    )
}
