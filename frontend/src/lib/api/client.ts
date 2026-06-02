import axios, { type AxiosError } from 'axios'
import type { ApiError } from './types'

const STORAGE_KEY = 'myanimeapi.session'

interface StoredSession {
  token?: string
}

function readToken(): string | null {
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return null
    const parsed = JSON.parse(raw) as StoredSession
    return parsed.token ?? null
  } catch {
    return null
  }
}

export const apiClient = axios.create({
  baseURL: import.meta.env.VITE_API_URL ?? '',
  timeout: 15_000,
  headers: { 'Content-Type': 'application/json' },
})

apiClient.interceptors.request.use((config) => {
  const token = readToken()
  if (token) {
    config.headers.set('Authorization', `Bearer ${token}`)
  }
  return config
})

apiClient.interceptors.response.use(
  (response) => response,
  (error: AxiosError<ApiError>) => {
    // 401 on a protected route — wipe local session so the UI redirects to /login.
    // We don't redirect from here to avoid a hard navigation; the page query
    // will fail and the ProtectedRoute / AuthProvider will react on next render.
    if (error.response?.status === 401) {
      window.localStorage.removeItem(STORAGE_KEY)
    }
    return Promise.reject(error)
  },
)

export function extractErrorMessage(error: unknown, fallback: string): string {
  if (axios.isAxiosError<ApiError>(error)) {
    return error.response?.data?.message ?? error.response?.data?.details ?? fallback
  }
  return fallback
}
