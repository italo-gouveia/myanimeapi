import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { ProtectedRoute } from './ProtectedRoute'

vi.mock('../lib/auth/AuthProvider', () => ({
  useAuth: vi.fn(),
}))

import { useAuth } from '../lib/auth/AuthProvider'

function renderRoute(opts: { isAuthenticated: boolean; isAdmin?: boolean; requireAdmin?: boolean }) {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated: opts.isAuthenticated,
    isAdmin: opts.isAdmin ?? false,
    token: null,
    username: null,
    role: 'user',
    userId: null,
    setSession: vi.fn(),
    clearSession: vi.fn(),
  })

  return render(
    <MemoryRouter initialEntries={['/protected']}>
      <Routes>
        <Route element={<ProtectedRoute requireAdmin={opts.requireAdmin} />}>
          <Route path="/protected" element={<div>Protected Content</div>} />
        </Route>
        <Route path="/login" element={<div>Login Page</div>} />
        <Route path="/" element={<div>Home Page</div>} />
      </Routes>
    </MemoryRouter>,
  )
}

describe('ProtectedRoute', () => {
  it('renders outlet when authenticated', () => {
    renderRoute({ isAuthenticated: true })
    expect(screen.getByText('Protected Content')).toBeInTheDocument()
  })

  it('redirects to /login when not authenticated', () => {
    renderRoute({ isAuthenticated: false })
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument()
    expect(screen.getByText('Login Page')).toBeInTheDocument()
  })

  it('renders outlet when authenticated and not admin-required', () => {
    renderRoute({ isAuthenticated: true, isAdmin: false, requireAdmin: false })
    expect(screen.getByText('Protected Content')).toBeInTheDocument()
  })

  it('redirects to / when authenticated but not admin and admin required', () => {
    renderRoute({ isAuthenticated: true, isAdmin: false, requireAdmin: true })
    expect(screen.queryByText('Protected Content')).not.toBeInTheDocument()
    expect(screen.getByText('Home Page')).toBeInTheDocument()
  })

  it('renders outlet when authenticated and admin and admin required', () => {
    renderRoute({ isAuthenticated: true, isAdmin: true, requireAdmin: true })
    expect(screen.getByText('Protected Content')).toBeInTheDocument()
  })
})
