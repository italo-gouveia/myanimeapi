import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AdminUsersPage } from './AdminUsersPage'

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
  useMutation: vi.fn(),
  useQueryClient: () => ({ invalidateQueries: vi.fn() }),
}))

vi.mock('../../lib/api/admin', () => ({
  listAdminUsers: vi.fn(),
  updateUserRole: vi.fn(),
}))

import { useQuery, useMutation } from '@tanstack/react-query'

const mockUser = {
  id: 1,
  username: 'alice',
  email: 'alice@example.com',
  role: 'user',
  is_active: true,
  created_at: '2024-01-01T00:00:00Z',
}

function renderPage() {
  vi.mocked(useMutation).mockReturnValue({
    mutate: vi.fn(), isPending: false, isError: false, error: null,
  } as unknown as ReturnType<typeof useMutation>)

  return render(
    <MemoryRouter>
      <AdminUsersPage />
    </MemoryRouter>,
  )
}

describe('AdminUsersPage', () => {
  it('renders Users heading', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: true, isError: false } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('heading', { name: /users/i })).toBeInTheDocument()
  })

  it('shows skeleton rows while loading', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: true, isError: false } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getAllByRole('row').length).toBeGreaterThan(1)
  })

  it('renders user rows when loaded', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: { users: [mockUser], total: 1 },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('alice')).toBeInTheDocument()
    expect(screen.getByText('alice@example.com')).toBeInTheDocument()
    expect(screen.getByText('Active')).toBeInTheDocument()
  })

  it('shows error alert on failure', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: false, isError: true } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/failed to load users/i)).toBeInTheDocument()
  })

  it('renders total count', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: { users: [mockUser], total: 42 },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/42/)).toBeInTheDocument()
  })

  it('renders inactive user status correctly', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: { users: [{ ...mockUser, is_active: false }], total: 1 },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('Inactive')).toBeInTheDocument()
  })
})
