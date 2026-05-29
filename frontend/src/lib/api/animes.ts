import { apiClient } from './client'
import type { Anime, PagedResponse } from './types'

export interface ListAnimesParams {
  page?: number
  limit?: number
  status?: string
  genre?: string
  sort_by?: 'title' | 'rating' | 'episodes' | 'created_at' | 'start_date'
  order?: 'asc' | 'desc'
}

export async function listAnimes(params: ListAnimesParams = {}): Promise<PagedResponse<Anime>> {
  const { data } = await apiClient.get<PagedResponse<Anime>>('/v1/animes', { params })
  return data
}

export async function searchAnimes(
  title: string,
  params: Omit<ListAnimesParams, 'genre' | 'status'> = {},
): Promise<PagedResponse<Anime>> {
  const { data } = await apiClient.get<PagedResponse<Anime>>('/v1/animes/search', {
    params: { title, ...params },
  })
  return data
}

export async function getAnime(id: number): Promise<Anime> {
  // GET /v1/animes/{id} returns the AnimeResponse directly (no envelope).
  const { data } = await apiClient.get<Anime>(`/v1/animes/${id}`)
  return data
}
