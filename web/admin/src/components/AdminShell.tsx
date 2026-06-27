import type { CSSProperties, ReactNode } from 'react'
import { useState } from 'react'
import { Link, useRouter, useRouteContext } from '@tanstack/react-router'
import './admin-shell.css'

type NavItem = {
  label: string
  to?: '/dashboard'
  disabled?: boolean
}

const NAV_ITEMS: NavItem[] = [
  { label: 'Dashboard', to: '/dashboard' },
  { label: 'Employees', disabled: true },
  { label: 'Settings', disabled: true },
]

type Props = {
  children: ReactNode
}

export function AdminShell({ children }: Props) {
  const { auth } = useRouteContext({ from: '__root__' })
  const router = useRouter()
  const [sidebarOpen, setSidebarOpen] = useState(false)

  function handleLogout() {
    auth.logout()
    void router.navigate({ to: '/login' })
  }

  return (
    <div style={styles.app}>
      {sidebarOpen ? (
        <div
          className="admin-shell-backdrop"
          role="presentation"
          onClick={() => setSidebarOpen(false)}
        />
      ) : null}
      <aside
        className={`admin-shell-sidebar${sidebarOpen ? ' is-open' : ''}`}
        style={styles.sidebar}
      >
        <div style={styles.brand}>OpenPresence</div>
        <nav style={styles.nav}>
          {NAV_ITEMS.map((item) =>
            item.disabled ? (
              <span key={item.label} style={styles.navDisabled} aria-disabled="true">
                {item.label}
              </span>
            ) : (
              <Link
                key={item.label}
                to={item.to!}
                style={styles.navLink}
                activeProps={{ style: styles.navLinkActive }}
                onClick={() => setSidebarOpen(false)}
              >
                {item.label}
              </Link>
            ),
          )}
        </nav>
      </aside>

      <div style={styles.main}>
        <header style={styles.header}>
          <button
            type="button"
            className="admin-shell-menu"
            style={styles.menuButton}
            aria-label="Toggle navigation"
            onClick={() => setSidebarOpen((open) => !open)}
          >
            ☰
          </button>
          <div style={styles.headerMeta}>
            <span style={styles.userName}>{auth.user?.displayName ?? 'User'}</span>
            <span style={styles.tenantId}>Tenant {shortId(auth.user?.tenantId)}</span>
          </div>
          <button type="button" style={styles.logoutButton} onClick={handleLogout}>
            Sign out
          </button>
        </header>
        <div style={styles.content}>{children}</div>
      </div>
    </div>
  )
}

function shortId(id: string | undefined): string {
  if (!id) return '—'
  return id.slice(0, 8)
}

const styles: Record<string, CSSProperties> = {
  app: {
    display: 'flex',
    minHeight: '100vh',
    fontFamily: 'system-ui, sans-serif',
    background: '#f8fafc',
    color: '#0f172a',
  },
  sidebar: {
    width: 240,
    background: '#0f172a',
    color: '#e2e8f0',
    padding: '1.25rem 1rem',
    flexShrink: 0,
    transition: 'transform 0.2s ease',
  },
  brand: {
    fontWeight: 700,
    fontSize: '1.1rem',
    marginBottom: '1.5rem',
    letterSpacing: '0.02em',
  },
  nav: { display: 'flex', flexDirection: 'column', gap: '0.35rem' },
  navLink: {
    display: 'block',
    padding: '0.55rem 0.75rem',
    borderRadius: 8,
    color: '#cbd5e1',
    textDecoration: 'none',
  },
  navLinkActive: {
    background: '#1e293b',
    color: '#fff',
    fontWeight: 600,
  },
  navDisabled: {
    display: 'block',
    padding: '0.55rem 0.75rem',
    borderRadius: 8,
    color: '#64748b',
    cursor: 'not-allowed',
  },
  main: { flex: 1, display: 'flex', flexDirection: 'column', minWidth: 0 },
  header: {
    display: 'flex',
    alignItems: 'center',
    gap: '1rem',
    padding: '0.85rem 1.25rem',
    background: '#fff',
    borderBottom: '1px solid #e2e8f0',
  },
  menuButton: {
    border: '1px solid #cbd5e1',
    background: '#fff',
    borderRadius: 6,
    padding: '0.35rem 0.6rem',
    cursor: 'pointer',
  },
  headerMeta: { flex: 1, display: 'flex', flexDirection: 'column', gap: '0.15rem' },
  userName: { fontWeight: 600 },
  tenantId: { fontSize: '0.85rem', color: '#64748b' },
  logoutButton: {
    border: '1px solid #cbd5e1',
    background: '#fff',
    borderRadius: 6,
    padding: '0.45rem 0.85rem',
    cursor: 'pointer',
  },
  content: { flex: 1, padding: '1.5rem' },
}
