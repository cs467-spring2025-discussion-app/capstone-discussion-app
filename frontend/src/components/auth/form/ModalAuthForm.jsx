import { useEffect } from "react";
import { useUser } from "@/hooks/useUser";
import CredentialsBox from "@/components/auth/form/CredentialsBox";
import DeleteBox from "@/components/auth/form/DeleteBox";

const ModalAuthForm = ({ action }) => {
  const formId = `${action}Form`;
  const { authModalRef, setCredentials } = useUser();

  // Clear credentials on escape press (which closes dialog by default)
  useEffect(() => {
    const handleKeyDown = (e) => {
      if (e.key === "Escape") {
        setCredentials({ email: "", password: "" });
      }
    };
    document.addEventListener("keydown", handleKeyDown);

    return () => {
      document.removeEventListener("keydown", handleKeyDown);
    };
  }, [authModalRef, setCredentials]);

  return (
    <dialog
      ref={authModalRef}
      id={formId}
      className="rounded-lg shadow-2xl border-white border-2 top-1/2 left-1/2 -translate-1/2"
    >
      <div className="flex flex-col bg-purple-100">
        <div className="font-bold rounded-lg p-4 mb-4">
          <div className="flex flex-col">
            {action === "delete" ? (
              <DeleteBox />
            ) : (
              <CredentialsBox action={action} />
            )}
          </div>
        </div>
      </div>
    </dialog>
  );
};

export default ModalAuthForm;
