/**
 * Register user with credentials and login
 * @param {Object} userData - User email and password
 * @returns {Promise<Object>} Response data
 * @throws {Error} On registration failure
 */
export async function Register(credentials) {
  const registerResponse = await fetch("/authapi/register", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });

  if (!registerResponse.ok) {
    const errorData = await registerResponse.json();
    throw new Error(errorData.error || "Registration failed");
  }

  // Try login after registration
  try {
    return await Login(credentials);
  } catch {
    throw new Error("Please try logging in manually.");
  }
}

/**
 * Login user with credentials
 * @param {Object} credentials - User email and password
 * @returns {Promise<Object>} User data
 * @throws {Error} On login failure
 */
export async function Login(credentials) {
  const response = await fetch("/authapi/login", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(credentials),
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || "Login failed");
  }

  // Return credentials
  return await Whoami();
}

/**
 * Update user credentials
 * @param {Object} userData - New user credentials
 * @returns {Promise<Object>} Updated user data
 * @throws {Error} On update failure
 */
export const UpdateUser = async (userData) => {
  const response = await fetch("/authapi/updateuser", {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(userData),
  });

  if (!response.ok) {
    const errorData = await response.json();
    throw new Error(errorData.error || "Update failed");
  }

  // Return new credentials
  return await Whoami();
};

/**
 * Get current user data
 * @returns {Promise<Object|null>} User data or null if not authenticated
 */
export async function Whoami() {
  try {
    const response = await fetch("/authapi/whoami", {
      method: "GET",
      headers: {
        "Content-Type": "application/json",
      },
    });

    if (!response.ok) {
      return null;
    }

    return await response.json();
  } catch (error) {
    console.error("Error fetching current user:", error);
    return null;
  }
}

/**
 * Logout current user
 * @returns {Promise<void>}
 */
export async function Logout() {
  await fetch("/authapi/logout", {
    method: "POST",
  });
}

/**
 * Delete user account permanently
 * @returns {Promise<void>}
 */
export async function DeleteAccount() {
  try {
    const response = await fetch("/authapi/deleteaccount", {
      method: "DELETE",
    });
    return await response.json();
  } catch (error) {
    console.error("Error fetching current user:", error);
    return null;
  }
}
