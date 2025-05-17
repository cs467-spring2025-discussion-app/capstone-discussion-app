import { useState } from "react";
import AuthContext from "./AuthContext";

const AuthProvider = ({ children }) => {
  const [user, setUser] = useState(null);
  const [credentials, setCredentials] = useState({email: "", password: ""});

  return (
    <AuthContext.Provider value={{
      user, setUser,
      credentials, setCredentials,
    }}>
      {children}
    </AuthContext.Provider>
  );
};

export default AuthProvider;
