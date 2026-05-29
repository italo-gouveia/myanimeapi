import { createContext, useContext, useState, useEffect, type ReactNode } from 'react'

interface AuthState {
  token: string | null
  username: string | null
}

interface AuthContextValue extends AuthState {
  setSession: (token: string, username: string) => void
  clearSession: () => void
  isAuthenticated: boolean
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

  const value: AuthContextValue = {
    ...state,
    setSession,
    clearSession,
    isAuthenticated: Boolean(state.token),
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
