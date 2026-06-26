import type { CSSProperties } from 'react'
import { Link, createFileRoute, useRouteContext } from '@tanstack/react-router'
import { apiBaseUrl } from '~/lib/env'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  const { auth } = useRouteContext({ from: '__root__' })

  return (
    <main style={styles.main}>
      <h1 style={styles.title}>OpenPresence Admin</h1>
      <p style={styles.lead}>Administrative panel — TanStack Start.</p>
      <dl style={styles.dl}>
        <dt>API</dt>
        <dd>
          <code>{apiBaseUrl()}</code>
        </dd>
        <dt>Session</dt>
        <dd>
          {auth.isLoading
            ? 'Loading…'
            : auth.isAuthenticated
              ? `${auth.user?.displayName} (${auth.user?.roles.join(', ')})`
              : 'Signed out'}
        </dd>
      </dl>
      <p>
        {auth.isAuthenticated ? (
          <Link to="/dashboard">Go to dashboard</Link>
        ) : (
          <Link to="/login">Sign in</Link>
        )}
      </p>
    </main>
  )
}

const styles: Record<string, CSSProperties> = {
  main: {
    fontFamily: 'system-ui, sans-serif',
    maxWidth: 640,
    margin: '4rem auto',
    padding: '0 1.5rem',
    lineHeight: 1.5,
  },
  title: { fontSize: '2rem', marginBottom: '0.5rem' },
  lead: { color: '#444', marginBottom: '1.5rem' },
  dl: { background: '#f4f4f5', padding: '1rem', borderRadius: 8, marginBottom: '1rem' },
}
