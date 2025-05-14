import React from "react";

const CredentialsForm = ({ formData, handleFormChange }) => {
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
          value={formData.email}
          onChange={handleFormChange}
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
          value={formData.password}
          onChange={handleFormChange}
          className="border border-black rounded-md p-2"
        />
      </div>
    </div>
  );
};

export default CredentialsForm;
