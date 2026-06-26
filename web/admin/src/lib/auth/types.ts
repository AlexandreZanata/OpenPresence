export type AuthRole = 'ORG_ADMIN' | 'MANAGER' | 'HR_ANALYST' | 'AUDITOR'

export type AuthUser = {
  id: string
  tenantId: string
  registrationId: string
  roles: AuthRole[]
  displayName: string
}

export type AuthTokens = {
  accessToken: string
  refreshToken: string
  expiresIn: number
}

export type LoginCredentials = {
  registrationId: string
  password: string
}

export type LoginResult = {
  user: AuthUser
  tokens: AuthTokens
}

export type AuthState = {
  user: AuthUser | null
  tokens: AuthTokens | null
  isAuthenticated: boolean
  isLoading: boolean
  login: (credentials: LoginCredentials) => Promise<void>
  logout: () => void
}
