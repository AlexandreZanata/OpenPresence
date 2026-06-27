import { useEffect, useState } from 'react'
import {
  createFileRoute,
  redirect,
  useNavigate,
  useRouteContext,
  useRouter,
} from '@tanstack/react-router'
import { LoginForm } from '~/components/login/LoginForm'
import { loginPageStyles as styles } from '~/components/login/form-styles'

type LoginSearch = {
  redirect?: string
}

export const Route = createFileRoute('/login')({
  validateSearch: (search: Record<string, unknown>): LoginSearch => ({
    redirect: typeof search.redirect === 'string' ? search.redirect : undefined,
  }),
  beforeLoad: ({ context, search }) => {
    if (context.auth.isAuthenticated) {
      throw redirect({ to: search.redirect ?? '/dashboard' })
    }
  },
  component: LoginPage,
})

function LoginPage() {
  const { auth } = useRouteContext({ from: '__root__' })
  const { redirect: redirectTo } = Route.useSearch()
  const navigate = useNavigate()
  const router = useRouter()
  const [hydrated, setHydrated] = useState(false)

  useEffect(() => {
    setHydrated(true)
  }, [])

  function onSuccess() {
    const target = redirectTo ?? '/dashboard'
    if (target.startsWith('http')) {
      router.history.push(target)
      return
    }
    void navigate({ to: target })
  }

  return (
    <div style={styles.page}>
      <div style={styles.card}>
        <h1 style={styles.title}>OpenPresence Admin</h1>
        <p style={styles.subtitle}>Sign in with your registration ID</p>
        {!hydrated ? (
          <p>Loading session…</p>
        ) : (
          <LoginForm auth={auth} onSuccess={onSuccess} />
        )}
      </div>
    </div>
  )
}
