import { apiClient } from './client'

export interface Genre {
  id: number
  name: string
}

export async function listGenres(): Promise<Genre[]> {
  const { data } = await apiClient.get<Genre[]>('/v1/genres?limit=100')
  return data
}
