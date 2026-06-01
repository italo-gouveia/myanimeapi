import { Link } from 'react-router-dom'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import type { Anime } from '../lib/api/types'
import { addFavorite, removeFavorite } from '../lib/api/favorites'
import {
  upsertWatchlistEntry, removeWatchlistEntry,
  WATCHLIST_STATUSES, WATCHLIST_STATUS_LABELS, type WatchlistStatus,
} from '../lib/api/watchlist'

/** Placeholder shown when an anime has no cover_url yet */
function CoverPlaceholder({ title }: { title: string }) {
  const initials = title
    .split(' ')
    .slice(0, 2)
    .map((w) => w[0]?.toUpperCase() ?? '')
    .join('')
  return (
    <div className="w-full h-full bg-gradient-to-br from-brand-100 to-brand-200 flex items-center justify-center">
      <span className="text-3xl font-bold text-brand-600 select-none">{initials}</span>
    </div>
  )
}

interface AnimeCardProps {
  anime: Anime
  /** Whether the current user has this anime in their favorites */
  isFavorited?: boolean
  /** Current watchlist status, or null if not on watchlist */
  watchlistStatus?: WatchlistStatus | null
  /** Whether a user is logged in (hides action bar when false) */
  isAuthenticated?: boolean
  /** Called after a successful favorite/watchlist action */
  onAction?: (msg: string) => void
}

export function AnimeCard({
  anime,
  isFavorited = false,
  watchlistStatus = null,
  isAuthenticated = false,
  onAction,
}: AnimeCardProps) {
  const queryClient = useQueryClient()

  const favMutation = useMutation<void, Error>({
    mutationFn: async () => { isFavorited ? await removeFavorite(anime.id) : await addFavorite(anime.id) },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['favorites'] })
      onAction?.(isFavorited ? 'Removed from favorites' : '♥ Added to favorites')
    },
  })

  const watchlistMutation = useMutation<void, Error, WatchlistStatus | null>({
    mutationFn: async (status) => {
      if (status) await upsertWatchlistEntry(anime.id, status)
      else await removeWatchlistEntry(anime.id)
    },
    onSuccess: (_, status) => {
      queryClient.invalidateQueries({ queryKey: ['watchlist'] })
      onAction?.(status ? `Added to ${WATCHLIST_STATUS_LABELS[status]}` : 'Removed from watchlist')
    },
  })

  const stopProp = (e: React.MouseEvent) => e.preventDefault()

  return (
    <div
      className="group flex flex-col bg-white border border-slate-200 rounded-lg overflow-hidden hover:shadow-md hover:border-brand-300 transition-shadow"
      data-testid={`anime-card-${anime.id}`}
    >
      {/* ── Cover + info — navigates to detail ── */}
      <Link to={`/animes/${anime.id}`} className="flex flex-col flex-1 min-w-0">
        <div className="relative w-full aspect-[3/4] bg-slate-100 overflow-hidden">
          {anime.cover_url ? (
            <img
              src={anime.cover_url}
              alt={`${anime.title} cover`}
              className="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
              loading="lazy"
            />
          ) : (
            <CoverPlaceholder title={anime.title} />
          )}
          {/* Rating badge */}
          <span className="absolute top-2 right-2 inline-block bg-black/60 text-yellow-400 text-xs font-bold px-2 py-0.5 rounded backdrop-blur-sm">
            ★ {anime.rating.toFixed(1)}
          </span>
          {/* Favorited badge */}
          {isFavorited && (
            <span className="absolute top-2 left-2 text-red-400 text-lg leading-none drop-shadow">
              ♥
            </span>
          )}
          {/* Watchlist badge */}
          {watchlistStatus && (
            <span className="absolute bottom-2 left-2 text-[10px] font-semibold bg-black/60 text-white px-1.5 py-0.5 rounded backdrop-blur-sm">
              {WATCHLIST_STATUS_LABELS[watchlistStatus]}
            </span>
          )}
        </div>

        <div className="p-3 flex flex-col gap-1">
          <h3 className="font-semibold text-slate-900 line-clamp-2 text-sm leading-snug">
            {anime.title}
          </h3>
          <p className="text-xs text-slate-500">
            {anime.status} · {anime.episodes} ep
          </p>
        </div>
      </Link>

      {/* ── Action bar — visible only when authenticated ── */}
      {isAuthenticated && (
        <div
          className="flex items-center gap-2 px-3 py-2 border-t border-slate-100 bg-slate-50"
          onClick={stopProp}
        >
          {/* Favorite toggle */}
          <button
            type="button"
            onClick={() => favMutation.mutate()}
            disabled={favMutation.isPending}
            className={`flex items-center gap-1 text-xs font-medium px-2 py-1 rounded-md transition-colors disabled:opacity-50 ${
              isFavorited
                ? 'text-red-500 bg-red-50 hover:bg-red-100'
                : 'text-slate-500 hover:text-red-500 hover:bg-red-50'
            }`}
            aria-label={isFavorited ? 'Remove from favorites' : 'Add to favorites'}
          >
            {isFavorited ? '♥' : '♡'} {isFavorited ? 'Saved' : 'Save'}
          </button>

          {/* Watchlist select */}
          <select
            value={watchlistStatus ?? ''}
            disabled={watchlistMutation.isPending}
            onChange={(e) => {
              const val = e.target.value
              watchlistMutation.mutate(val ? (val as WatchlistStatus) : null)
            }}
            className="flex-1 text-xs border border-slate-200 rounded-md px-2 py-1 bg-white text-slate-600 focus:outline-none focus:ring-1 focus:ring-brand-500 disabled:opacity-50 truncate"
            aria-label="Watchlist status"
          >
            <option value="">{watchlistStatus ? '✕ Remove' : '+ Watchlist'}</option>
            {WATCHLIST_STATUSES.map((s) => (
              <option key={s} value={s}>{WATCHLIST_STATUS_LABELS[s]}</option>
            ))}
          </select>
        </div>
      )}
    </div>
  )
}
