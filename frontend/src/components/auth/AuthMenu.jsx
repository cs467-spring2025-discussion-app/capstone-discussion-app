import { useUser } from "@/hooks/useUser";
import { useEffect, useState } from "react";
import AuthButton from "@/components/auth/AuthButton";
import ModalAuthForm from "@/components/auth/form/ModalAuthForm";
import * as authService from "@/api/authService"

const AuthMenu = () => {
  const [isCheckingStatus, setIsCheckingStatus] = useState(true);
  const { user, setUser } = useUser();

  useEffect(() => {
    const checkLoginStatus = async () => {
      try {
        const userData = await authService.Whoami();
        setUser(userData)
      } catch (error) {
        console.error("Error checking login status:", error);
        setUser(null)
      } finally {
        setIsCheckingStatus(false);
      }
    };

    checkLoginStatus();
  }, [setUser]);

  return (
    <div>
      {isCheckingStatus ? (
        <span className="inline-block p-1">&nbsp;</span>
      ) : !user ? (
        <div>
          <AuthButton action="login"/>
          <AuthButton action="register"/>
        </div>
      ) : (
        <AuthButton action="logout"/>
      )}

      <ModalAuthForm action="login" />
      <ModalAuthForm action="register" />
    </div>
  );
};

export default AuthMenu;
