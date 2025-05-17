import React from "react";
import { useUser } from "@/hooks/useUser";

const CredentialsInput = () => {
  const { credentials, setCredentials } = useUser();
  function handleCredentialsChange(e) {
    const { name, value } = e.target;
    setCredentials((formValueState) => ({
      ...formValueState,
      [name]: value,
    }));
  }

  return (
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
          className="text-black border border-black rounded-md p-2"
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
          className="border border-black rounded-md p-2"
        />
      </div>
    </div>
  );
};

export default CredentialsInput;
