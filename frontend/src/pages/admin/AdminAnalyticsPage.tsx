import { useQuery } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { getAdminAnalytics, type AnalyticsRow, type DailyCount } from '../../lib/api/export'

function SkeletonRow() {
  return (
    <div className="flex items-center gap-3 py-2">
      <div className="h-3 w-4 bg-slate-100 rounded animate-pulse" />
      <div className="h-3 flex-1 bg-slate-100 rounded animate-pulse" />
      <div className="h-3 w-10 bg-slate-100 rounded animate-pulse" />
    </div>
  )
}

interface TopListProps {
  title: string
  rows: AnalyticsRow[]
  isLoading: boolean
  extra?: (row: AnalyticsRow) => React.ReactNode
}

function TopList({ title, rows, isLoading, extra }: TopListProps) {
  return (
    <div className="bg-white border border-slate-200 rounded-xl p-5 flex flex-col gap-3">
      <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">{title}</p>
      {isLoading
        ? Array.from({ length: 5 }).map((_, i) => <SkeletonRow key={i} />)
        : rows.length === 0
          ? <p className="text-sm text-slate-400">No data yet.</p>
          : rows.map((row, i) => (
              <div key={row.anime_id} className="flex items-center gap-3 text-sm">
                <span className="text-slate-400 tabular-nums w-4 text-right">{i + 1}.</span>
                <Link
                  to={`/animes/${row.anime_id}`}
                  className="flex-1 font-medium text-slate-800 hover:text-brand-700 truncate transition-colors"
                >
                  {row.title}
                </Link>
                <div className="flex items-center gap-2 shrink-0">
                  {extra?.(row)}
                  <span className="bg-slate-100 text-slate-600 text-xs font-semibold px-2 py-0.5 rounded-full">
                    {row.count}
                  </span>
                </div>
              </div>
            ))
      }
    </div>
  )
}

interface SparkBarProps {
  days: DailyCount[]
  isLoading: boolean
  label: string
  color: string
}

function SparkBar({ days, isLoading, label, color }: SparkBarProps) {
  const max = Math.max(...days.map((d) => d.count), 1)
  return (
    <div className="bg-white border border-slate-200 rounded-xl p-5 flex flex-col gap-3">
      <p className="text-xs font-semibold uppercase tracking-wider text-slate-400">{label} — last 7 days</p>
      {isLoading ? (
        <div className="flex items-end gap-1.5 h-16">
          {Array.from({ length: 7 }).map((_, i) => (
            <div key={i} className="flex-1 bg-slate-100 rounded animate-pulse" style={{ height: `${30 + i * 5}%` }} />
          ))}
        </div>
      ) : days.length === 0 ? (
        <p className="text-sm text-slate-400">No data yet.</p>
      ) : (
        <div className="flex items-end gap-1.5 h-20">
          {[...days].reverse().map((d) => {
            const pct = Math.round((d.count / max) * 100)
            const dateLabel = new Date(d.date + 'T00:00:00').toLocaleDateString('en-US', { month: 'short', day: 'numeric' })
            return (
              <div key={d.date} className="flex-1 flex flex-col items-center gap-1 group">
                <span className="text-[10px] text-slate-400 opacity-0 group-hover:opacity-100 transition-opacity tabular-nums">
                  {d.count}
                </span>
                <div
                  className={`w-full rounded-t transition-all duration-300 ${color}`}
                  style={{ height: `${Math.max(pct, 4)}%` }}
                  title={`${dateLabel}: ${d.count}`}
                />
                <span className="text-[9px] text-slate-400 leading-none">{dateLabel.split(' ')[0]}</span>
              </div>
            )
          })}
        </div>
      )}
    </div>
  )
}

export function AdminAnalyticsPage() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['admin', 'analytics'],
    queryFn: getAdminAnalytics,
    staleTime: 60_000,
  })

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold">Analytics</h1>

      {isError && <p className="text-red-600" role="alert">Failed to load analytics.</p>}

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <SparkBar
          label="Reviews"
          days={data?.recent_reviews_per_day ?? []}
          isLoading={isLoading}
          color="bg-brand-500"
        />
        <SparkBar
          label="Signups"
          days={data?.recent_signups_per_day ?? []}
          isLoading={isLoading}
          color="bg-green-500"
        />
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
        <TopList
          title="Most favorited"
          rows={data?.top_favorited ?? []}
          isLoading={isLoading}
          extra={() => <span className="text-red-400 text-xs">♥</span>}
        />
        <TopList
          title="Most reviewed"
          rows={data?.top_reviewed ?? []}
          isLoading={isLoading}
          extra={(row) =>
            row.avg_rating != null ? (
              <span className="text-yellow-500 text-xs font-semibold">★ {row.avg_rating}</span>
            ) : null
          }
        />
      </div>
    </div>
  )
}
