import { useRef, useState } from "react";
import AuthContext from "./AuthContext";

const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [credentials, setCredentials] = useState({ email: "", password: "" });
  const [dropdown, setDropdown] = useState(null);
  const [authError, setAuthError] = useState(null);
  const authModalRef = useRef(null);

  return (
    <AuthContext.Provider
      value={{
        authModalRef,
        authError,
        setAuthError,
        credentials,
        setCredentials,
        dropdown,
        setDropdown,
        user,
        setUser,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
