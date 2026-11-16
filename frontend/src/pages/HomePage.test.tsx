import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { HomePage } from './HomePage'

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
  keepPreviousData: undefined,
}))

vi.mock('../lib/api/animes', () => ({
  listAnimes: vi.fn(),
  searchAnimes: vi.fn(),
}))

vi.mock('../lib/api/genres', () => ({ listGenres: vi.fn() }))
vi.mock('../lib/api/favorites', () => ({ listFavorites: vi.fn() }))
vi.mock('../lib/api/watchlist', () => ({ getWatchlist: vi.fn() }))
vi.mock('../lib/auth/AuthProvider', () => ({ useAuth: vi.fn() }))

vi.mock('../components/AnimeCard', () => ({
  AnimeCard: ({ anime }: { anime: { title: string } }) => <div data-testid="anime-card">{anime.title}</div>,
}))

import { useQuery } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/AuthProvider'

function mockAuth(isAuthenticated = false) {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated, isAdmin: false, token: null, username: null,
    role: 'user', userId: null, setSession: vi.fn(), clearSession: vi.fn(),
  })
}

const mockAnime = { id: 1, title: 'Naruto', rating: 8.5, episodes: 220, status: 'Completed' }
const mockGenre = { id: 1, name: 'Action' }

function setupQueries({
  animesData = undefined as unknown,
  animesLoading = false,
  genresData = undefined as unknown,
} = {}) {
  vi.mocked(useQuery).mockImplementation(({ queryKey }: { queryKey: unknown[] }) => {
    const key = Array.isArray(queryKey) ? queryKey[0] : ''
    if (key === 'genres') {
      return { data: genresData, isLoading: false, isError: false } as ReturnType<typeof useQuery>
    }
    if (key === 'animes' || key === 'search') {
      return { data: animesData, isLoading: animesLoading, isError: false, isFetching: false } as ReturnType<typeof useQuery>
    }
    return { data: undefined, isLoading: false, isError: false } as ReturnType<typeof useQuery>
  })
}

function renderPage() {
  return render(<MemoryRouter><HomePage /></MemoryRouter>)
}

describe('HomePage', () => {
  it('renders Genres button', () => {
    mockAuth()
    setupQueries({ genresData: [mockGenre] })
    renderPage()
    expect(screen.getByTestId('catalog-genre-btn')).toBeInTheDocument()
  })

  it('renders sort options', () => {
    mockAuth()
    setupQueries()
    renderPage()
    expect(screen.getByRole('combobox')).toBeInTheDocument()
  })

  it('renders search input', () => {
    mockAuth()
    setupQueries()
    renderPage()
    expect(screen.getByTestId('catalog-search')).toBeInTheDocument()
  })

  it('renders status filter buttons', () => {
    mockAuth()
    setupQueries()
    renderPage()
    expect(screen.getByRole('button', { name: /airing/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /completed/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /upcoming/i })).toBeInTheDocument()
  })

  it('renders anime cards when data is loaded', () => {
    mockAuth()
    setupQueries({
      animesData: { data: [mockAnime], total: 1 },
    })
    renderPage()
    expect(screen.getByTestId('anime-card')).toBeInTheDocument()
    expect(screen.getByText('Naruto')).toBeInTheDocument()
  })

  it('shows loading indicator while fetching', () => {
    mockAuth()
    setupQueries({ animesLoading: true })
    renderPage()
    // No anime cards displayed during loading
    expect(screen.queryByTestId('anime-card')).not.toBeInTheDocument()
  })

  it('shows empty state when no animes returned', () => {
    mockAuth()
    setupQueries({ animesData: { data: [], total: 0 } })
    renderPage()
    expect(screen.getByText(/no anime found/i)).toBeInTheDocument()
  })

  it('renders genres in dropdown button', () => {
    mockAuth()
    setupQueries({ genresData: [mockGenre, { id: 2, name: 'Romance' }] })
    renderPage()
    expect(screen.getByTestId('catalog-genre-btn')).toHaveTextContent('Genres')
  })
})
