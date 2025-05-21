export default function Badges() {
    return (
        <div>
            <div class="max-w-md mt-15 mx-auto p-5 bg-white border-solid border-gray-700 border-t-4 border-t-purple-500 rounded-lg shadow-md">
                <h2 class="text-lg font-semibold text-gray-800 mb-2">My Badges</h2>
                <p class="text-sm text-gray-600 mb-4">Badges are earned based on your participation and contribution quality.</p>
                
                <div class="flex flex-wrap gap-2">
                    <span class="bg-purple-100 text-purple-700 text-sm font-medium px-3 py-1 rounded-full">
                    Insightful
                    </span>

                    <span class="flex bg-gray-200 text-gray-500 text-sm font-medium px-3 py-1 rounded-full">
                    Top Contributor <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" strokeWidth={1.5} stroke="currentColor" className="size-6">
                    <path strokeLinecap="round" strokeLinejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 1 0-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 0 0 2.25-2.25v-6.75a2.25 2.25 0 0 0-2.25-2.25H6.75a2.25 2.25 0 0 0-2.25 2.25v6.75a2.25 2.25 0 0 0 2.25 2.25Z" />
                    </svg>
                    </span>
                </div>
            </div>
        </div>
    )
}