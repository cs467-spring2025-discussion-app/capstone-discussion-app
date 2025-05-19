import { BrowserRouter, Routes, Route } from "react-router-dom";
import DiscussionPage from "@/pages/DiscussionPage";
import Home from "@/pages/Home";
import NotFound from "@/pages/NotFound";
import AuthMenu from "@/components/auth/AuthMenu";
import AuthProvider from "@/contexts/AuthProvider";

export default function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <AuthMenu />
        <Routes>
          <Route path="/" element={<Home />} />
          <Route
            path="/discussions/:discussionId"
            element={<DiscussionPage />}
          />

          <Route path="/not-found" element={<NotFound />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}
