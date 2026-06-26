import type { CSSProperties } from 'react'
import { createFileRoute } from '@tanstack/react-router'
import { apiBaseUrl } from '~/lib/env'

export const Route = createFileRoute('/')({
  component: HomePage,
})

function HomePage() {
  return (
    <main style={styles.main}>
      <h1 style={styles.title}>OpenPresence Admin</h1>
      <p style={styles.lead}>
        Administrative panel scaffold — TanStack Start + Router + Query.
      </p>
      <dl style={styles.dl}>
        <dt>API base URL</dt>
        <dd>
          <code>{apiBaseUrl()}</code>
        </dd>
      </dl>
      <p style={styles.hint}>
        Next: login and authenticated routes (admin-02 … admin-05).
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
  dl: { background: '#f4f4f5', padding: '1rem', borderRadius: 8 },
  hint: { marginTop: '2rem', fontSize: '0.9rem', color: '#666' },
}
