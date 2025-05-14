import DiscussionPage from '@/pages/DiscussionPage'
import Home from '@/pages/Home'
import AuthControl from '@/components/auth/AuthControl'
import AuthProvider from '@/contexts/AuthProvider'

export default function App() {
  return (
    <>
      <AuthProvider>
        <AuthControl />
        <Home />
        <DiscussionPage />
      </AuthProvider>
    </>
  )
}
