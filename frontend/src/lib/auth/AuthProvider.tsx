import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'

interface AuthState {
  token: string | null
  username: string | null
}

interface AuthContextValue extends AuthState {
  setSession: (token: string, username: string) => void
  clearSession: () => void
  isAuthenticated: boolean
  isAdmin: boolean
  role: string
}

/** Decode JWT payload without a library (read-only — no signature verification). */
function decodeJWT(token: string): { is_admin?: boolean; role?: string } {
  try {
    const payload = token.split('.')[1]
    const base64 = payload.replace(/-/g, '+').replace(/_/g, '/')
    const padded = base64.padEnd(base64.length + ((4 - (base64.length % 4)) % 4), '=')
    return JSON.parse(atob(padded)) as { is_admin?: boolean; role?: string }
  } catch {
    return {}
  }
}

const STORAGE_KEY = 'myanimeapi.session'

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

function readSession(): AuthState {
  try {
    const raw = window.localStorage.getItem(STORAGE_KEY)
    if (!raw) return { token: null, username: null }
    const parsed = JSON.parse(raw) as Partial<AuthState>
    return {
      token: parsed.token ?? null,
      username: parsed.username ?? null,
    }
  } catch {
    return { token: null, username: null }
  }
}

export function AuthProvider({ children }: { children: ReactNode }) {
  const [state, setState] = useState<AuthState>(() => readSession())

  useEffect(() => {
    if (state.token) {
      window.localStorage.setItem(STORAGE_KEY, JSON.stringify(state))
    } else {
      window.localStorage.removeItem(STORAGE_KEY)
    }
  }, [state])

  const setSession = (token: string, username: string) => {
    setState({ token, username })
  }

  const clearSession = () => {
    setState({ token: null, username: null })
  }

  const claims = state.token ? decodeJWT(state.token) : {}

  const value: AuthContextValue = {
    ...state,
    setSession,
    clearSession,
    isAuthenticated: Boolean(state.token),
    isAdmin: Boolean(claims.is_admin),
    role: claims.role ?? 'user',
  }

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
}

// eslint-disable-next-line react-refresh/only-export-components
export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext)
  if (!ctx) {
    throw new Error('useAuth must be used inside <AuthProvider>')
  }
  return ctx
}
