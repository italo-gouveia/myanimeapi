import axios from 'axios';
import type { User, Anime, Review, UserCredentials, UserCreateRequest, AnimeCreateRequest, ReviewCreateRequest, AuthResponse } from '@/types/models';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/v1';

const api = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Add token to requests if it exists
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Auth endpoints
export const auth = {
  login: (credentials: UserCredentials) => 
    api.post<AuthResponse>('/auth/authenticate', credentials),
  register: (userData: UserCreateRequest) => 
    api.post<User>('/auth/register', userData),
};

// User endpoints
export const users = {
  getAll: () => api.get<User[]>('/users'),
  getById: (id: number) => api.get<User>(`/users/${id}`),
  create: (userData: UserCreateRequest) => api.post<User>('/users', userData),
  update: (id: number, userData: Partial<UserCreateRequest> & { currentPassword?: string; newPassword?: string }) =>
    api.put<User>(`/users/${id}`, userData),
  delete: (id: number) => api.delete(`/users/${id}`),
};

// Anime endpoints
export const anime = {
  getAll: () => api.get<Anime[]>('/anime'),
  getById: (id: number) => api.get<Anime>(`/anime/${id}`),
  create: (animeData: AnimeCreateRequest) => 
    api.post<Anime>('/anime', animeData),
  update: (id: number, animeData: Partial<AnimeCreateRequest>) => 
    api.put<Anime>(`/anime/${id}`, animeData),
  delete: (id: number) => api.delete(`/anime/${id}`),
};

// Review endpoints
export const reviews = {
  getAll: () => api.get<Review[]>('/reviews'),
  getById: (id: number) => api.get<Review>(`/reviews/${id}`),
  create: (reviewData: ReviewCreateRequest) => 
    api.post<Review>('/reviews', reviewData),
  update: (id: number, reviewData: Partial<ReviewCreateRequest>) => 
    api.put<Review>(`/reviews/${id}`, reviewData),
  delete: (id: number) => api.delete(`/reviews/${id}`),
  getByAnimeId: (animeId: number) => 
    api.get<Review[]>(`/reviews/anime/${animeId}`),
  getByUserId: (userId: number) => 
    api.get<Review[]>(`/reviews/user/${userId}`),
}; 