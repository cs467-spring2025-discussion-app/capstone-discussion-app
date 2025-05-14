import React from "react";

const LoginButton = ({ onClick }) => {
  return (
    <button
      type="button"
      className="text-purple-800 bg-purple-100 inline-block p-1 rounded-lg cursor-pointer"
      onClick={onClick}
    >
      Login
    </button>
  );
};

export default LoginButton;
