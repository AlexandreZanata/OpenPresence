import { useEffect } from 'react'
import { RouterProvider } from '@tanstack/react-router'
import type { AnyRouter } from '@tanstack/react-router'
import { useAuth } from '~/lib/auth/AuthProvider'

type Props = {
  router: AnyRouter
}

export function AuthRouterProvider({ router }: Props) {
  const auth = useAuth()
  const queryClient = router.options.context.queryClient

  useEffect(() => {
    void router.invalidate()
  }, [auth.isAuthenticated, auth.user?.id, router])

  return <RouterProvider router={router} context={{ queryClient, auth }} />
}
