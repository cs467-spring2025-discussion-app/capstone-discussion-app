import { useState } from "react";
import { useUser } from "@/hooks/useUser";
import FormButton from "@/components/auth/form/FormButton";
import * as authService from "@/api/authService";

const DeleteBox = () => {
  const { setUser } = useUser();
  const [confirmDelete, setConfirmDelete] = useState(false);

  const handleConfirmDelete = (e) => {
    setConfirmDelete(e.target.checked);
  };

  const closeDeleteBox = () => {
    document.getElementById("deleteForm").close();
  };

  /**
   * Handle delete button press
   */
  async function handleDelete(e) {
    if (e) {
      e.preventDefault();
    } else {
      return;
    }
    try {
      await authService.DeleteAccount();
      setUser(null);
      closeDeleteBox();
    } catch (error) {
      console.log(error);
    }
  }

  return (
    <div>
      <div className="flex-col items-center text-lg">
        <div className="flex justify-center">Delete account?</div>
        <div className="mt-4 flex items-center">
          <input
            type="checkbox"
            id="confirmDeleteCheckbox"
            checked={confirmDelete}
            onChange={handleConfirmDelete}
            className="mr-2"
          />
          <label htmlFor="confirmDeleteCheckbox">
            Check to confirm. This cannot be undone!!
          </label>
        </div>
      </div>

      <div className="flex gap-4 justify-center mt-4">
        <button
          type="button"
          onClick={handleDelete}
          disabled={!confirmDelete}
          className={`py-2 px-4 rounded-md text-white
                  ${
                    !confirmDelete
                      ? "cursor-not-allowed bg-red-300"
                      : "cursor-pointer bg-red-500 hover:bg-red-600"
                  }`}
        >
          Delete
        </button>
        <FormButton text="Close" func={closeDeleteBox} />
      </div>
    </div>
  );
};

export default DeleteBox;
