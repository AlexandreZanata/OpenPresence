import { ApiError } from '~/lib/api/errors'
import { createMockAccessToken } from './jwt'
import { DEV_MOCK_USERS } from './dev-users'
import type { LoginCredentials, LoginResult } from './types'

export async function mockLogin(credentials: LoginCredentials): Promise<LoginResult> {
  const match = DEV_MOCK_USERS.find(
    (user) =>
      user.registrationId === credentials.registrationId &&
      user.password === credentials.password,
  )

  if (!match) {
    throw new ApiError('UNAUTHORIZED', 'Invalid registration ID or password', 401)
  }

  const user = {
    id: match.id,
    tenantId: match.tenantId,
    registrationId: match.registrationId,
    roles: match.roles,
    displayName: match.displayName,
  }

  return {
    user,
    tokens: {
      accessToken: createMockAccessToken(user),
      refreshToken: `mock-refresh-${match.registrationId}`,
      expiresIn: 900,
    },
  }
}
