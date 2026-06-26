import { StrictMode, startTransition, Suspense, use } from 'react'
import { hydrateRoot } from 'react-dom/client'
import { hydrateStart } from '@tanstack/react-start/client'
import { AuthProvider } from '~/lib/auth/AuthProvider'
import { AuthRouterProvider } from '~/components/AuthRouterProvider'

const routerPromise = hydrateStart()

function ClientRoot() {
  const router = use(routerPromise)
  return (
    <AuthProvider>
      <AuthRouterProvider router={router} />
    </AuthProvider>
  )
}

startTransition(() => {
  hydrateRoot(
    document,
    <StrictMode>
      <Suspense fallback={null}>
        <ClientRoot />
      </Suspense>
    </StrictMode>,
  )
})
