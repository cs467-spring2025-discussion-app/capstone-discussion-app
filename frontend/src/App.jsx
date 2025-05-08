import DiscussionPage from '@/pages/DiscussionPage'
import Home from '@/pages/Home'
import AuthButton from '@/components/AuthButton'
import AuthProvider from '@/contexts/AuthProvider'

export default function App() {
  return (
    <>
      <AuthProvider>
        <AuthButton />
        <Home />
        <DiscussionPage />
      </AuthProvider>
    </>
  )
}
