import {
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ReactNode,
} from 'react'
import { loginRequest } from './login-api'
import {
  clearStoredSession,
  loadStoredSession,
  saveStoredSession,
} from './storage'
import type { AuthState, LoginCredentials } from './types'

const AuthContext = createContext<AuthState | null>(null)

function readInitialAuth(): Pick<AuthState, 'user' | 'tokens'> {
  const session = loadStoredSession()
  if (!session) {
    return { user: null, tokens: null }
  }
  return { user: session.user, tokens: session.tokens }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [initialAuth] = useState(readInitialAuth)
  const [user, setUser] = useState(initialAuth.user)
  const [tokens, setTokens] = useState(initialAuth.tokens)

  const login = useCallback(async (credentials: LoginCredentials) => {
    const result = await loginRequest(credentials)
    saveStoredSession(result)
    setUser(result.user)
    setTokens(result.tokens)
  }, [])

  const logout = useCallback(() => {
    clearStoredSession()
    setUser(null)
    setTokens(null)
  }, [])

  const value = useMemo<AuthState>(
    () => ({
      user,
      tokens,
      isAuthenticated: user !== null && tokens !== null,
      isLoading: false,
      login,
      logout,
    }),
    [user, tokens, login, logout],
  )

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

export function useAuth(): AuthState {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used within AuthProvider')
  }
  return ctx
}
