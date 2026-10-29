import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { WatchlistPage } from './WatchlistPage'

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
  useMutation: vi.fn(),
  useQueryClient: () => ({ invalidateQueries: vi.fn() }),
}))

vi.mock('../lib/api/watchlist', () => ({
  getWatchlist: vi.fn(),
  upsertWatchlistEntry: vi.fn(),
  removeWatchlistEntry: vi.fn(),
  WATCHLIST_STATUSES: ['watching', 'plan_to_watch', 'completed', 'on_hold', 'dropped'],
  WATCHLIST_STATUS_LABELS: {
    watching: 'Watching',
    plan_to_watch: 'Plan to Watch',
    completed: 'Completed',
    on_hold: 'On Hold',
    dropped: 'Dropped',
  },
  WATCHLIST_STATUS_COLORS: {
    watching: 'bg-blue-100 text-blue-700',
    plan_to_watch: 'bg-slate-100 text-slate-700',
    completed: 'bg-green-100 text-green-700',
    on_hold: 'bg-yellow-100 text-yellow-700',
    dropped: 'bg-red-100 text-red-700',
  },
}))

import { useQuery, useMutation } from '@tanstack/react-query'

function setupMutation() {
  vi.mocked(useMutation).mockReturnValue({
    mutate: vi.fn(),
    isPending: false,
    isError: false,
    error: null,
  } as unknown as ReturnType<typeof useMutation>)
}

function renderPage() {
  return render(
    <MemoryRouter>
      <WatchlistPage />
    </MemoryRouter>,
  )
}

describe('WatchlistPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setupMutation()
  })

  it('shows loading state', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: true, data: undefined, isError: false } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
  })

  it('shows error state', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: true, data: undefined } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/failed to load watchlist/i)).toBeInTheDocument()
  })

  it('shows empty state when watchlist has no entries', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { entries: [], grouped: {} },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/your watchlist is empty/i)).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /browse catalog/i })).toBeInTheDocument()
  })

  it('renders grouped watchlist entries', () => {
    const entry = {
      id: 1,
      anime_id: 5,
      status: 'watching',
      anime: { id: 5, title: 'One Piece', rating: 8.7, episodes: 1000 },
    }
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: {
        entries: [entry],
        grouped: { watching: [entry] },
      },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getAllByText('Watching').length).toBeGreaterThan(0)
    expect(screen.getByText('One Piece')).toBeInTheDocument()
  })

  it('shows total count in header', () => {
    const entry = {
      id: 1,
      anime_id: 5,
      status: 'watching',
      anime: { id: 5, title: 'One Piece', rating: 8.7, episodes: 1000 },
    }
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { entries: [entry], grouped: { watching: [entry] } },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/1 anime/)).toBeInTheDocument()
  })

  it('renders the page heading', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: { entries: [], grouped: {} } } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('heading', { name: /watchlist/i })).toBeInTheDocument()
  })

  it('covers onSuccess callbacks for move and remove mutations', async () => {
    const callbacks: Array<((...args: any[]) => any) | undefined> = []
    vi.mocked(useMutation).mockImplementation(({ onSuccess }: any = {}) => {
      callbacks.push(onSuccess)
      return { mutate: vi.fn(), isPending: false, isError: false, error: null } as unknown as ReturnType<typeof useMutation>
    })
    const entry = { id: 1, anime_id: 5, status: 'watching', anime: { id: 5, title: 'One Piece', rating: 8.7, episodes: 1000 } }
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false, isError: false,
      data: { entries: [entry], grouped: { watching: [entry] } },
    } as ReturnType<typeof useQuery>)
    renderPage()
    await act(async () => {
      for (const cb of callbacks) await cb?.(undefined, 'watching', undefined)
    })
  })
})
