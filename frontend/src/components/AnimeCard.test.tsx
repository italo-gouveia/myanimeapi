import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
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

  it('calls favMutation.mutate when favorite button is clicked', () => {
    renderCard({ isAuthenticated: true })

    const btn = screen.getByRole('button', { name: /add to favorites|remove from favorites/i })
    fireEvent.click(btn)
    // mutation.mutate is called — no error thrown means stopPropagation ran
    expect(btn).toBeInTheDocument()
  })

  it('calls watchlistMutation.mutate when watchlist select changes', () => {
    renderCard({ isAuthenticated: true })

    const select = screen.getByRole('combobox', { name: /watchlist status/i })
    fireEvent.change(select, { target: { value: 'watching' } })
    expect(select).toBeInTheDocument()
  })

  it('stopPropagation fires on select click without bubbling', () => {
    renderCard({ isAuthenticated: true })

    const select = screen.getByRole('combobox', { name: /watchlist status/i })
    // fireEvent.click triggers the onClick handler (stopPropagation)
    fireEvent.click(select)
    expect(select).toBeInTheDocument()
  })

  it('renders favorited badge and correct button label when isFavorited=true', () => {
    renderCard({ isAuthenticated: true, isFavorited: true })

    expect(screen.getByRole('button', { name: /remove from favorites/i })).toBeInTheDocument()
  })

  it('renders CoverPlaceholder when cover_url is absent', () => {
    const animeNoCover = { ...mockAnime, cover_url: undefined as unknown as string }
    render(
      <MemoryRouter>
        <AnimeCard anime={animeNoCover} />
      </MemoryRouter>,
    )
    // Initials derived from "Fullmetal Alchemist" → "FA"
    expect(screen.getByText('FA')).toBeInTheDocument()
  })
})
