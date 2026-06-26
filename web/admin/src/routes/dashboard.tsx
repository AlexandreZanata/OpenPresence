import type { CSSProperties } from 'react'
import { createFileRoute, Link, useRouteContext } from '@tanstack/react-router'

export const Route = createFileRoute('/dashboard')({
  component: DashboardPage,
})

function DashboardPage() {
  const { auth } = useRouteContext({ from: '__root__' })

  return (
    <main style={styles.main}>
      <h1 style={styles.title}>Dashboard</h1>
      <p style={styles.lead}>
        {auth.isAuthenticated
          ? `Welcome, ${auth.user?.displayName ?? 'user'}.`
          : 'You are not signed in.'}
      </p>
      <p style={styles.hint}>Admin shell layout arrives in admin-05.</p>
      <Link to="/login">Back to login</Link>
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
}
