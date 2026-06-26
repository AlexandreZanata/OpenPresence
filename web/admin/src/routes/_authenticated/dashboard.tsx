import type { CSSProperties } from 'react'
import { createFileRoute, Link, useRouteContext } from '@tanstack/react-router'

export const Route = createFileRoute('/_authenticated/dashboard')({
  component: DashboardPage,
})

function DashboardPage() {
  const { auth } = useRouteContext({ from: '__root__' })

  return (
    <main style={styles.main}>
      <h1 style={styles.title}>Dashboard</h1>
      <p style={styles.lead}>Welcome, {auth.user?.displayName ?? 'user'}.</p>
      <p style={styles.hint}>Admin shell layout arrives in admin-05.</p>
      <p>
        <button type="button" style={styles.button} onClick={() => auth.logout()}>
          Sign out
        </button>
        {' · '}
        <Link to="/login">Back to login</Link>
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
  },
  title: { fontSize: '2rem' },
  lead: { color: '#444' },
  hint: { color: '#666', fontSize: '0.9rem' },
  button: {
    cursor: 'pointer',
    padding: '0.35rem 0.75rem',
    borderRadius: 6,
    border: '1px solid #ccc',
    background: '#fff',
  },
}
