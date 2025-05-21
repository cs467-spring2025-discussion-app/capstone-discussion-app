import { useState } from 'react';

const ATTRIBUTE_CATEGORIES = {
  'Political Leaning': ['Progressive', 'Liberal', 'Moderate', 'Conservative', 'Libertarian'],
  'Professional Background': ['Tech', 'Healthcare', 'Education', 'Finance', 'Arts', 'Trade', 'Student'],
  Values: ['Environmentalist', 'Privacy Advocate', 'Free Speech', 'Internationalist', 'Traditionalist'],
};


export default function Attributes() {
    const [selected, setSelected] = useState([]);
    const [customAttr, setCustomAttr] = useState('');

    const toggleAttribute = (attr) => {
        if (selected.includes(attr)) {
        setSelected(selected.filter(a => a !== attr)); // "untoggle" attribute
        } else if (selected.length < 5) {
        setSelected([...selected, attr]);
        }
    };

    const addCustomAttribute = () => {
        // add customer attr if it's not already added
        if (customAttr && !selected.includes(customAttr) && selected.length < 5) {
        setSelected([...selected, customAttr]);
        setCustomAttr('');
        }
    };
    return (
        <div className="max-w-md mt-15 mx-auto p-5 bg-white border-solid border-gray-700 border-t-4 border-t-purple-500 rounded-lg shadow-md">
            <h2 className="text-xl font-semibold mb-1">My Attributes</h2>
            <p className="text-sm text-gray-600 mb-4">Select up to 5 attributes that represent your perspective. These will be shown with your comments.</p>

            <div className="flex flex-wrap gap-2 mb-4">
                {selected.length === 0 ? (
                <span className="text-gray-400 text-sm">No attributes selected</span>
                ) : (
                selected.map(attr => (
                <span
                    key={attr}
                    className="px-3 py-1 bg-purple-100 text-purple-700 rounded-full text-sm cursor-pointer hover:bg-purple-200"
                    onClick={() => toggleAttribute(attr)}
                >
                    {attr}
                </span>
                ))
                )}
        </div>

        <div className="flex gap-2 mb-6">
            <input
            type="text"
            placeholder="Add custom attribute..."
            value={customAttr}
            onChange={(e) => setCustomAttr(e.target.value)}
            className="flex-grow px-3 py-2 border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-purple-300"
            />
            <button
            onClick={addCustomAttribute}
            className="bg-purple-500 text-white px-4 py-2 rounded-lg hover:bg-purple-600 transition"
            >
            + Add
            </button>
        </div>

        {Object.entries(ATTRIBUTE_CATEGORIES).map(([category, items]) => (
            <div key={category} className="mb-5">
            <h3 className="text-sm font-semibold text-gray-700 mb-2">{category}</h3>
            <div className="flex flex-wrap gap-2">
                {items.map(attr => (
                <button
                    key={attr}
                    onClick={() => toggleAttribute(attr)}
                    className={`px-3 py-1 rounded-full text-sm border transition ${
                    selected.includes(attr)
                        ? 'bg-purple-500 text-white border-transparent'
                        : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
                    }`}
                >
                    {attr}
                </button>
                ))}
            </div>
            </div>
        ))}
        </div>
    )
    
}

