import { apiClient } from './client'

export interface AdminStats {
  users_total: number
  anime_total: number
  with_cover_total: number
  reviews_total: number
  favorites_total: number
  genres_total: number
  tags_total: number
}

export interface AdminUser {
  id: number
  username: string
  email: string
  role: string
  is_admin: boolean
  is_active: boolean
  created_at: string
}

export interface AdminUsersResponse {
  users: AdminUser[]
  total: number
  page: number
  limit: number
}

export interface ETLSyncResult {
  imported: number
  updated: number
  skipped: number
  errors: number
  pages: number
}

export async function getAdminStats(): Promise<AdminStats> {
  const res = await apiClient.get<{ data: AdminStats }>('/v1/admin/stats')
  return res.data.data
}

export async function listAdminUsers(page = 1, limit = 20): Promise<AdminUsersResponse> {
  const res = await apiClient.get<{ data: AdminUsersResponse }>('/v1/admin/users', {
    params: { page, limit },
  })
  return res.data.data
}

export async function updateUserRole(userId: number, role: string): Promise<void> {
  await apiClient.put(`/v1/admin/users/${userId}/role`, { role })
}

export async function triggerETLSync(pages: number): Promise<ETLSyncResult> {
  const res = await apiClient.post<{ data: ETLSyncResult }>(`/v1/etl/sync?pages=${pages}`)
  return res.data.data
}
