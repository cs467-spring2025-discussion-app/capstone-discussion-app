import { useUser } from "@/hooks/useUser";
import { useEffect, useState } from "react";
import AuthButton from "@/components/auth/AuthButton";
import ModalAuthForm from "@/components/auth/form/ModalAuthForm";
import * as authService from "@/api/authService";

const AuthMenu = () => {
  const [isCheckingStatus, setIsCheckingStatus] = useState(true);
  const [isDropdownOpen, setIsDropdownOpen] = useState(false);
  const { user, setUser } = useUser();

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

  /**
   * Toggle dropdown menu state
   */
  function toggleDropdown() {
    setIsDropdownOpen(!isDropdownOpen);
  }

  return (
    <div className="m-2">
      {isCheckingStatus ? ( // waiting for logged-in check
        <span className="inline-block p-1">&nbsp;</span>
      ) : !user ? ( // no user is logged in
        <div className="absolute right-0 p-2 flex gap-2 items-center">
          <AuthButton action="login" />
        </div>
      ) : (
        // user is logged in
        <div className="absolute right-0 p-2">
          <button
            className="bg-blue-300 p-2 rounded-lg shadow-lg font-bold"
            type="button"
            id="authDropdown"
            onClick={toggleDropdown}
          >
            <div>{user.email}</div>
          </button>
          {isDropdownOpen && (
            <div className="absolute right-0 m-2 flex flex-col max-w-max ">
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
