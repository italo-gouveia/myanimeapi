import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { FavoritesPage } from './FavoritesPage'

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
  useMutation: vi.fn(),
  useQueryClient: () => ({ invalidateQueries: vi.fn() }),
}))

vi.mock('../lib/api/favorites', () => ({
  listFavorites: vi.fn(),
  removeFavorite: vi.fn(),
}))

import { useQuery, useMutation } from '@tanstack/react-query'

const mockMutate = vi.fn()

function setupMutation(overrides = {}) {
  vi.mocked(useMutation).mockReturnValue({
    mutate: mockMutate,
    isPending: false,
    isError: false,
    error: null,
    ...overrides,
  } as unknown as ReturnType<typeof useMutation>)
}

function renderPage() {
  return render(
    <MemoryRouter>
      <FavoritesPage />
    </MemoryRouter>,
  )
}

describe('FavoritesPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setupMutation()
  })

  it('shows loading state', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: true, data: undefined, isError: false } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/loading favorites/i)).toBeInTheDocument()
  })

  it('shows error state', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: true,
      error: new Error('Network error'),
      data: undefined,
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/could not load favorites/i)).toBeInTheDocument()
  })

  it('shows empty state when no favorites', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: [] } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByTestId('favorites-empty')).toBeInTheDocument()
    expect(screen.getByText(/no favorites yet/i)).toBeInTheDocument()
  })

  it('renders favorites list', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: [
        {
          id: 1,
          anime: {
            id: 10,
            title: 'Naruto',
            rating: 8.5,
            status: 'Completed',
            episodes: 220,
            description: 'A ninja story',
          },
        },
      ],
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByTestId('favorites-list')).toBeInTheDocument()
    expect(screen.getByTestId('favorite-card-10')).toBeInTheDocument()
    expect(screen.getByText('Naruto')).toBeInTheDocument()
  })

  it('calls remove mutation when remove button is clicked', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: [
        {
          id: 1,
          anime: { id: 10, title: 'Naruto', rating: 8.5, status: 'Completed', episodes: 220 },
        },
      ],
    } as ReturnType<typeof useQuery>)
    renderPage()

    const removeBtn = screen.getByTestId('remove-favorite-10')
    fireEvent.click(removeBtn)
    expect(mockMutate).toHaveBeenCalledWith(10)
  })

  it('shows mutation error alert', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: [] } as ReturnType<typeof useQuery>)
    setupMutation({ isError: true, error: new Error('Failed') } as unknown as ReturnType<typeof useMutation>)
    renderPage()
    expect(screen.getByText(/failed to remove/i)).toBeInTheDocument()
  })

  it('renders page heading', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: [] } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('heading', { name: /favorites/i })).toBeInTheDocument()
  })
})
