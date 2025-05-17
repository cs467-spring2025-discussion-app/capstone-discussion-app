import { useState } from "react";
import { useUser } from "@/hooks/useUser";
import CredentialsInput from "@/components/auth/form/CredentialsInput";
import FormButton from "@/components/auth/form/FormButton";
import * as authService from "@/api/authService";

const CredentialsBox = ({ action }) => {
  const { setUser, credentials, setCredentials, setDropdown } = useUser();

  const [errorMsg, setAuthErrorMsg] = useState("");

  /**
   * Update credentials state on form input change
   */
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
      const userData = await authService.Login(credentials);
      setUser(userData);
      closeAuthForm();
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
      const userData = await authService.Register(credentials);
      setUser(userData);
      closeAuthForm();
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
      const userData = await authService.UpdateUser(credentials);
      setUser(userData);
      closeAuthForm();
    } catch (error) {
      setAuthErrorMsg(error.message || "Credentials update failed");
    }
  }

  /**
   * Close the current auth box and reset form data
   */
  function closeAuthForm() {
    const modal = document.getElementById(`${action}Form`);
    if (modal) {
      setCredentials({email: "", password: ""})
      setDropdown(false)
      modal.close();
    } else {
      console.error(`could not get element with id ${action}Form`)
    }
  }

  // Set handle function and label based on action
  const handleSubmit =
    (action === "login" && handleLogin) ||
    (action === "update" && handleUpdate);

  // Capitalize action for label
  const label = `${action.charAt(0).toUpperCase()}${action.slice(1)}`;

  return (
    <div className="flex flex-col items-center">
      <CredentialsInput />

      {/* Error message display */}
      {errorMsg && <div className="text-red-600 text-center">{errorMsg}</div>}

      {/* Action buttons */}
      <div className="flex gap-4 justify-center mt-4">
        <FormButton text={label} func={handleSubmit} />

        {/* Show additional register button for login */}
        {action === "login" && (
          <FormButton text="Register" func={handleRegister} />
        )}

        <FormButton text="Close" func={closeAuthForm} />
      </div>
    </div>
  );
};

export default CredentialsBox;
