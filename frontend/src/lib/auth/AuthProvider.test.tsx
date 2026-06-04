import { render, screen, act } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, beforeEach } from 'vitest'
import { AuthProvider, useAuth } from './AuthProvider'

// Minimal JWT with known payload for testing.
// Payload: { is_admin: true, role: "admin", user_id: "42" }
const ADMIN_TOKEN =
  'eyJhbGciOiJIUzI1NiJ9.' +
  btoa(JSON.stringify({ is_admin: true, role: 'admin', user_id: '42' }))
    .replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '') +
  '.sig'

const USER_TOKEN =
  'eyJhbGciOiJIUzI1NiJ9.' +
  btoa(JSON.stringify({ is_admin: false, role: 'user', user_id: '7' }))
    .replace(/\+/g, '-').replace(/\//g, '_').replace(/=+$/, '') +
  '.sig'

function TestConsumer() {
  const auth = useAuth()
  return (
    <div>
      <span data-testid="authenticated">{String(auth.isAuthenticated)}</span>
      <span data-testid="admin">{String(auth.isAdmin)}</span>
      <span data-testid="role">{auth.role}</span>
      <span data-testid="userId">{String(auth.userId)}</span>
      <button onClick={() => auth.setSession(ADMIN_TOKEN, 'italo')}>login-admin</button>
      <button onClick={() => auth.setSession(USER_TOKEN, 'bob')}>login-user</button>
      <button onClick={() => auth.clearSession()}>logout</button>
    </div>
  )
}

describe('AuthProvider', () => {
  beforeEach(() => {
    window.localStorage.clear()
  })

  it('starts unauthenticated when localStorage is empty', () => {
    render(<AuthProvider><TestConsumer /></AuthProvider>)
    expect(screen.getByTestId('authenticated')).toHaveTextContent('false')
    expect(screen.getByTestId('admin')).toHaveTextContent('false')
  })

  it('sets isAuthenticated and decodes admin claims after setSession', async () => {
    render(<AuthProvider><TestConsumer /></AuthProvider>)
    await userEvent.click(screen.getByText('login-admin'))
    expect(screen.getByTestId('authenticated')).toHaveTextContent('true')
    expect(screen.getByTestId('admin')).toHaveTextContent('true')
    expect(screen.getByTestId('role')).toHaveTextContent('admin')
    expect(screen.getByTestId('userId')).toHaveTextContent('42')
  })

  it('clears session after clearSession', async () => {
    render(<AuthProvider><TestConsumer /></AuthProvider>)
    await userEvent.click(screen.getByText('login-admin'))
    await userEvent.click(screen.getByText('logout'))
    expect(screen.getByTestId('authenticated')).toHaveTextContent('false')
    expect(screen.getByTestId('userId')).toHaveTextContent('null')
  })

  it('decodes non-admin user claims correctly', async () => {
    render(<AuthProvider><TestConsumer /></AuthProvider>)
    await userEvent.click(screen.getByText('login-user'))
    expect(screen.getByTestId('admin')).toHaveTextContent('false')
    expect(screen.getByTestId('role')).toHaveTextContent('user')
    expect(screen.getByTestId('userId')).toHaveTextContent('7')
  })

  it('persists session to localStorage after login', async () => {
    render(<AuthProvider><TestConsumer /></AuthProvider>)
    await userEvent.click(screen.getByText('login-admin'))
    const stored = JSON.parse(window.localStorage.getItem('myanimeapi.session')!)
    expect(stored.token).toBe(ADMIN_TOKEN)
    expect(stored.username).toBe('italo')
  })

  it('removes session from localStorage after logout', async () => {
    render(<AuthProvider><TestConsumer /></AuthProvider>)
    await userEvent.click(screen.getByText('login-admin'))
    await userEvent.click(screen.getByText('logout'))
    expect(window.localStorage.getItem('myanimeapi.session')).toBeNull()
  })

  it('hydrates session from localStorage on mount', () => {
    window.localStorage.setItem(
      'myanimeapi.session',
      JSON.stringify({ token: USER_TOKEN, username: 'bob' }),
    )
    render(<AuthProvider><TestConsumer /></AuthProvider>)
    expect(screen.getByTestId('authenticated')).toHaveTextContent('true')
    expect(screen.getByTestId('role')).toHaveTextContent('user')
  })

  it('throws when useAuth is used outside AuthProvider', () => {
    // Suppress React's error boundary output in test console
    const spy = vi.spyOn(console, 'error').mockImplementation(() => {})
    expect(() => render(<TestConsumer />)).toThrow('useAuth must be used inside <AuthProvider>')
    spy.mockRestore()
  })
})
