import { useUser } from "@/hooks/useUser";
import { useEffect, useState, useRef } from "react";
import AuthButton from "@/components/auth/AuthButton";
import ModalAuthForm from "@/components/auth/form/ModalAuthForm";
import * as authService from "@/api/authService";

const AuthMenu = () => {
  const [isCheckingStatus, setIsCheckingStatus] = useState(true);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const { user, setUser } = useUser();
  const dropdownRef = useRef(null);

  useEffect(() => {
    /**
     * Check login status before displaying menu options
     */
    async function checkLoginStatus() {
      try {
        const userData = await authService.Whoami();
        setUser(userData);
      } catch (error) {
        console.error("Error checking login status:", error);
        setUser(null);
      } finally {
        setIsCheckingStatus(false);
      }
    }

    checkLoginStatus();
  }, [setUser]);

  useEffect(() => {
    // Close dropdown when clicking outside of it
    // Ref: https://www.codemzy.com/blog/reactjs-dropdown-component
    function handleOutsideClick(e) {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target)) {
        setIsDropdownOpen(false);
      }
    }
    // Close dropdown when pressing <Escape>
    function handleEscapeKey(event) {
      if (event.key === "Escape") {
        setIsDropdownOpen(false);
      }
    }

    if (isDropdownOpen) {
      document.addEventListener("mousedown", handleOutsideClick);
      document.addEventListener("keyup", handleEscapeKey);
    }
    return function removeListener() {
      document.removeEventListener("mousedown", handleOutsideClick);
      document.removeEventListener("keyup", handleEscapeKey);
    };
  }, [isDropdownOpen]);

  /**
   * Toggle dropdown menu state
   */
  function toggleDropdown() {
    setIsDropdownOpen(!isDropdownOpen);
  }

  return (
    <div className="flex justify-end m-2">
      {isCheckingStatus ? ( // waiting for logged-in check
        <span className="inline-block">&nbsp;</span>
      ) : !user ? ( // no user is logged in
        <div className="bg-blue-300 rounded-lg shadow-lg font-bold">
          <AuthButton action="login" />
        </div>
      ) : (
        // user is logged in
        <div className="relative" ref={dropdownRef}>
          <button
            className="bg-blue-300 font-bold p-1
              border-1 border-transparent rounded-md shadow-lg
              cursor-pointer
              hover:border-black"
            type="button"
            onClick={toggleDropdown}
          >
            <div>{user.email}</div>
          </button>
          {isDropdownOpen && (
            <div
              id="authDropdown"
              className="text-lg bg-blue-300 max-w-max 
                border-1 rounded-md shadow-xl 
                absolute right-0 mt-2 flex flex-col"
            >
              <AuthButton action="update" />
              <AuthButton action="logout" />
              <AuthButton action="delete" />
            </div>
          )}
        </div>
      )}

      <ModalAuthForm action="login" />
      <ModalAuthForm action="update" />
      <ModalAuthForm action="delete" />
    </div>
  );
};

export default AuthMenu;
