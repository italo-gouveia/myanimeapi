import { render, screen } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { AnimeDetailPage } from './AnimeDetailPage'
import type { Anime } from '../lib/api/types'

// ── Mock react-router-dom useParams ──────────────────────────────────────────
// We wrap with MemoryRouter + Routes so useParams receives the real :id param.

// ── Mock TanStack Query ───────────────────────────────────────────────────────
type UseQueryOptions = {
  queryKey?: unknown[]
  queryFn?: () => unknown
  enabled?: boolean
}

let mockAnimeQueryResult: {
  data: Anime | undefined
  isLoading: boolean
  isError: boolean
  error: Error | null
} = { data: undefined, isLoading: false, isError: false, error: null }

vi.mock('@tanstack/react-query', () => ({
  useQuery: (opts: UseQueryOptions) => {
    // The anime query uses queryKey ['anime', animeId]
    if (Array.isArray(opts.queryKey) && opts.queryKey[0] === 'anime') {
      return mockAnimeQueryResult
    }
    // favorites / watchlist queries — return empty/disabled state
    return { data: undefined, isLoading: false, isError: false, error: null }
  },
  useMutation: () => ({
    mutate: vi.fn(),
    mutateAsync: vi.fn(),
    isPending: false,
    isError: false,
    error: null,
  }),
  useQueryClient: () => ({
    invalidateQueries: vi.fn(),
  }),
}))

// ── Mock API modules ──────────────────────────────────────────────────────────
vi.mock('../lib/api/animes', () => ({ getAnime: vi.fn() }))
vi.mock('../lib/api/favorites', () => ({
  listFavorites: vi.fn(),
  addFavorite: vi.fn(),
  removeFavorite: vi.fn(),
}))
vi.mock('../lib/api/reviews', () => ({
  createReview: vi.fn(),
  updateReview: vi.fn(),
  deleteReview: vi.fn(),
}))
vi.mock('../lib/api/watchlist', () => ({
  getWatchlist: vi.fn(),
  upsertWatchlistEntry: vi.fn(),
  removeWatchlistEntry: vi.fn(),
  WATCHLIST_STATUSES: ['plan_to_watch', 'watching', 'completed', 'dropped'],
  WATCHLIST_STATUS_LABELS: {
    plan_to_watch: 'Plan to Watch',
    watching: 'Watching',
    completed: 'Completed',
    dropped: 'Dropped',
  },
}))

// ── Mock useAuth ──────────────────────────────────────────────────────────────
let mockAuthValue = {
  isAuthenticated: false,
  userId: null as number | null,
  token: null as string | null,
  username: null as string | null,
  isAdmin: false,
  role: 'user',
  setSession: vi.fn(),
  clearSession: vi.fn(),
}

vi.mock('../lib/auth/AuthProvider', () => ({
  useAuth: () => mockAuthValue,
}))

// ── Helpers ───────────────────────────────────────────────────────────────────
function renderPage(animeId: string | number = '1') {
  return render(
    <MemoryRouter initialEntries={[`/animes/${animeId}`]}>
      <Routes>
        <Route path="/animes/:id" element={<AnimeDetailPage />} />
        <Route path="/" element={<div>Catalog</div>} />
      </Routes>
    </MemoryRouter>,
  )
}

const mockAnime: Anime = {
  id: 1,
  title: 'Cowboy Bebop',
  description: 'A ragtag crew of bounty hunters chase down criminals.',
  rating: 8.9,
  episodes: 26,
  status: 'Finished Airing',
  start_date: '1998-04-03T00:00:00Z',
  end_date: '1999-04-24T00:00:00Z',
  cover_url: 'https://example.com/cowboy-bebop.jpg',
  created_at: '2024-01-01T00:00:00Z',
  updated_at: '2024-01-01T00:00:00Z',
  reviews: [],
  genres: [{ id: 1, name: 'Action' }],
  tags: [{ id: 1, name: 'space' }],
}

// ── Tests ─────────────────────────────────────────────────────────────────────
describe('AnimeDetailPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    // Reset to unauthenticated, no-data state before each test
    mockAuthValue = {
      isAuthenticated: false,
      userId: null,
      token: null,
      username: null,
      isAdmin: false,
      role: 'user',
      setSession: vi.fn(),
      clearSession: vi.fn(),
    }
    mockAnimeQueryResult = { data: undefined, isLoading: false, isError: false, error: null }
  })

  // ── Test 1: loading state ────────────────────────────────────────────────────
  it('renders the loading state while the anime query is pending', () => {
    mockAnimeQueryResult = { data: undefined, isLoading: true, isError: false, error: null }

    renderPage('1')

    expect(screen.getByText(/loading/i)).toBeInTheDocument()
    expect(screen.queryByTestId('anime-detail')).not.toBeInTheDocument()
  })

  // ── Test 2: anime details rendered ──────────────────────────────────────────
  it('renders title, rating, and description when data is available', () => {
    mockAnimeQueryResult = { data: mockAnime, isLoading: false, isError: false, error: null }

    renderPage('1')

    // Title
    expect(screen.getByTestId('anime-detail-title')).toHaveTextContent('Cowboy Bebop')

    // Rating displayed with one decimal place
    expect(screen.getByText(/8\.9/)).toBeInTheDocument()

    // Description / synopsis
    expect(screen.getByText(/ragtag crew of bounty hunters/i)).toBeInTheDocument()

    // Genre and tag chips
    expect(screen.getByTestId('anime-detail-genre')).toHaveTextContent('Action')
    expect(screen.getByTestId('anime-detail-tag')).toHaveTextContent('#space')
  })

  // ── Test 3: Number.NaN – invalid ID shows error instead of crashing ─────────
  it('shows an error message for a non-numeric anime ID (Number.NaN path)', () => {
    // Pass a non-numeric segment — Number.parseInt('abc', 10) === NaN,
    // which equals Number.NaN. The page should guard this with Number.isFinite.
    renderPage('abc')

    const alert = screen.getByRole('alert')
    expect(alert).toHaveTextContent(/invalid anime id/i)

    // The anime detail article must NOT be rendered
    expect(screen.queryByTestId('anime-detail')).not.toBeInTheDocument()
  })

  // ── Test 4: null / undefined anime data handled gracefully ──────────────────
  it('shows an error UI when the query resolves with no anime data', () => {
    // Simulate a query that finished without data (network error or 404)
    mockAnimeQueryResult = {
      data: undefined,
      isLoading: false,
      isError: true,
      error: new Error('Not found'),
    }

    renderPage('999')

    const errorContainer = screen.getByTestId('anime-detail-error')
    expect(errorContainer).toBeInTheDocument()
    expect(errorContainer).toHaveTextContent(/could not load this anime/i)

    // A back-to-catalog link should be present
    const backLink = screen.getByRole('link', { name: /back to catalog/i })
    expect(backLink).toBeInTheDocument()
  })
})
