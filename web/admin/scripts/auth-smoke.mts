import { mockLogin } from '../src/lib/auth/dev-mock.ts'
import { DEV_MOCK_USERS } from '../src/lib/auth/dev-users.ts'
import { decodeAccessToken } from '../src/lib/auth/jwt.ts'
import { ApiError } from '../src/lib/api/errors.ts'

function assert(condition: unknown, message: string): void {
  if (!condition) throw new Error(message)
}

async function main(): Promise<void> {
  const admin = DEV_MOCK_USERS[0]!
  const ok = await mockLogin({
    registrationId: admin.registrationId,
    password: admin.password,
  })
  assert(ok.user.roles.includes('ORG_ADMIN'), 'mock user must be ORG_ADMIN')
  assert(ok.tokens.accessToken.length > 0, 'access token required')

  const decoded = decodeAccessToken(ok.tokens.accessToken)
  assert(decoded?.id === ok.user.id, 'JWT payload must match user id')

  const manager = DEV_MOCK_USERS.find((u) => u.registrationId === 'manager')!
  const managerLogin = await mockLogin({
    registrationId: manager.registrationId,
    password: manager.password,
  })
  assert(managerLogin.user.roles.includes('MANAGER'), 'manager mock must be MANAGER')

  let rejected = false
  try {
    await mockLogin({ registrationId: 'wrong', password: 'wrong' })
  } catch (err) {
    rejected = err instanceof ApiError && err.code === 'UNAUTHORIZED'
  }
  assert(rejected, 'invalid credentials must return UNAUTHORIZED')

  console.log('AUTH SMOKE OK')
}

main().catch((err) => {
  console.error('AUTH SMOKE FAILED:', err)
  process.exit(1)
})
