import { useUser } from "@/hooks/useUser";

const LogoutButton = () => {
  const { setUser } = useUser();

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

  return (
    <button
      type="button"
      className="text-purple-800 bg-purple-100 inline-block p-1 rounded-lg cursor-pointer"
      onClick={handleLogout}
    >
      Logout
    </button>
  );
};

export default LogoutButton;
