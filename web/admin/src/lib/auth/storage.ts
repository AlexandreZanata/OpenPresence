import type { AuthTokens, AuthUser } from './types'
import { decodeAccessToken } from './jwt'

const STORAGE_KEY = 'openpresence.admin.auth'

type StoredSession = {
  tokens: AuthTokens
  user: AuthUser
}

function canUseStorage(): boolean {
  return typeof window !== 'undefined' && typeof window.localStorage !== 'undefined'
}

export function loadStoredSession(): StoredSession | null {
  if (!canUseStorage()) return null
  const raw = window.localStorage.getItem(STORAGE_KEY)
  if (!raw) return null
  try {
    const parsed = JSON.parse(raw) as StoredSession
    if (!parsed.tokens?.accessToken) return null
    const user = parsed.user ?? decodeAccessToken(parsed.tokens.accessToken)
    if (!user) return null
    return { tokens: parsed.tokens, user }
  } catch {
    return null
  }
}

export function saveStoredSession(session: StoredSession): void {
  if (!canUseStorage()) return
  window.localStorage.setItem(STORAGE_KEY, JSON.stringify(session))
}

export function clearStoredSession(): void {
  if (!canUseStorage()) return
  window.localStorage.removeItem(STORAGE_KEY)
}

export function getStoredAccessToken(): string | null {
  return loadStoredSession()?.tokens.accessToken ?? null
}
