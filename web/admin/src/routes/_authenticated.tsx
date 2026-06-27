import { Outlet, createFileRoute, redirect, useNavigate, useRouteContext } from '@tanstack/react-router'
import { useEffect, useState } from 'react'
import { AdminShell } from '~/components/AdminShell'

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: ({ context, location }) => {
    if (!context.auth.isAuthenticated) {
      throw redirect({
        to: '/login',
        search: { redirect: location.href },
      })
    }
  },
  component: AuthenticatedLayout,
})

function AuthenticatedLayout() {
  const { auth } = useRouteContext({ from: '__root__' })
  const navigate = useNavigate()
  const [hydrated, setHydrated] = useState(false)

  useEffect(() => {
    setHydrated(true)
  }, [])

  useEffect(() => {
    if (hydrated && !auth.isAuthenticated) {
      void navigate({
        to: '/login',
        search: { redirect: window.location.href },
      })
    }
  }, [hydrated, auth.isAuthenticated, navigate])

  if (!hydrated || !auth.isAuthenticated) {
    return <p style={{ fontFamily: 'system-ui, sans-serif', padding: '2rem' }}>Loading session…</p>
  }

  return (
    <AdminShell>
      <Outlet />
    </AdminShell>
  )
}
