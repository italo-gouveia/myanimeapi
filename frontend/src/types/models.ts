export interface User {
  id: number;
  username: string;
  email?: string;
  is_admin: boolean;
  created_at: string;
  updated_at: string;
}

export interface Anime {
  id: number;
  title: string;
  description: string;
  rating: number;
  created_at: string;
  updated_at: string;
}

export interface Review {
  id: number;
  userId: number;
  animeId: number;
  content: string;
  rating: number;
  created_at: string;
  updated_at: string;
  user?: User;
  anime?: Anime;
}

export interface UserCredentials {
  username: string;
  password: string;
}

export interface UserCreateRequest {
  username: string;
  email: string;
  password: string;
}

export interface AnimeCreateRequest {
  title: string;
  description: string;
  rating: number;
}

export interface ReviewCreateRequest {
  userId: number;
  animeId: number;
  content: string;
  rating: number;
}

export interface AuthResponse {
  token: string;
  user: User;
} 