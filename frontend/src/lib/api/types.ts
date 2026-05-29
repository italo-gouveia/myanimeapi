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

export interface ApiError {
  code?: string
  message?: string
  details?: string
}
