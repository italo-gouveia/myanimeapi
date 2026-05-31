import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getAdminStats, type AdminStats } from '../../lib/api/admin'

interface StatCardProps {
  label: string
  value: number | undefined
  sub?: string
  isLoading?: boolean
}

function StatCard({ label, value, sub, isLoading }: StatCardProps) {
  return (
    <div className="bg-white border border-slate-200 rounded-xl p-5 flex flex-col gap-1">
      <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">{label}</p>
      {isLoading ? (
        <div className="h-8 w-16 bg-slate-100 rounded animate-pulse" />
      ) : (
        <p className="text-3xl font-bold text-slate-800">{value?.toLocaleString() ?? '—'}</p>
      )}
      {sub && <p className="text-xs text-slate-500">{sub}</p>}
    </div>
  )
}

export function AdminDashboard() {
  const { data: stats, isLoading, isError } = useQuery<AdminStats>({
    queryKey: ['admin', 'stats'],
    queryFn: getAdminStats,
  })

  return (
    <div className="flex flex-col gap-6">
      <div className="flex items-center justify-between">
        <h1 className="text-2xl font-bold">Dashboard</h1>
      </div>

      {isError && (
        <p className="text-red-600" role="alert">
          Failed to load stats.
        </p>
      )}

      <div className="grid grid-cols-2 sm:grid-cols-3 lg:grid-cols-4 gap-4">
        <StatCard label="Users" value={stats?.users_total} isLoading={isLoading} />
        <StatCard
          label="Anime"
          value={stats?.anime_total}
          sub={`${stats?.with_cover_total ?? 0} with cover`}
          isLoading={isLoading}
        />
        <StatCard label="Reviews" value={stats?.reviews_total} isLoading={isLoading} />
        <StatCard label="Favorites" value={stats?.favorites_total} isLoading={isLoading} />
        <StatCard label="Genres" value={stats?.genres_total} isLoading={isLoading} />
        <StatCard label="Tags" value={stats?.tags_total} isLoading={isLoading} />
      </div>

      <div className="bg-white border border-slate-200 rounded-xl p-5">
        <h2 className="text-sm font-semibold text-slate-700 mb-3">Quick actions</h2>
        <div className="flex flex-wrap gap-3">
          <Link
            to="/admin/etl"
            className="px-4 py-2 rounded-lg bg-brand-600 text-white text-sm font-medium hover:bg-brand-700 transition-colors"
          >
            🔄 Run ETL sync
          </Link>
          <Link
            to="/admin/users"
            className="px-4 py-2 rounded-lg border border-slate-200 text-slate-700 text-sm font-medium hover:bg-slate-50 transition-colors"
          >
            👥 Manage users
          </Link>
        </div>
      </div>
    </div>
  )
}
