import React from "react";
import { useUser } from "@/hooks/useUser";
import * as authService from "@/api/authService";

const AuthButton = ({ action }) => {
  const { setUser } = useUser();

  const handleAction = async () => {
    // Logout
    if (action === "logout") {
      try {
        await authService.Logout();
        setUser(null);
      } catch (error) {
        console.error("Logout error:", error);
      }
    } else {
      // Login/Register/Delete by showing modal
      const modal = document.getElementById(`${action}Form`);
      if (modal) {
        modal.showModal();
      } else {
        console.error(`Modal with ID ${action}Form not found`);
      }
    }
  };

  return (
    <button
      type="button"
      className="text-purple-800 bg-purple-100 inline-block p-1 rounded-lg cursor-pointer"
      onClick={handleAction}
    >
      {action.charAt(0).toUpperCase() + action.slice(1)}
    </button>
  );
};

export default AuthButton;
