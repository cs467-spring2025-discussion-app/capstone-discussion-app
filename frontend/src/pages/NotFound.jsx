import  HomeLink  from "@/components/HomeLink"

export default function NotFound() {
  return (
    <div
      className="
        flex flex-col items-center
        gap-4
      "
    >
      <h1>404 - Not Found</h1>
      <HomeLink />
    </div>
  );
}
