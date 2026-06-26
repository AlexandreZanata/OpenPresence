import { apiBaseUrl } from '~/lib/env'
import { getStoredAccessToken } from '~/lib/auth/storage'
import { ApiError, parseApiError } from './errors'

export type ApiRequestOptions = {
  method?: string
  body?: unknown
  token?: string | null
  headers?: Record<string, string>
}

export async function apiFetch<T>(path: string, options: ApiRequestOptions = {}): Promise<T> {
  const token = options.token === undefined ? getStoredAccessToken() : options.token
  const headers: Record<string, string> = {
    Accept: 'application/json',
    ...options.headers,
  }
  if (options.body !== undefined) {
    headers['Content-Type'] = 'application/json'
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }

  const response = await fetch(`${apiBaseUrl()}${path}`, {
    method: options.method ?? (options.body !== undefined ? 'POST' : 'GET'),
    headers,
    body: options.body !== undefined ? JSON.stringify(options.body) : undefined,
  })

  if (!response.ok) {
    throw await parseApiError(response)
  }

  if (response.status === 204) {
    return undefined as T
  }

  return (await response.json()) as T
}

export { ApiError }
