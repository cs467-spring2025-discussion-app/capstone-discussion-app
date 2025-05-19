import { Link } from "react-router-dom";

const HomeLink = () => {
  return (
    <Link
      to="/"
      className="
          text-purple-700
          bg-purple-100
          rounded-lg
          max-w-max
          p-1
        "
    >
      Home
    </Link>
  );
};

export default HomeLink;
