import React, { useState } from 'react';

export default function TabToggle({ activeTab, setActiveTab }) {
  return (
    <div className="bg-gray-100 rounded-full justify-center shadow-inner mb-3">
      <button
        onClick={() => setActiveTab('daily')}
        className={`px-5 py-2 font-semibold rounded-full transition cursor-pointer
          ${activeTab === 'daily' ? 'bg-white text-purple-600 shadow' : 'text-gray-500'}`}
      >
        Daily Prompts
      </button>
      <button
        onClick={() => setActiveTab('attributes')}
        className={`px-5 py-2 font-semibold rounded-full transition cursor-pointer
          ${activeTab === 'attributes' ? 'bg-white text-purple-600 shadow' : 'text-gray-500'}`}
      >
        My Attributes
      </button>
    </div>
  );
}
