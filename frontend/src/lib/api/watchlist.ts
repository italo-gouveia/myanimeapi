import { apiClient } from './client'
import type { AnimeResponse } from './types'

export type WatchlistStatus = 'plan_to_watch' | 'watching' | 'completed' | 'dropped' | 'on_hold'

export const WATCHLIST_STATUS_LABELS: Record<WatchlistStatus, string> = {
  plan_to_watch: 'Plan to Watch',
  watching:      'Watching',
  completed:     'Completed',
  dropped:       'Dropped',
  on_hold:       'On Hold',
}

export const WATCHLIST_STATUS_COLORS: Record<WatchlistStatus, string> = {
  plan_to_watch: 'bg-slate-100 text-slate-700',
  watching:      'bg-blue-100 text-blue-700',
  completed:     'bg-green-100 text-green-700',
  dropped:       'bg-red-100 text-red-600',
  on_hold:       'bg-yellow-100 text-yellow-700',
}

export const WATCHLIST_STATUSES = Object.keys(WATCHLIST_STATUS_LABELS) as WatchlistStatus[]

export interface WatchlistEntry {
  id: number
  user_id: number
  anime_id: number
  status: WatchlistStatus
  created_at: string
  updated_at: string
  anime?: AnimeResponse
}

export interface WatchlistResponse {
  entries: WatchlistEntry[]
  grouped: Record<WatchlistStatus, WatchlistEntry[]>
}

export async function getWatchlist(): Promise<WatchlistResponse> {
  const res = await apiClient.get<WatchlistResponse>('/v1/watchlist')
  return res.data
}

export async function upsertWatchlistEntry(animeId: number, status: WatchlistStatus): Promise<WatchlistEntry> {
  const res = await apiClient.put<WatchlistEntry>(`/v1/watchlist/${animeId}`, { status })
  return res.data
}

export async function removeWatchlistEntry(animeId: number): Promise<void> {
  await apiClient.delete(`/v1/watchlist/${animeId}`)
}
