import { apiClient } from './client'

export type ExportFormat = 'json' | 'csv'

export async function exportUserData(format: ExportFormat = 'json'): Promise<void> {
  const res = await apiClient.get(`/v1/users/export`, {
    params: { format },
    responseType: 'blob',
  })

  const blob = new Blob([res.data as BlobPart], {
    type: format === 'csv' ? 'text/csv' : 'application/json',
  })
  const url  = URL.createObjectURL(blob)
  const a    = document.createElement('a')
  a.href     = url
  a.download = `myanimeapi-export.${format}`
  a.click()
  URL.revokeObjectURL(url)
}

export interface AnalyticsRow {
  anime_id: number
  title: string
  count: number
  avg_rating?: number
}

export interface DailyCount {
  date: string
  count: number
}

export interface AdminAnalytics {
  top_favorited: AnalyticsRow[]
  top_reviewed: AnalyticsRow[]
  recent_reviews_per_day: DailyCount[]
  recent_signups_per_day: DailyCount[]
}

export async function getAdminAnalytics(): Promise<AdminAnalytics> {
  const res = await apiClient.get<AdminAnalytics>('/v1/admin/analytics')
  return res.data
}
