import { useUser } from "@/hooks/useUser";
import { useState } from "react";
import CredentialsInput from "@/components/auth/form/CredentialsInput";
import FormButton from "@/components/auth/form/FormButton";
import * as authService from "@/api/authService";

const ModalAuthForm = ({action}) => {
  const formId = `${action}Form`;
  const { setUser } = useUser();

  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const [errorMsg, setAuthErrorMsg] = useState("");

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData((formValueState) => ({
      ...formValueState,
      [name]: value,
    }));
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
      document.getElementById(formId).close();
      window.location.reload();  // Reload to set auth
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
      document.getElementById(formId).close();
    } catch (error) {
      setAuthErrorMsg(error.message || "Registration failed");
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

  // Set handle function and label based on action
  const handleSubmit = action === "login" ? handleLogin : handleRegister;
  // Capitalize action for label
  const label = `${action.charAt(0).toUpperCase()}${action.slice(1)}`

  return (
    <dialog
      id={formId}
      className="rounded-lg shadow-2xl border-white border-2 top-1/2 left-1/2 -translate-1/2"
    >
      <div className="flex flex-col bg-purple-100">
        <div className="font-bold rounded-lg p-4 mb-4">
          <div className="flex flex-col">
            {/* Email and password input */}
            <CredentialsInput formData={formData} handleChange={handleChange} />

            {/* Error message display */}
            {errorMsg && <div className="text-red-600">{errorMsg}</div>}

            {/* Action buttons */}
            <div className="flex gap-4 justify-center mt-4">
                <FormButton text={label} func={handleSubmit} />
                <FormButton text="Close" func={closeAuthBox} />
            </div>
          </div>
        </div>
      </div>
    </dialog>
  );
};

export default ModalAuthForm;
