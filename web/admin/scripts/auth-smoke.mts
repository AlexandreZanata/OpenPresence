import { mockLogin } from '../src/lib/auth/dev-mock.ts'
import { decodeAccessToken } from '../src/lib/auth/jwt.ts'
import { ApiError } from '../src/lib/api/errors.ts'

function assert(condition: unknown, message: string): void {
  if (!condition) throw new Error(message)
}

async function main(): Promise<void> {
  const ok = await mockLogin({ registrationId: 'admin', password: 'admin' })
  assert(ok.user.roles.includes('ORG_ADMIN'), 'mock user must be ORG_ADMIN')
  assert(ok.tokens.accessToken.length > 0, 'access token required')

  const decoded = decodeAccessToken(ok.tokens.accessToken)
  assert(decoded?.id === ok.user.id, 'JWT payload must match user id')

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
