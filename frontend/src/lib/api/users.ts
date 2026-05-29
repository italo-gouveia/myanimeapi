import { apiClient } from './client'
import type { ApiResponse, User } from './types'

export async function getProfile(): Promise<User> {
  const { data } = await apiClient.get<ApiResponse<User>>('/v1/users/profile')
  return data.data
}
