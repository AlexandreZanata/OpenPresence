import { apiFetch } from '~/lib/api/client'
import { isAuthMockEnabled } from '~/lib/env'
import { mockLogin } from './dev-mock'
import type { LoginCredentials, LoginResult } from './types'

type LoginResponse = {
  accessToken: string
  refreshToken: string
  expiresIn: number
  user?: LoginResult['user']
}

export async function loginRequest(credentials: LoginCredentials): Promise<LoginResult> {
  if (isAuthMockEnabled()) {
    return mockLogin(credentials)
  }

  const response = await apiFetch<LoginResponse>('/v1/auth/login', {
    method: 'POST',
    body: credentials,
    token: null,
  })

  if (!response.user) {
    throw new Error('Auth API must return user profile with tokens')
  }

  return {
    user: response.user,
    tokens: {
      accessToken: response.accessToken,
      refreshToken: response.refreshToken,
      expiresIn: response.expiresIn,
    },
  }
}
