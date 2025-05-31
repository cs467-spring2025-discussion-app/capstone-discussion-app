export default function CommentFlair({ tagName }) {
    return (
        <span className='text-xs mr-2 px-5 text-purple-800 border border-purple-400 border-solid rounded-lg bg-gray-50 font-semibold shadow-md'>
            {tagName}
        </span>
    )
}   