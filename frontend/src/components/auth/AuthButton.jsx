import React from "react";
import { useUser } from "@/hooks/useUser";
import * as authService from "@/api/authService";

const AuthButton = ({ action }) => {
  const { setUser } = useUser();

  /**
   * Show auth box appropriate for `action` or logout
   */
  async function handleAction() {
    // Logout
    if (action === "logout") {
      try {
        await authService.Logout();
        setUser(null);
        window.location.href="/";
      } catch (error) {
        console.error("Logout error:", error);
      }
    } else {
      // Login/Register/Update/Delete by showing modal
      const modal = document.getElementById(`${action}Form`);
      if (modal) {
        modal.showModal();
      } else {
        console.error(`Modal with ID ${action}Form not found`);
      }
    }
  }

  return (
    <button
      type="button"
      className="p-1 rounded-md hover:bg-blue-400 cursor-pointer"
      onClick={handleAction}
    >
      {action.charAt(0).toUpperCase() + action.slice(1)}
    </button>
  );
};

export default AuthButton;
