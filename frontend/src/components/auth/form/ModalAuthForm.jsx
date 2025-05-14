import CredentialsBox from "@/components/auth/form/CredentialsBox";
import DeleteBox from "@/components/auth/form/DeleteBox";

const ModalAuthForm = ({ action }) => {
  const formId = `${action}Form`;

  return (
    <dialog
      id={formId}
      className="rounded-lg shadow-2xl border-white border-2 top-1/2 left-1/2 -translate-1/2"
    >
      <div className="flex flex-col bg-purple-100">
        <div className="font-bold rounded-lg p-4 mb-4">
          <div className="flex flex-col">
            {action === "delete" ? (
              // Delete account with confirmation
              <DeleteBox />
            ) : (
              // Email and password input (Register, Login, Update)
              <CredentialsBox action={action} />
            )}
          </div>
        </div>
      </div>
    </dialog>
  );
};

export default ModalAuthForm;
