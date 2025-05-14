import React from "react";

const LogoutButton = ({ onClick }) => {
  return (
    <button
      type="button"
      className="text-purple-800 bg-purple-100 inline-block p-1 rounded-lg cursor-pointer"
      onClick={onClick}
    >
      Logout
    </button>
  );
};

export default LogoutButton;
