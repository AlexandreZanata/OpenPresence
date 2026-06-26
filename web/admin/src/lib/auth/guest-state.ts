import type { AuthState } from './types'

/** SSR / pre-hydration placeholder — replaced by AuthProvider on the client. */
export function createGuestAuthState(): AuthState {
  return {
    user: null,
    tokens: null,
    isAuthenticated: false,
    isLoading: false,
    login: async () => {
      throw new Error('AuthProvider is not mounted')
    },
    logout: () => {},
  }
}
