import { useState } from "react";
import { useUser } from "@/hooks/useUser";
import CredentialsInput from "@/components/auth/form/CredentialsInput";
import FormButton from "@/components/auth/form/FormButton";
import * as authService from "@/api/authService";

const CredentialsBox = ({ action }) => {
  const formId = `${action}Form`;
  const { setUser } = useUser();

  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });

  const [errorMsg, setAuthErrorMsg] = useState("");

  /**
   * Update formData state on form input change
   */
  function handleCredentialsFormChange(e) {
    const { name, value } = e.target;
    setFormData((formValueState) => ({
      ...formValueState,
      [name]: value,
    }));
  }

  /**
   * Handle login button press
   */
  async function handleLogin(e) {
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
      closeAuthBox();
    } catch (error) {
      setAuthErrorMsg(error.message || "Login failed");
    }
  }

  /**
   * Handle login button press
   */
  async function handleRegister(e) {
    if (e) {
      e.preventDefault();
    } else {
      return;
    }
    try {
      const userData = await authService.Register(formData);
      setUser(userData);
      closeAuthBox();
    } catch (error) {
      setAuthErrorMsg(error.message || "Registration failed");
    }
  }

  /**
   * Handle update button press
   */
  async function handleUpdate(e) {
    if (e) {
      e.preventDefault();
    } else {
      return;
    }
    try {
      const userData = await authService.UpdateUser(formData);
      setUser(userData);
      closeAuthBox();
    } catch (error) {
      setAuthErrorMsg(error.message || "Credentials update failed");
    }
  }

  /**
   * Close the current auth box and reset form data
   */
  function closeAuthBox() {
    setFormData({
      email: "",
      password: "",
    });
    setAuthErrorMsg("");
    document.getElementById(formId).close();
  }

  // Set handle function and label based on action
  const handleSubmit =
    (action === "login" && handleLogin) ||
    (action === "update" && handleUpdate);

  // Capitalize action for label
  const label = `${action.charAt(0).toUpperCase()}${action.slice(1)}`;

  return (
    <div className="flex flex-col items-center">
      <CredentialsInput
        formData={formData}
        handleFormChange={handleCredentialsFormChange}
      />

      {/* Error message display */}
      {errorMsg && <div className="text-red-600">{errorMsg}</div>}

      {/* Action buttons */}
      <div className="flex gap-4 justify-center mt-4">
        <FormButton text={label} func={handleSubmit} />

        {/* Show additional register button for login */}
        {action === "login" && (
          <FormButton text="Register" func={handleRegister} />
        )}

        <FormButton text="Close" func={closeAuthBox} />
      </div>
    </div>
  );
};

export default CredentialsBox;
