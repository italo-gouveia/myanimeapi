import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { ToastContainer, useToast } from '../components/Toast'
import {
  getWatchlist,
  removeWatchlistEntry,
  upsertWatchlistEntry,
  WATCHLIST_STATUSES,
  WATCHLIST_STATUS_COLORS,
  WATCHLIST_STATUS_LABELS,
  type WatchlistEntry,
  type WatchlistStatus,
} from '../lib/api/watchlist'

const SECTION_ICONS: Record<WatchlistStatus, string> = {
  watching:      '▶',
  plan_to_watch: '📋',
  completed:     '✅',
  on_hold:       '⏸',
  dropped:       '🗑',
}

function WatchlistCard({ entry, onToast }: { entry: WatchlistEntry; onToast: (msg: string) => void }) {
  const queryClient = useQueryClient()

  const moveMutation = useMutation({
    mutationFn: (status: WatchlistStatus) => upsertWatchlistEntry(entry.anime_id, status),
    onSuccess: (_, status) => {
      queryClient.invalidateQueries({ queryKey: ['watchlist'] })
      onToast(`Moved to ${WATCHLIST_STATUS_LABELS[status]}`)
    },
  })

  const removeMutation = useMutation({
    mutationFn: () => removeWatchlistEntry(entry.anime_id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['watchlist'] })
      onToast('Removed from watchlist')
    },
  })

  const anime = entry.anime

  const isBusy = moveMutation.isPending || removeMutation.isPending

  return (
    <div className={`group flex gap-3 bg-white border border-slate-200 rounded-xl p-3 hover:border-brand-300 transition-all duration-200 animate-fade-in ${isBusy ? 'opacity-60 scale-[0.99]' : ''}`}>
      {/* Cover */}
      {anime?.cover_url ? (
        <img
          src={anime.cover_url}
          alt={anime?.title}
          className="w-12 h-16 object-cover rounded-md shrink-0"
        />
      ) : (
        <div className="w-12 h-16 bg-slate-100 rounded-md shrink-0 flex items-center justify-center text-slate-400 text-xs">
          ?
        </div>
      )}

      {/* Info */}
      <div className="flex-1 min-w-0 flex flex-col justify-between">
        <div>
          <Link
            to={`/animes/${entry.anime_id}`}
            className="text-sm font-semibold text-slate-800 hover:text-brand-700 transition-colors line-clamp-1"
          >
            {anime?.title ?? `Anime #${entry.anime_id}`}
          </Link>
          {anime && (
            <p className="text-xs text-slate-500 mt-0.5">
              {anime.episodes} ep · ★ {anime.rating?.toFixed(1)}
            </p>
          )}
        </div>

        {/* Status + actions */}
        <div className="flex items-center gap-2 mt-2">
          <select
            value={entry.status}
            disabled={moveMutation.isPending}
            onChange={(e) => moveMutation.mutate(e.target.value as WatchlistStatus)}
            className="text-xs border border-slate-200 rounded px-2 py-1 focus:outline-none focus:ring-1 focus:ring-brand-500 disabled:opacity-50"
            aria-label="Change status"
          >
            {WATCHLIST_STATUSES.map((s) => (
              <option key={s} value={s}>{WATCHLIST_STATUS_LABELS[s]}</option>
            ))}
          </select>

          <button
            type="button"
            onClick={() => removeMutation.mutate()}
            disabled={removeMutation.isPending}
            className="text-xs text-slate-400 hover:text-red-500 transition-colors disabled:opacity-50 ml-auto"
            aria-label="Remove from watchlist"
          >
            Remove
          </button>
        </div>
      </div>
    </div>
  )
}

export function WatchlistPage() {
  const { data, isLoading, isError } = useQuery({
    queryKey: ['watchlist'],
    queryFn: getWatchlist,
  })
  const { toasts, push: pushToast, remove: removeToast } = useToast()

  const total = data?.entries.length ?? 0

  return (
    <>
    <ToastContainer toasts={toasts} onDone={removeToast} />
    <div className="flex flex-col gap-6">
      <header>
        <h1 className="text-3xl font-bold">Watchlist</h1>
        {data && (
          <p className="text-slate-500 text-sm mt-1">{total} {total === 1 ? 'anime' : 'animes'}</p>
        )}
      </header>

      {isLoading && <p className="text-slate-500">Loading…</p>}
      {isError && <p className="text-red-600" role="alert">Failed to load watchlist.</p>}

      {data && total === 0 && (
        <div className="text-center py-16 flex flex-col items-center gap-3">
          <p className="text-slate-500">Your watchlist is empty.</p>
          <Link to="/" className="text-brand-700 underline text-sm">Browse catalog →</Link>
        </div>
      )}

      {data && total > 0 && (
        <div className="flex flex-col gap-8">
          {WATCHLIST_STATUSES.map((status) => {
            const entries = data.grouped[status] ?? []
            if (entries.length === 0) return null
            return (
              <section key={status}>
                <div className="flex items-center gap-2 mb-3">
                  <span>{SECTION_ICONS[status]}</span>
                  <h2 className="font-semibold text-slate-700">{WATCHLIST_STATUS_LABELS[status]}</h2>
                  <span className={`text-xs font-medium px-2 py-0.5 rounded-full ${WATCHLIST_STATUS_COLORS[status]}`}>
                    {entries.length}
                  </span>
                </div>
                <ul className="grid gap-2 sm:grid-cols-2 lg:grid-cols-3">
                  {entries.map((entry) => (
                    <li key={entry.id}>
                      <WatchlistCard entry={entry} onToast={(msg) => pushToast(msg)} />
                    </li>
                  ))}
                </ul>
              </section>
            )
          })}
        </div>
      )}
    </div>
    </>
  )
}
