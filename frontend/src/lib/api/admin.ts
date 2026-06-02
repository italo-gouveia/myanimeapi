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

// GET /v1/admin/stats — returns AdminStats directly (no envelope)
export async function getAdminStats(): Promise<AdminStats> {
  const res = await apiClient.get<AdminStats>('/v1/admin/stats')
  return res.data
}

// GET /v1/admin/users — returns { data: AdminUser[], total, page, limit }
export async function listAdminUsers(page = 1, limit = 20): Promise<AdminUsersResponse> {
  const res = await apiClient.get<{ data: AdminUser[]; total: number; page: number; limit: number }>(
    '/v1/admin/users',
    { params: { page, limit } },
  )
  return {
    users: res.data.data,
    total: res.data.total,
    page: res.data.page,
    limit: res.data.limit,
  }
}

export async function updateUserRole(userId: number, role: string): Promise<void> {
  await apiClient.put(`/v1/admin/users/${userId}/role`, { role })
}

// POST /v1/etl/sync — returns ETLSyncResult directly (no envelope)
export async function triggerETLSync(pages: number): Promise<ETLSyncResult> {
  const res = await apiClient.post<ETLSyncResult>(`/v1/etl/sync?pages=${pages}`, null, {
    timeout: 300_000, // 5 min — ETL needs ~500ms per page + Jikan latency
  })
  return res.data
}
