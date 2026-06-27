import type { ReactNode } from 'react'
import { useState } from 'react'
import { Link, useRouter, useRouteContext } from '@tanstack/react-router'
import './admin-shell.css'

type NavItem = {
  label: string
  to?: '/dashboard'
  disabled?: boolean
  icon: ReactNode
}

const NAV_SECTIONS: { title: string; items: NavItem[] }[] = [
  {
    title: 'Overview',
    items: [{ label: 'Dashboard', to: '/dashboard', icon: <IconDashboard /> }],
  },
  {
    title: 'Management',
    items: [
      { label: 'Employees', disabled: true, icon: <IconUsers /> },
      { label: 'Settings', disabled: true, icon: <IconSettings /> },
    ],
  },
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
    <div className="admin-shell">
      {sidebarOpen ? (
        <div
          className="admin-shell-backdrop"
          role="presentation"
          onClick={() => setSidebarOpen(false)}
        />
      ) : null}

      <aside className={`admin-shell-sidebar${sidebarOpen ? ' is-open' : ''}`}>
        <div className="admin-shell-brand">
          <div className="admin-shell-brand-mark" aria-hidden="true">
            OP
          </div>
          <div className="admin-shell-brand-text">
            <span className="admin-shell-brand-name">OpenPresence</span>
            <span className="admin-shell-brand-sub">Admin panel</span>
          </div>
        </div>

        <nav className="admin-shell-nav" aria-label="Main navigation">
          {NAV_SECTIONS.map((section) => (
            <div key={section.title} className="admin-shell-nav-section">
              <p className="admin-shell-nav-label">{section.title}</p>
              {section.items.map((item) =>
                item.disabled ? (
                  <span
                    key={item.label}
                    className="admin-shell-nav-item is-disabled"
                    aria-disabled="true"
                  >
                    <span className="admin-shell-nav-icon">{item.icon}</span>
                    <span className="admin-shell-nav-label-text">{item.label}</span>
                    <span className="admin-shell-nav-badge">Soon</span>
                  </span>
                ) : (
                  <Link
                    key={item.label}
                    to={item.to!}
                    className="admin-shell-nav-item"
                    activeProps={{ className: 'admin-shell-nav-item is-active' }}
                    onClick={() => setSidebarOpen(false)}
                  >
                    <span className="admin-shell-nav-icon">{item.icon}</span>
                    <span className="admin-shell-nav-label-text">{item.label}</span>
                  </Link>
                ),
              )}
            </div>
          ))}
        </nav>

        <div className="admin-shell-sidebar-footer">
          Signed in as {auth.user?.displayName ?? 'user'}
        </div>
      </aside>

      <div className="admin-shell-main">
        <header className="admin-shell-header">
          <button
            type="button"
            className="admin-shell-menu"
            aria-label="Toggle navigation"
            onClick={() => setSidebarOpen((open) => !open)}
          >
            ☰
          </button>
          <div className="admin-shell-header-meta">
            <span className="admin-shell-user-name">{auth.user?.displayName ?? 'User'}</span>
            <span className="admin-shell-tenant-id">Tenant {shortId(auth.user?.tenantId)}</span>
          </div>
          <button type="button" className="admin-shell-logout" onClick={handleLogout}>
            Sign out
          </button>
        </header>
        <div className="admin-shell-content">{children}</div>
      </div>
    </div>
  )
}

function shortId(id: string | undefined): string {
  if (!id) return '—'
  return id.slice(0, 8)
}

function IconDashboard() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <rect x="3" y="3" width="8" height="8" rx="1.5" stroke="currentColor" strokeWidth="1.75" />
      <rect x="13" y="3" width="8" height="5" rx="1.5" stroke="currentColor" strokeWidth="1.75" />
      <rect x="13" y="10" width="8" height="11" rx="1.5" stroke="currentColor" strokeWidth="1.75" />
      <rect x="3" y="13" width="8" height="8" rx="1.5" stroke="currentColor" strokeWidth="1.75" />
    </svg>
  )
}

function IconUsers() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <circle cx="9" cy="8" r="3.25" stroke="currentColor" strokeWidth="1.75" />
      <path
        d="M4 19c0-2.76 2.24-5 5-5s5 2.24 5 5"
        stroke="currentColor"
        strokeWidth="1.75"
        strokeLinecap="round"
      />
      <path
        d="M16 11.5a2.5 2.5 0 1 0 0-5"
        stroke="currentColor"
        strokeWidth="1.75"
        strokeLinecap="round"
      />
      <path
        d="M19 19c0-2.21-1.57-4.05-3.65-4.47"
        stroke="currentColor"
        strokeWidth="1.75"
        strokeLinecap="round"
      />
    </svg>
  )
}

function IconSettings() {
  return (
    <svg width="18" height="18" viewBox="0 0 24 24" fill="none" aria-hidden="true">
      <circle cx="12" cy="12" r="3" stroke="currentColor" strokeWidth="1.75" />
      <path
        d="M12 2v2M12 20v2M4.93 4.93l1.41 1.41M17.66 17.66l1.41 1.41M2 12h2M20 12h2M4.93 19.07l1.41-1.41M17.66 6.34l1.41-1.41"
        stroke="currentColor"
        strokeWidth="1.75"
        strokeLinecap="round"
      />
    </svg>
  )
}
