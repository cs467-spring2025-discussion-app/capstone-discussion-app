const FormButton = ({ text, func }) => {
  return (
    <button
      type="button"
      onClick={func}
      className="bg-purple-600 text-white p-2 rounded-md hover:bg-purple-700 cursor-pointer transition-colors duration-200"
    >
      {text}
    </button>
  );
};

export default FormButton;
