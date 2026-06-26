import { useState, type CSSProperties, type FormEvent } from 'react'
import { createFileRoute, useRouteContext } from '@tanstack/react-router'
import { apiBaseUrl, isAuthMockEnabled } from '~/lib/env'
import { ApiError } from '~/lib/api/client'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  const { auth } = useRouteContext({ from: '__root__' })
  const [registrationId, setRegistrationId] = useState('admin')
  const [password, setPassword] = useState('admin')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)

  async function onLogin(event: FormEvent) {
    event.preventDefault()
    setBusy(true)
    setError(null)
    try {
      await auth.login({ registrationId, password })
    } catch (err) {
      setError(err instanceof ApiError ? err.message : 'Login failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <main style={styles.main}>
      <h1 style={styles.title}>OpenPresence Admin</h1>
      <p style={styles.lead}>Auth core — mock login until POST /v1/auth/login exists.</p>

      <dl style={styles.dl}>
        <dt>API base URL</dt>
        <dd>
          <code>{apiBaseUrl()}</code>
        </dd>
        <dt>Auth mock</dt>
        <dd>
          <code>{String(isAuthMockEnabled())}</code>
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

      {!auth.isAuthenticated ? (
        <form onSubmit={onLogin} style={styles.form}>
          <label style={styles.label}>
            Registration ID
            <input
              value={registrationId}
              onChange={(e) => setRegistrationId(e.target.value)}
              autoComplete="username"
              style={styles.input}
            />
          </label>
          <label style={styles.label}>
            Password
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              autoComplete="current-password"
              style={styles.input}
            />
          </label>
          {error ? <p style={styles.error}>{error}</p> : null}
          <button type="submit" disabled={busy} style={styles.button}>
            {busy ? 'Signing in…' : 'Sign in (mock)'}
          </button>
        </form>
      ) : (
        <button type="button" onClick={() => auth.logout()} style={styles.button}>
          Sign out
        </button>
      )}
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
  dl: { background: '#f4f4f5', padding: '1rem', borderRadius: 8, marginBottom: '1.5rem' },
  form: { display: 'grid', gap: '0.75rem' },
  label: { display: 'grid', gap: '0.25rem', fontSize: '0.9rem' },
  input: { padding: '0.5rem', fontSize: '1rem' },
  button: {
    width: 'fit-content',
    padding: '0.5rem 1rem',
    fontSize: '1rem',
    cursor: 'pointer',
  },
  error: { color: '#b91c1c', margin: 0 },
}
