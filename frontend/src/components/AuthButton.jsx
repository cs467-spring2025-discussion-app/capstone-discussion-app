import { useEffect, useState } from "react";
import { useUser } from "@/hooks/useUser";

const AuthButton = () => {
  const [formData, setFormData] = useState({
    email: "",
    password: "",
  });
  const [isCheckingStatus, setIsCheckingStatus] = useState(true);
  const { user, setUser } = useUser();
  const [errorMsg, setErrorMsg] = useState("");

  useEffect(() => {
    const checkLoginStatus = async () => {
      try {
        const response = await fetch("/authapi/whoami", {
          method: "GET",
          headers: {
            "Content-Type": "application/json",
          },
        });
        if (response.ok) {
          const userData = await response.json();
          setUser(userData);
        } else {
          setUser(null);
        }
      } catch (error) {
        console.error("Error checking login status:", error);
      } finally {
        setIsCheckingStatus(false);
      }
    };

    checkLoginStatus();
  }, [setUser]);

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
      setErrorMsg("Please enter a valid email address");
      return;
    }

    try {
      const response = await fetch("/authapi/login", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      console.log("Login Status:", response.status);
      console.log("Login Email:", formData.email);
      console.log("Response", response);

      if (!response.ok) {
        const resp = await response.json();
        setErrorMsg(resp.error || "Login failed");
        console.error("Login request rejected:", response.status, resp);
        return;
      }

      const userData = await fetch("/authapi/whoami").then((res) => res.json());
      setUser(userData);
      document.getElementById("loginForm").close();
      window.location.reload();
    } catch (error) {
      setErrorMsg("Login failed");
      console.error("Login error:", error);
    }
  };

  const handleRegister = async (e) => {
    if (e) {
      e.preventDefault();
    } else {
      return;
    }
    try {
      const response = await fetch("/authapi/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(formData),
      });

      console.log("Register status:", response.status);
      console.log("Register email:", formData.email);
      console.log("Register response:", response);

      if (!response.ok) {
        setErrorMsg("Registration failed");
        const resp = await response.json();
        setErrorMsg(resp.error || "Registration failed");
        console.error("Registration request rejected:", response.status, resp);
        return;
      }

      // NOTE: Reload and force new login - but eventually want email confirmation
      window.location.reload();
    } catch (error) {
      setErrorMsg("Registration failed");
      console.error("Registration error:", error);
    }
  };

  const handleLogout = async () => {
    try {
      await fetch("/authapi/logout", {
        method: "POST",
      });
      setUser(null);
      window.location.reload();
    } catch (error) {
      console.error("Logout error:", error);
    }
  };

  const closeAuthBox = () => {
    setFormData({
      email: "",
      password: "",
    });
    setErrorMsg("");
    document.getElementById("loginForm").close();
  };

  return (
    <div>
      {isCheckingStatus ? (
        <span className="inline-block p-1">&nbsp;</span>
      ) : !user ? (
        <button
          type="button"
          className="text-purple-800 bg-purple-100 inline-block p-1 rounded-lg cursor-pointer"
          onClick={() => document.getElementById("loginForm").showModal()}
        >
          Login
        </button>
      ) : (
        <button
          type="button"
          className="text-purple-800 bg-purple-100 inline-block p-1 rounded-lg cursor-pointer"
          onClick={handleLogout}
        >
          Logout
        </button>
      )}

      <dialog
        id="loginForm"
        className="rounded-lg shadow-2xl border-white border-2 top-1/2 left-1/2 -translate-1/2"
      >
        <div className="flex flex-col bg-purple-100">
          <div className="font-bold rounded-lg p-4 mb-4">
            <div className="flex flex-col">
              <div className="flex flex-col">
                <label htmlFor="email" className="text-black">
                  Email:
                </label>
                <input
                  type="email"
                  id="email"
                  name="email"
                  required
                  value={formData.email}
                  onChange={handleChange}
                  className="text-black border border-black rounded-md p-2"
                />
              </div>
              <div className="flex flex-col">
                <label htmlFor="password">Password:</label>
                <input
                  type="password"
                  id="password"
                  name="password"
                  value={formData.password}
                  onChange={handleChange}
                  required
                  className="border border-black rounded-md p-2"
                />
              </div>
              {errorMsg && <div className="text-red-600">{errorMsg}</div>}
              <div className="flex justify-center space-x-4 pt-4">
                <button
                  type="submit"
                  onClick={handleLogin}
                  className="bg-black text-white p-1 rounded-md hover:bg-blue-500 cursor-pointer"
                >
                  Login
                </button>
                <button
                  type="submit"
                  onClick={handleRegister}
                  className="bg-black text-white p-1 rounded-md hover:bg-blue-500 cursor-pointer"
                >
                  Register
                </button>
              </div>
            </div>

            <div className="flex justify-center pt-2">
              <button
                type="button"
                onClick={closeAuthBox}
                className="bg-black text-white p-1 rounded-md hover:bg-blue-500 cursor-pointer"
              >
                Close
              </button>
            </div>
          </div>
        </div>
      </dialog>
    </div>
  );
};

export default AuthButton;
