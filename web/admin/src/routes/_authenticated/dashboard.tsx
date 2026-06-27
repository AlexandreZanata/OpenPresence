import type { CSSProperties } from 'react'
import { useQuery } from '@tanstack/react-query'
import { createFileRoute, useRouteContext } from '@tanstack/react-router'
import { fetchLiveHealth } from '~/lib/api/health'
import { apiBaseUrl } from '~/lib/env'

export const Route = createFileRoute('/_authenticated/dashboard')({
  component: DashboardPage,
})

function DashboardPage() {
  const { auth } = useRouteContext({ from: '__root__' })
  const healthQuery = useQuery({
    queryKey: ['health', 'live'],
    queryFn: fetchLiveHealth,
    retry: 1,
    refetchInterval: 30_000,
  })

  const apiStatus = healthQuery.isLoading
    ? 'Checking…'
    : healthQuery.isError
      ? 'Offline'
      : 'Online'

  return (
    <div style={styles.page}>
      <h1 style={styles.title}>Dashboard</h1>
      <div style={styles.card}>
        <h2 style={styles.cardTitle}>Welcome back</h2>
        <p style={styles.lead}>
          Signed in as <strong>{auth.user?.displayName ?? 'user'}</strong> (
          {auth.user?.roles.join(', ') ?? 'no roles'}).
        </p>
        <p style={styles.hint}>Employee, enrollment, and punch modules will appear here in later phases.</p>
      </div>
      <div style={styles.card}>
        <h2 style={styles.cardTitle}>API status</h2>
        <dl style={styles.dl}>
          <dt>Base URL</dt>
          <dd>
            <code>{apiBaseUrl()}</code>
          </dd>
          <dt>Health</dt>
          <dd>
            <span style={statusStyle(healthQuery.isError)}>{apiStatus}</span>
            {healthQuery.data?.status ? ` (${healthQuery.data.status})` : null}
          </dd>
        </dl>
        <button type="button" style={styles.button} onClick={() => void healthQuery.refetch()}>
          Refresh status
        </button>
      </div>
    </div>
  )
}

function statusStyle(isError: boolean): CSSProperties {
  return {
    fontWeight: 600,
    color: isError ? '#b91c1c' : '#15803d',
  }
}

const styles: Record<string, CSSProperties> = {
  page: { maxWidth: 720 },
  title: { fontSize: '1.75rem', marginBottom: '1rem' },
  card: {
    background: '#fff',
    border: '1px solid #e2e8f0',
    borderRadius: 10,
    padding: '1.25rem',
    marginBottom: '1rem',
  },
  cardTitle: { fontSize: '1.1rem', marginBottom: '0.75rem' },
  lead: { color: '#334155', marginBottom: '0.5rem' },
  hint: { color: '#64748b', fontSize: '0.9rem' },
  dl: {
    display: 'grid',
    gridTemplateColumns: 'auto 1fr',
    gap: '0.35rem 1rem',
    marginBottom: '1rem',
  },
  button: {
    cursor: 'pointer',
    padding: '0.4rem 0.85rem',
    borderRadius: 6,
    border: '1px solid #cbd5e1',
    background: '#fff',
  },
}
