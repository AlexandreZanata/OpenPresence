import { Outlet, createFileRoute, redirect, useNavigate, useRouteContext } from '@tanstack/react-router'
import { useEffect } from 'react'
import { AdminShell } from '~/components/AdminShell'

export const Route = createFileRoute('/_authenticated')({
  beforeLoad: ({ context, location }) => {
    if (context.auth.isLoading) {
      return
    }
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

  useEffect(() => {
    if (!auth.isLoading && !auth.isAuthenticated) {
      void navigate({
        to: '/login',
        search: { redirect: window.location.href },
      })
    }
  }, [auth.isLoading, auth.isAuthenticated, navigate])

  if (auth.isLoading) {
    return <p style={{ fontFamily: 'system-ui, sans-serif', padding: '2rem' }}>Loading session…</p>
  }

  if (!auth.isAuthenticated) {
    return null
  }

  return (
    <AdminShell>
      <Outlet />
    </AdminShell>
  )
}
