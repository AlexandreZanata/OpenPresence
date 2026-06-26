import { ApiError } from '~/lib/api/errors'
import { createMockAccessToken } from './jwt'
import type { AuthRole, LoginCredentials, LoginResult } from './types'

const MOCK_REGISTRATION = 'admin'
const MOCK_PASSWORD = 'admin'

export async function mockLogin(credentials: LoginCredentials): Promise<LoginResult> {
  if (
    credentials.registrationId !== MOCK_REGISTRATION ||
    credentials.password !== MOCK_PASSWORD
  ) {
    throw new ApiError('UNAUTHORIZED', 'Invalid registration ID or password', 401)
  }

  const user = {
    id: '33333333-3333-3333-3333-333333333333',
    tenantId: '11111111-1111-1111-1111-111111111111',
    registrationId: MOCK_REGISTRATION,
    roles: ['ORG_ADMIN'] as AuthRole[],
    displayName: 'Admin User',
  }

  return {
    user,
    tokens: {
      accessToken: createMockAccessToken(user),
      refreshToken: 'mock-refresh-token',
      expiresIn: 900,
    },
  }
}
