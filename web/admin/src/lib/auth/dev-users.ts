import type { AuthRole } from './types'

export type DevMockUser = {
  registrationId: string
  password: string
  displayName: string
  roles: AuthRole[]
  id: string
  tenantId: string
}

/** Mirrors infra/dev/seed-admin-users.sql — mock auth only until POST /v1/auth/login exists. */
export const DEV_MOCK_USERS: DevMockUser[] = [
  {
    registrationId: 'admin',
    password: 'admin',
    displayName: 'Admin User',
    roles: ['ORG_ADMIN'],
    id: '33333333-3333-3333-3333-333333333333',
    tenantId: '11111111-1111-1111-1111-111111111111',
  },
  {
    registrationId: 'manager',
    password: 'manager',
    displayName: 'Manager User',
    roles: ['MANAGER'],
    id: '34444444-4444-4444-4444-444444444444',
    tenantId: '11111111-1111-1111-1111-111111111111',
  },
  {
    registrationId: 'hr',
    password: 'hr',
    displayName: 'HR Analyst',
    roles: ['HR_ANALYST'],
    id: '35555555-5555-5555-5555-555555555555',
    tenantId: '11111111-1111-1111-1111-111111111111',
  },
  {
    registrationId: 'auditor',
    password: 'auditor',
    displayName: 'Auditor User',
    roles: ['AUDITOR'],
    id: '36666666-6666-6666-6666-666666666666',
    tenantId: '11111111-1111-1111-1111-111111111111',
  },
]
