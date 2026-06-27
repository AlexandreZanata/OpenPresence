import { apiFetch } from './client'

export type LiveHealth = {
  status: string
}

export async function fetchLiveHealth(): Promise<LiveHealth> {
  return apiFetch<LiveHealth>('/health/live', { token: null })
}
