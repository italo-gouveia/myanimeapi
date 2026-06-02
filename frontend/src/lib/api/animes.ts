import { apiClient } from './client'
import type { Anime, PagedResponse } from './types'

export interface ListAnimesParams {
  page?: number
  limit?: number
  status?: string
  statuses?: string[]  // OR semantics — sent as comma-separated ?statuses=Airing,Completed
  genre?: string
  genres?: string[]    // AND semantics — sent as comma-separated ?genres=Action,Adventure
  sort_by?: 'title' | 'rating' | 'episodes' | 'created_at' | 'start_date'
  order?: 'asc' | 'desc'
}

/** Serialise array fields to the comma-separated format the backend expects. */
function serializeParams(params: ListAnimesParams): Record<string, unknown> {
  const { genres, statuses, ...rest } = params
  return {
    ...rest,
    ...(genres && genres.length > 0 ? { genres: genres.join(',') } : {}),
    ...(statuses && statuses.length > 0 ? { statuses: statuses.join(',') } : {}),
  }
}

export async function listAnimes(params: ListAnimesParams = {}): Promise<PagedResponse<Anime>> {
  const { data } = await apiClient.get<PagedResponse<Anime>>('/v1/animes', {
    params: serializeParams(params),
  })
  return data
}

export async function searchAnimes(
  title: string,
  params: Omit<ListAnimesParams, 'genre' | 'status'> = {},
): Promise<PagedResponse<Anime>> {
  const { data } = await apiClient.get<PagedResponse<Anime>>('/v1/animes/search', {
    params: { title, ...serializeParams(params) },
  })
  return data
}

export async function getAnime(id: number): Promise<Anime> {
  const { data } = await apiClient.get<Anime>(`/v1/animes/${id}`)
  return data
}
