import type { AuthUser } from './types'

type JwtPayload = {
  sub?: string
  tenant_id?: string
  registration_id?: string
  roles?: string[]
  name?: string
}

export function decodeAccessToken(token: string): AuthUser | null {
  try {
    const [, payload] = token.split('.')
    if (!payload) return null
    const json = JSON.parse(atob(payload.replace(/-/g, '+').replace(/_/g, '/'))) as JwtPayload
    if (!json.sub || !json.tenant_id) return null
    return {
      id: json.sub,
      tenantId: json.tenant_id,
      registrationId: json.registration_id ?? '',
      roles: (json.roles ?? []) as AuthUser['roles'],
      displayName: json.name ?? json.registration_id ?? 'User',
    }
  } catch {
    return null
  }
}

export function createMockAccessToken(user: AuthUser): string {
  const header = btoa(JSON.stringify({ alg: 'none', typ: 'JWT' }))
  const body = btoa(
    JSON.stringify({
      sub: user.id,
      tenant_id: user.tenantId,
      registration_id: user.registrationId,
      roles: user.roles,
      name: user.displayName,
    }),
  )
  return `${header}.${body}.mock`
}
