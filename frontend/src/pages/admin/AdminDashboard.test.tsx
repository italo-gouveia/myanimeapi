import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AdminDashboard } from './AdminDashboard'

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
}))

vi.mock('../../lib/api/admin', () => ({
  getAdminStats: vi.fn(),
}))

import { useQuery } from '@tanstack/react-query'

function renderDashboard() {
  return render(
    <MemoryRouter>
      <AdminDashboard />
    </MemoryRouter>,
  )
}

describe('AdminDashboard', () => {
  it('renders dashboard heading', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: true, isError: false } as ReturnType<typeof useQuery>)
    renderDashboard()
    expect(screen.getByRole('heading', { name: /dashboard/i })).toBeInTheDocument()
  })

  it('renders stat cards while loading', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: true, isError: false } as ReturnType<typeof useQuery>)
    renderDashboard()
    expect(screen.getAllByText(/users/i).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/anime/i).length).toBeGreaterThan(0)
    expect(screen.getAllByText(/reviews/i).length).toBeGreaterThan(0)
  })

  it('renders stat values when loaded', () => {
    vi.mocked(useQuery).mockReturnValue({
      data: {
        users_total: 42,
        anime_total: 100,
        reviews_total: 300,
        favorites_total: 500,
        genres_total: 20,
        tags_total: 15,
        with_cover_total: 80,
      },
      isLoading: false,
      isError: false,
    } as ReturnType<typeof useQuery>)
    renderDashboard()
    expect(screen.getByText('42')).toBeInTheDocument()
    expect(screen.getByText('100')).toBeInTheDocument()
  })

  it('shows error alert when query fails', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: false, isError: true } as ReturnType<typeof useQuery>)
    renderDashboard()
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/failed to load stats/i)).toBeInTheDocument()
  })

  it('renders quick action links', () => {
    vi.mocked(useQuery).mockReturnValue({ data: undefined, isLoading: false, isError: false } as ReturnType<typeof useQuery>)
    renderDashboard()
    expect(screen.getByRole('link', { name: /run etl sync/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /manage users/i })).toBeInTheDocument()
  })
})
