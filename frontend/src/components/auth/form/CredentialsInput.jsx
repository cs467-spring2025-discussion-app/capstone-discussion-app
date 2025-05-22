import { useUser } from "@/hooks/useUser";
import FormButton from "@/components/auth/form/FormButton";
import * as authService from "@/api/authService";

const CredentialsInput = ({ action }) => {
  const {
    authError,
    setAuthError,
    setUser,
    credentials,
    setCredentials,
    setDropdown,
  } = useUser();

  // Capitalize action for label
  const label = `${action.charAt(0).toUpperCase()}${action.slice(1)}`;

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
      setAuthError("Please enter a valid email address");
      return;
    }

    try {
      const userData = await authService.Login(credentials);
      setUser(userData);
      closeAuthForm();
    } catch (error) {
      setAuthError(error.message || "Login failed");
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
      setAuthError(error.message || "Registration failed");
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
      setAuthError(error.message || "Credentials update failed");
    }
  }

  /**
   * Close the current auth box and reset form data
   */
  function closeAuthForm() {
    const modal = document.getElementById(`${action}Form`);
    if (modal) {
      setCredentials({ email: "", password: "" });
      setDropdown(false);
      setAuthError("");
      modal.close();
    } else {
      console.error(`could not get element with id ${action}Form`);
    }
  }

  // Set handle function and label based on action
  const handleSubmit =
    (action === "login" && handleLogin) ||
    (action === "update" && handleUpdate);

  function handleCredentialsChange(e) {
    const { name, value } = e.target;
    setCredentials((formValueState) => ({
      ...formValueState,
      [name]: value,
    }));
  }

  return (
    <div className="flex flex-col items-center">
      <div>
        <div className="flex flex-col">
          <label htmlFor="email" className="text-black">
            Email:
          </label>
          <input
            type="email"
            id="email"
            name="email"
            required
            value={credentials.email}
            onChange={handleCredentialsChange}
            className="text-gray-900 border border-purple-400 rounded-md p-2 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
          />
        </div>
        <div className="flex flex-col">
          <label htmlFor="password">Password:</label>
          <input
            type="password"
            id="password"
            name="password"
            required
            value={credentials.password}
            onChange={handleCredentialsChange}
            className="text-gray-900 border border-purple-400 rounded-md p-2 placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-purple-500"
          />
        </div>
      </div>

      {/* Error message display */}
      {authError && <div className="text-red-600 text-center font-semibold">{authError}</div>}

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

export default CredentialsInput;
