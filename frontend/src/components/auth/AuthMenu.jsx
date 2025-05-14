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
    <div className="m-2">
      {isCheckingStatus ? (
        <span className="inline-block p-1">&nbsp;</span>
      ) : !user ? (
        <div className="flex gap-2 items-center">
          <AuthButton action="login"/>
          <AuthButton action="register"/>
        </div>
      ) : (
        <div className="flex gap-2 items-center">
          <AuthButton action="update"/>
          <AuthButton action="logout"/>
          <AuthButton action="delete"/>
          <span>{user.email}</span>
        </div>
      )}

      <ModalAuthForm action="login" />
      <ModalAuthForm action="register" />
      <ModalAuthForm action="update" />
      <ModalAuthForm action="delete" />
    </div>
  );
};

export default AuthMenu;
