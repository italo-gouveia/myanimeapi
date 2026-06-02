import { apiClient } from './client'
import type { Anime } from './types'

export interface Favorite {
  id: number
  user_id: number
  anime_id: number
  created_at: string
  updated_at: string
  anime?: Anime
}

export async function listFavorites(): Promise<Favorite[]> {
  // GET /v1/favorites returns the array directly, no envelope.
  const { data } = await apiClient.get<Favorite[]>('/v1/favorites')
  return data
}

export async function addFavorite(animeId: number): Promise<Favorite> {
  const { data } = await apiClient.post<Favorite>(`/v1/favorites/${animeId}`)
  return data
}

export async function removeFavorite(animeId: number): Promise<void> {
  await apiClient.delete(`/v1/favorites/${animeId}`)
}
