import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter, Routes, Route } from 'react-router-dom'
import { Layout } from './Layout'

vi.mock('../lib/auth/AuthProvider', () => ({ useAuth: vi.fn() }))
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => vi.fn() }
})

import { useAuth } from '../lib/auth/AuthProvider'

function mockAuth(overrides = {}) {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated: false,
    isAdmin: false,
    token: null,
    username: null,
    role: 'user',
    userId: null,
    setSession: vi.fn(),
    clearSession: vi.fn(),
    ...overrides,
  })
}

function renderLayout() {
  return render(
    <MemoryRouter>
      <Routes>
        <Route element={<Layout />}>
          <Route index element={<div>home content</div>} />
        </Route>
      </Routes>
    </MemoryRouter>,
  )
}

describe('Layout', () => {
  it('renders brand link', () => {
    mockAuth()
    renderLayout()
    expect(screen.getByRole('link', { name: /myanimeapi/i })).toBeInTheDocument()
  })

  it('renders Catalog and Characters nav links', () => {
    mockAuth()
    renderLayout()
    expect(screen.getByRole('link', { name: /catalog/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /characters/i })).toBeInTheDocument()
  })

  it('shows Sign in and Register links when not authenticated', () => {
    mockAuth()
    renderLayout()
    expect(screen.getByRole('link', { name: /sign in/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /register/i })).toBeInTheDocument()
  })

  it('shows username and Sign out button when authenticated', () => {
    mockAuth({ isAuthenticated: true, username: 'italo' })
    renderLayout()
    expect(screen.getByTestId('header-username')).toHaveTextContent('italo')
    expect(screen.getByTestId('logout-button')).toBeInTheDocument()
  })

  it('shows authenticated nav links when logged in', () => {
    mockAuth({ isAuthenticated: true, username: 'italo' })
    renderLayout()
    expect(screen.getByRole('link', { name: /favorites/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /watchlist/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /profile/i })).toBeInTheDocument()
  })

  it('shows Admin link for admin users', () => {
    mockAuth({ isAuthenticated: true, isAdmin: true, username: 'admin' })
    renderLayout()
    expect(screen.getByRole('link', { name: /admin/i })).toBeInTheDocument()
  })

  it('does not show Admin link for regular users', () => {
    mockAuth({ isAuthenticated: true, isAdmin: false, username: 'user' })
    renderLayout()
    expect(screen.queryByRole('link', { name: /^admin$/i })).not.toBeInTheDocument()
  })

  it('calls clearSession on Sign out click', async () => {
    const clearSession = vi.fn()
    mockAuth({ isAuthenticated: true, username: 'italo', clearSession })
    renderLayout()
    await userEvent.click(screen.getByTestId('logout-button'))
    expect(clearSession).toHaveBeenCalledOnce()
  })

  it('renders footer text', () => {
    mockAuth()
    renderLayout()
    expect(screen.getByRole('contentinfo')).toHaveTextContent('MyAnimeAPI')
  })

  it('renders outlet content', () => {
    mockAuth()
    renderLayout()
    expect(screen.getByText('home content')).toBeInTheDocument()
  })
})
