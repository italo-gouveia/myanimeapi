import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AnimeCard } from './AnimeCard'
import type { Anime } from '../lib/api/types'

// Mock TanStack Query
vi.mock('@tanstack/react-query', () => ({
  useMutation: () => ({
    mutate: vi.fn(),
    isPending: false,
  }),
  useQueryClient: () => ({
    invalidateQueries: vi.fn(),
  }),
}))

// Mock favorites API
vi.mock('../lib/api/favorites', () => ({
  addFavorite: vi.fn(),
  removeFavorite: vi.fn(),
}))

// Mock watchlist API
vi.mock('../lib/api/watchlist', () => ({
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

const mockAnime: Anime = {
  id: 42,
  title: 'Fullmetal Alchemist',
  cover_url: 'https://example.com/fma.jpg',
  rating: 9.1,
  status: 'Completed',
  episodes: 64,
}

function renderCard(props: Partial<React.ComponentProps<typeof AnimeCard>> = {}) {
  return render(
    <MemoryRouter>
      <AnimeCard anime={mockAnime} {...props} />
    </MemoryRouter>,
  )
}

describe('AnimeCard', () => {
  beforeEach(() => {
    vi.clearAllMocks()
  })

  it('renders anime title, rating, status, and episodes', () => {
    renderCard()

    expect(screen.getByText('Fullmetal Alchemist')).toBeInTheDocument()
    expect(screen.getByText(/9\.1/)).toBeInTheDocument()
    // status and episodes appear in the subtitle line
    expect(screen.getAllByText(/Completed/i).length).toBeGreaterThan(0)
    expect(screen.getByText(/64 ep/)).toBeInTheDocument()
  })

  it('does NOT show the action bar when unauthenticated', () => {
    renderCard({ isAuthenticated: false })

    expect(screen.queryByRole('button', { name: /favorites/i })).not.toBeInTheDocument()
    expect(screen.queryByRole('combobox', { name: /watchlist/i })).not.toBeInTheDocument()
  })

  it('shows the action bar when authenticated', () => {
    renderCard({ isAuthenticated: true })

    expect(
      screen.getByRole('button', { name: /add to favorites|remove from favorites/i }),
    ).toBeInTheDocument()
  })

  it('renders the watchlist select when authenticated', () => {
    renderCard({ isAuthenticated: true })

    const select = screen.getByRole('combobox', { name: /watchlist status/i })
    expect(select).toBeInTheDocument()
  })

  it('renders the cover image with correct src and alt', () => {
    renderCard()

    const img = screen.getByRole('img', { name: /fullmetal alchemist cover/i })
    expect(img).toBeInTheDocument()
    expect(img).toHaveAttribute('src', 'https://example.com/fma.jpg')
    expect(img).toHaveAttribute('alt', 'Fullmetal Alchemist cover')
  })

  it('links to the anime detail page', () => {
    renderCard()

    const link = screen.getByRole('link')
    expect(link).toHaveAttribute('href', '/animes/42')
  })
})
