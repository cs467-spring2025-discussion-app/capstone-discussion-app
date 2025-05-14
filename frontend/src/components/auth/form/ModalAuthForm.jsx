import { useUser } from "@/hooks/useUser";
import { useState } from "react";
import CredentialsInput from "@/components/auth/form/CredentialsInput";
import FormButton from "@/components/auth/form/FormButton";
import * as authService from "@/api/authService";

const ModalAuthForm = ({ action }) => {
  const formId = `${action}Form`;
  const { setUser } = useUser();

  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [confirmDelete, setConfirmDelete] = useState(false);

  const [errorMsg, setAuthErrorMsg] = useState("");

  const handleCredentialsFormChange = (e) => {
    const { name, value } = e.target;
    setFormData((formValueState) => ({
      ...formValueState,
      [name]: value,
    }));
  };
  const handleConfirmDelete = (e) => {
    setConfirmDelete(e.target.checked);
  };

  const handleLogin = async (e) => {
    if (e) {
      e.preventDefault();
    } else {
      return;
    }

    if (!document.getElementById("email").checkValidity()) {
      setAuthErrorMsg("Please enter a valid email address");
      return;
    }

    try {
      const userData = await authService.Login(formData);
      setUser(userData);
      closeAuthBox()
    } catch (error) {
      setAuthErrorMsg(error.message || "Login failed");
    }
  };

  const handleRegister = async (e) => {
    if (e) {
      e.preventDefault();
    } else {
      return;
    }
    try {
      const userData = await authService.Register(formData);
      setUser(userData);
      closeAuthBox()
    } catch (error) {
      setAuthErrorMsg(error.message || "Registration failed");
    }
  };

  const handleUpdate = async (e) => {
    if (e) {
      e.preventDefault();
    } else {
      return;
    }
    try {
      const userData = await authService.UpdateUser(formData);
      setUser(userData);
      closeAuthBox()
    } catch (error) {
      setAuthErrorMsg(error.message || "Credentials update failed");
    }
  };

  const handleDelete = async (e) => {
    if (e) {
      e.preventDefault();
    } else {
      return;
    }
    try {
      await authService.DeleteAccount();
      setUser(null);
      closeAuthBox()
    } catch (error) {
      setAuthErrorMsg(error.message || "Account delete failed");
    }
  };

  const closeAuthBox = () => {
    setFormData({
      email: "",
      password: "",
    });
    setAuthErrorMsg("");
    document.getElementById(formId).close();
  };

  // // Set handle function and label based on action
  const handleSubmit =
    (action === "login" && handleLogin) ||
    (action === "register" && handleRegister) ||
    (action === "update" && handleUpdate) ||
    (action === "delete" && handleDelete);

  //   action === "register" ? handleRegister;

  // Capitalize action for label
  const label = `${action.charAt(0).toUpperCase()}${action.slice(1)}`;

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
              <div className="flex-col items-center text-lg">
                <div className="flex justify-center">Delete account?</div>
                <div className="mt-4 flex items-center">
                  <input
                    type="checkbox"
                    id="confirmDeleteCheckbox"
                    checked={confirmDelete}
                    onChange={handleConfirmDelete}
                    className="mr-2"
                  />
                  <label htmlFor="confirmDeleteCheckbox">
                    Check to confirm. This cannot be undone!!
                  </label>
                </div>
              </div>
            ) : (
              // Email and password input
              <CredentialsInput
                formData={formData}
                handleChange={handleCredentialsFormChange}
              />
            )}

            {/* Error message display */}
            {errorMsg && <div className="text-red-600">{errorMsg}</div>}

            {/* Action buttons */}
            <div className="flex gap-4 justify-center mt-4">
              {action === "delete" ? (
                <button
                  type="button"
                  onClick={handleSubmit}
                  disabled={!confirmDelete}
                  className={`py-2 px-4 rounded-md text-white
                  ${
                    !confirmDelete
                      ? "cursor-not-allowed bg-red-300"
                      : "cursor-pointer bg-red-500 hover:bg-red-600"
                  }`}
                >
                  {label}
                </button>
              ) : (
                <FormButton text={label} func={handleSubmit} />
              )}
              <FormButton text="Close" func={closeAuthBox} />
            </div>
          </div>
        </div>
      </div>
    </dialog>
  );
};

export default ModalAuthForm;
