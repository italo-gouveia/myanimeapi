import { apiClient } from './client'

export interface ReviewPayload {
  animeId: number
  content: string
  rating: number
}

export interface ReviewUpdatePayload {
  content?: string
  rating?: number
}

export interface ReviewResponse {
  id: number
  animeId: number
  userId: number
  content: string
  rating: number
  created_at: string
  updated_at: string
}

export async function createReview(payload: ReviewPayload): Promise<ReviewResponse> {
  const res = await apiClient.post<ReviewResponse>('/v1/reviews', payload)
  return res.data
}

export async function updateReview(id: number, payload: ReviewUpdatePayload): Promise<ReviewResponse> {
  const res = await apiClient.put<ReviewResponse>(`/v1/reviews/${id}`, payload)
  return res.data
}

export async function deleteReview(id: number): Promise<void> {
  await apiClient.delete(`/v1/reviews/${id}`)
}
