import { describe, it, vi, beforeEach } from 'vitest'
import { render, act } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AnimeCard } from './AnimeCard'
import type { Anime } from '../lib/api/types'

vi.mock('@tanstack/react-query', () => ({
  useMutation: vi.fn(),
  useQueryClient: () => ({ invalidateQueries: vi.fn().mockResolvedValue(undefined) }),
}))

vi.mock('../lib/api/favorites', () => ({
  addFavorite: vi.fn(),
  removeFavorite: vi.fn(),
}))

vi.mock('../lib/api/watchlist', () => ({
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
  WATCHLIST_STATUS_COLORS: {},
}))

import { useMutation } from '@tanstack/react-query'

const mockAnime = { id: 1, title: 'Naruto', rating: 8.5, episodes: 220, status: 'Completed' } as unknown as Anime

describe('AnimeCard onSuccess callbacks', () => {
  beforeEach(() => vi.clearAllMocks())

  it('covers favMutation and watchlistMutation onSuccess', async () => {
    const callbacks: Array<((...args: any[]) => any) | undefined> = []
    vi.mocked(useMutation).mockImplementation(({ onSuccess }: any = {}) => {
      callbacks.push(onSuccess)
      return { mutate: vi.fn(), isPending: false, isError: false, error: null } as unknown as ReturnType<typeof useMutation>
    })
    render(
      <MemoryRouter>
        <AnimeCard anime={mockAnime} isAuthenticated />
      </MemoryRouter>,
    )
    await act(async () => {
      // favMutation onSuccess
      await callbacks[0]?.(undefined, undefined, undefined)
      // watchlistMutation onSuccess
      await callbacks[1]?.(undefined, 'watching', undefined)
    })
  })
})
