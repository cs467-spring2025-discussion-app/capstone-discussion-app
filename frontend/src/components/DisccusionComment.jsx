import React, { useState } from 'react';
import CommentFlair from './CommentFlair';

export default function DiscussionComment() {
    return (
        <div className='pl-5 shadow-lg rounded-lg mb-5 pt-5 bg-white w-full max-w-5xl h-35 text-left'>
            <CommentFlair />
            <CommentFlair />
            <CommentFlair />
            <p className='mt-2'>This is a multi-line comment to be updated dynamically in the near future. This is a multi-line comment to be updated dynamically in the near future. This is a multi-line comment to be updated dynamically in the near future. This is a multi-line comment to be updated dynamically in the near future.</p>
        </div>
    )
}