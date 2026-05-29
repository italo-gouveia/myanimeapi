export interface ApiResponse<T> {
  status: 'success' | 'error'
  message: string
  data: T
}

export interface AuthData {
  token: string
}

export interface User {
  id: number
  username: string
  email: string
  is_active: boolean
  is_admin: boolean
  profile_pic?: string
  bio?: string
  created_at: string
  updated_at: string
}

export interface Genre {
  id: number
  name: string
}

export interface Tag {
  id: number
  name: string
}

export interface Review {
  id: number
  content: string
  rating: number
  userId?: number
  user_id?: number
  animeId?: number
  anime_id?: number
  created_at: string
  updated_at: string
}

export interface Anime {
  id: number
  title: string
  description: string
  rating: number
  episodes: number
  status: string
  start_date: string
  end_date: string
  created_at: string
  updated_at: string
  reviews?: Review[]
  genres?: Genre[]
  tags?: Tag[]
}

export interface PagedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}

export interface ApiError {
  code?: string
  message?: string
  details?: string
}
