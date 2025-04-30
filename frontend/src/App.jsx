import DiscussionPage from '@/pages/DiscussionPage'
import Home from '@/pages/Home'
import AuthButton from '@/components/AuthButton'

export default function App() {
  return (
    <>
      <AuthButton />
      <Home />
      <DiscussionPage />
    </>
  )
}
