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

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<AuthState['user']>(null)
  const [tokens, setTokens] = useState<AuthState['tokens']>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const session = loadStoredSession()
    if (session) {
      setUser(session.user)
      setTokens(session.tokens)
    }
    setIsLoading(false)
  }, [])

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
      isLoading,
      login,
      logout,
    }),
    [user, tokens, isLoading, login, logout],
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
