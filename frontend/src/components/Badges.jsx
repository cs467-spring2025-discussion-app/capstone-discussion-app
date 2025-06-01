import React, { useEffect, useState } from "react";

export default function Badges({ userId }) {
    const [badges, setBadges] = useState([]);

    useEffect(() => {
        async function loadBadges() {
            try {
                // get all the badges and the ones the user has earned
                const [allRes, userRes] = await Promise.all([
                    fetch("/fastapi/badges"),
                    fetch("/fastapi/badges/user", {
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({ user_id: userId }),
                    }),
                ]);
                
                // parse the response
                const all = await allRes.json();
                const { user_badges } = await userRes.json();

                // put earned badges into a set so we can look them up
                const unlockedIds = new Set(user_badges.map(b => b.id));

                // merge in whether each badge is unlocked
                const merged = all.map(b => ({
                    ...b,
                    unlocked: unlockedIds.has(b.id),
                }));

                setBadges(merged);
            } catch (err) {
                console.error("Failed to load badges:", err);
            } 
        }
        loadBadges();
    }, [userId]);

    return (
        <div>
            <div className="max-w-md mt-15 mx-auto p-5 bg-white border-solid border-gray-700 border-t-4 border-t-purple-500 rounded-lg shadow-md">
                <h2 className="text-lg font-semibold text-gray-800 mb-2">My Badges</h2>
                <p className="text-sm text-gray-600 mb-4">Badges are earned based on your participation and contribution quality.</p>
                
                <div className="flex flex-wrap gap-2">
                    {badges.map(badge => (
                        <span 
                            key={badge.id}
                            className={`flex text-sm font-medium px-3 py-1 rounded-full ${
                                badge.unlocked
                                    ? "bg-purple-100 text-purple-700"
                                    : "bg-gray-200 text-gray-500"
                            }`}
                        >
                            {badge.name}
                            {!badge.unlocked && (
                                <svg
                                xmlns="http://www.w3.org/2000/svg"
                                fill="none"
                                viewBox="0 0 24 24"
                                strokeWidth={1.5}
                                stroke="currentColor"
                                className="size-6 ml-1"
                              >
                                <path
                                  strokeLinecap="round"
                                  strokeLinejoin="round"
                                  d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z"
                                />
                              </svg>
                            )}
                        </span>
                    ))}
                </div>
            </div>
        </div>
    )
}