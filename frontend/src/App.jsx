import { BrowserRouter, Routes, Route, useLocation } from "react-router-dom";
import DiscussionPage from "@/pages/DiscussionPage";
import Home from "@/pages/Home";
import HomeLink from "@/components/HomeLink";
import NotFound from "@/pages/NotFound";
import AuthMenu from "@/components/auth/AuthMenu";
import AuthProvider from "@/contexts/AuthProvider";

export default function App() {
  function AppLayout() {
    const location = useLocation();

    return (
      <>
        <div className="flex m-2">
          {location.pathname !== "/" && <HomeLink />}
          <div className="flex-grow"></div>
          <AuthMenu />
        </div>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/discussions/:discussionId" element={<DiscussionPage />} />
          <Route path="/not-found" element={<NotFound />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </>
    );
  }

  return (
    <AuthProvider>
      <BrowserRouter>
        <AppLayout />
      </BrowserRouter>
    </AuthProvider>
  );
}
