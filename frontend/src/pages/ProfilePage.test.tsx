import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { MemoryRouter } from 'react-router-dom'
import { ProfilePage } from './ProfilePage'

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
}))

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => vi.fn() }
})

vi.mock('../lib/api/users', () => ({ getProfile: vi.fn() }))
vi.mock('../lib/api/export', () => ({ exportUserData: vi.fn() }))
vi.mock('../lib/auth/AuthProvider', () => ({
  useAuth: vi.fn(),
}))

import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/AuthProvider'

const mockUser = {
  id: 1,
  username: 'italo',
  email: 'italo@example.com',
  is_active: true,
  is_admin: false,
  bio: 'Anime fan',
  profile_pic: null,
  created_at: '2024-03-15T12:00:00Z',
}

function renderPage() {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated: true,
    isAdmin: false,
    token: 'tok',
    username: 'italo',
    role: 'user',
    userId: 1,
    setSession: vi.fn(),
    clearSession: vi.fn(),
  })
  return render(<MemoryRouter><ProfilePage /></MemoryRouter>)
}

describe('ProfilePage', () => {
  it('shows loading text', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: true, data: undefined, isError: false } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/loading profile/i)).toBeInTheDocument()
  })

  it('shows error state', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: true,
      data: undefined,
      error: new Error('Network error'),
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/could not load profile/i)).toBeInTheDocument()
  })

  it('renders username and email', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: mockUser,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByTestId('profile-username')).toHaveTextContent('italo')
    expect(screen.getByTestId('profile-email')).toHaveTextContent('italo@example.com')
  })

  it('renders bio when present', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: mockUser } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByTestId('profile-bio')).toHaveTextContent('Anime fan')
  })

  it('renders active status', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: mockUser } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('Active')).toBeInTheDocument()
  })

  it('renders inactive status', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { ...mockUser, is_active: false },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('Inactive')).toBeInTheDocument()
  })

  it('renders admin badge for admin users', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { ...mockUser, is_admin: true },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('Admin')).toBeInTheDocument()
  })

  it('renders export buttons', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: mockUser } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByTestId('export-json')).toBeInTheDocument()
    expect(screen.getByTestId('export-csv')).toBeInTheDocument()
  })

  it('calls clearSession when sign out clicked', async () => {
    const clearSession = vi.fn()
    vi.mocked(useAuth).mockReturnValue({
      isAuthenticated: true, isAdmin: false, token: 'tok', username: 'italo',
      role: 'user', userId: 1, setSession: vi.fn(), clearSession,
    })
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: mockUser } as ReturnType<typeof useQuery>)
    render(<MemoryRouter><ProfilePage /></MemoryRouter>)
    await userEvent.click(screen.getByTestId('profile-logout'))
    expect(clearSession).toHaveBeenCalledOnce()
  })
})
