import DiscussionPage from '@/pages/DiscussionPage'
import Home from '@/pages/Home'
import AuthMenu from '@/components/auth/AuthMenu'
import AuthProvider from '@/contexts/AuthProvider'

export default function App() {
  return (
    <>
      <AuthProvider>
        <AuthMenu />
        <Home />
        <DiscussionPage />
      </AuthProvider>
    </>
  )
}
