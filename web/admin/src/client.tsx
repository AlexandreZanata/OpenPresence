import { StrictMode, startTransition } from 'react'
import { hydrateRoot } from 'react-dom/client'
import { Await } from '@tanstack/react-router'
import { hydrateStart } from '@tanstack/react-start/client'
import { AuthProvider } from '~/lib/auth/AuthProvider'
import { AuthRouterProvider } from '~/components/AuthRouterProvider'

const routerPromise = hydrateStart()

function ClientRoot() {
  return (
    <AuthProvider>
      <Await promise={routerPromise}>
        {(router) => <AuthRouterProvider router={router} />}
      </Await>
    </AuthProvider>
  )
}

startTransition(() => {
  hydrateRoot(
    document,
    <StrictMode>
      <ClientRoot />
    </StrictMode>,
  )
})
