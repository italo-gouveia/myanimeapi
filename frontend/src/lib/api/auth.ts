import { apiClient } from './client'
import type { ApiResponse, AuthData, User } from './types'

export interface LoginRequest {
  username: string
  password: string
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
}

export async function login(body: LoginRequest): Promise<AuthData> {
  const { data } = await apiClient.post<ApiResponse<AuthData>>('/v1/users/login', body)
  return data.data
}

export async function register(body: RegisterRequest): Promise<User> {
  const { data } = await apiClient.post<ApiResponse<User>>('/v1/users/register', body)
  return data.data
}

export async function getProfile(): Promise<User> {
  const { data } = await apiClient.get<ApiResponse<User>>('/v1/users/profile')
  return data.data
}
