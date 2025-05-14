const FormButton = ({ text, func }) => {
  return (
    <button
      type="button"
      onClick={func}
      className="bg-black text-white p-1 rounded-md hover:bg-blue-500 cursor-pointer"
    >
      {text}
    </button>
  );
};

export default FormButton;
