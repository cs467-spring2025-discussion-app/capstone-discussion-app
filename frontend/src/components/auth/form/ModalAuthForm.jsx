import { useEffect } from "react";
import { useUser } from "@/hooks/useUser";
import CredentialsInput from "@/components/auth/form/CredentialsInput";
import DeleteBox from "@/components/auth/form/DeleteBox";

const ModalAuthForm = ({ action }) => {
  const formId = `${action}Form`;
  const { authModalRef, setAuthError, setCredentials } = useUser();

  // Clear credentials on escape press (which closes dialog by default)
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === "Escape") {
        setCredentials({ email: "", password: "" });
        setAuthError("")
      }
    };
    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [authModalRef, setAuthError, setCredentials]);

  return (
    <dialog
      ref={authModalRef}
      id={formId}
      className="rounded-lg shadow-lg border border-purple-500 bg-white p-6 fixed top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 max-w-md w-full"
    >
      <div className="flex flex-col bg-purple-100">
        <div className="font-bold rounded-lg p-4 mb-4">
          <div className="flex flex-col">
            {action === "delete" ? (
              <DeleteBox />
            ) : (
              <CredentialsInput action={action} />
            )}
          </div>
        </div>
      </div>
    </dialog>
  );
};

export default ModalAuthForm;
