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
  username?: string
  animeId?: number
  anime_id?: number
  created_at: string
  updated_at: string
}

export interface Character {
  id: number
  name: string
  description?: string
  voice_actor?: string
  image_url?: string
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
  cover_url?: string
  mal_id?: number
  created_at: string
  updated_at: string
  reviews?: Review[]
  genres?: Genre[]
  tags?: Tag[]
  characters?: Character[]
}

// AnimeResponse is the lean shape returned in nested contexts (watchlist, etc.)
export type AnimeResponse = Omit<Anime, 'reviews'>

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
