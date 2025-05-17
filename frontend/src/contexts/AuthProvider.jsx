import { useState } from "react";
import AuthContext from "./AuthContext";

const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [credentials, setCredentials] = useState({email: "", password: ""});
  const [dropdown, setDropdown] = useState(null);

  return (
    <AuthContext.Provider value={{
      user, setUser,
      credentials, setCredentials,
      dropdown, setDropdown,
    }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
