import { apiClient } from './client'
import type { Character, PagedResponse } from './types'

export interface ListCharactersParams {
  page?: number
  limit?: number
}

export async function listCharacters(params: ListCharactersParams = {}): Promise<PagedResponse<Character>> {
  const { data } = await apiClient.get<PagedResponse<Character>>('/v1/characters', { params })
  return data
}

export async function searchCharacters(
  name: string,
  params: ListCharactersParams = {},
): Promise<PagedResponse<Character>> {
  const { data } = await apiClient.get<PagedResponse<Character>>('/v1/characters/search', {
    params: { name, ...params },
  })
  return data
}

export async function getCharacter(id: number): Promise<Character> {
  const { data } = await apiClient.get<Character>(`/v1/characters/${id}`)
  return data
}

export async function listCharactersByAnime(
  animeId: number,
  params: ListCharactersParams = {},
): Promise<PagedResponse<Character>> {
  const { data } = await apiClient.get<PagedResponse<Character>>(
    `/v1/animes/${animeId}/characters`,
    { params },
  )
  return data
}

export interface CreateCharacterPayload {
  name: string
  description?: string
  voice_actor?: string
  image_url?: string
}

export async function createCharacter(payload: CreateCharacterPayload): Promise<Character> {
  const { data } = await apiClient.post<Character>('/v1/characters', payload)
  return data
}

export async function updateCharacter(id: number, payload: Partial<CreateCharacterPayload>): Promise<Character> {
  const { data } = await apiClient.put<Character>(`/v1/characters/${id}`, payload)
  return data
}

export async function deleteCharacter(id: number): Promise<void> {
  await apiClient.delete(`/v1/characters/${id}`)
}

export async function addCharactersToAnime(animeId: number, characterIds: number[]): Promise<void> {
  await apiClient.post(`/v1/animes/${animeId}/characters`, { character_ids: characterIds })
}

export async function removeCharacterFromAnime(animeId: number, characterId: number): Promise<void> {
  await apiClient.delete(`/v1/animes/${animeId}/characters/${characterId}`)
}
