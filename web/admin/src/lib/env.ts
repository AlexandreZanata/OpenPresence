/** Vite-injected env vars for the admin panel. */
export function apiBaseUrl(): string {
  return import.meta.env.VITE_API_BASE_URL ?? 'http://127.0.0.1:8088'
}

export function isAuthMockEnabled(): boolean {
  return import.meta.env.VITE_AUTH_MOCK === 'true'
}
