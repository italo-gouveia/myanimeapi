import { useState, useEffect, useCallback } from 'react'
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

/** Brief inline success flash — shows a label for `duration`ms then clears */
function useFlash(duration = 1400) {
  const [flash, setFlash] = useState<string | null>(null)
  const trigger = useCallback((msg: string) => {
    setFlash(msg)
    const id = window.setTimeout(() => setFlash(null), duration)
    return () => clearTimeout(id)
  }, [duration])
  return { flash, trigger }
}

interface AnimeCardProps {
  anime: Anime
  isFavorited?: boolean
  watchlistStatus?: WatchlistStatus | null
  isAuthenticated?: boolean
}

export function AnimeCard({
  anime,
  isFavorited = false,
  watchlistStatus = null,
  isAuthenticated = false,
}: AnimeCardProps) {
  const queryClient         = useQueryClient()
  const favFlash            = useFlash()
  const watchFlash          = useFlash()

  // Clean up flash timers on unmount
  useEffect(() => () => { /* timers self-clean */ }, [])

  const favMutation = useMutation<void, Error>({
    mutationFn: async () => {
      if (isFavorited) await removeFavorite(anime.id)
      else await addFavorite(anime.id)
    },
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['favorites'] })
      favFlash.trigger(isFavorited ? 'Removed' : '✓ Saved!')
    },
  })

  const watchlistMutation = useMutation<void, Error, WatchlistStatus | null>({
    mutationFn: async (status) => {
      if (status) await upsertWatchlistEntry(anime.id, status)
      else await removeWatchlistEntry(anime.id)
    },
    onSuccess: (_, status) => {
      queryClient.invalidateQueries({ queryKey: ['watchlist'] })
      watchFlash.trigger(status ? `✓ ${WATCHLIST_STATUS_LABELS[status]}` : '✓ Removed')
    },
  })

  return (
    <div
      className="group flex flex-col bg-white border border-slate-200 rounded-lg overflow-hidden hover:shadow-md hover:border-brand-300 transition-shadow"
      data-testid={`anime-card-${anime.id}`}
    >
      {/* ── Cover + info → navigates to detail ── */}
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
            <span className="absolute top-2 left-2 text-red-400 text-lg leading-none drop-shadow">♥</span>
          )}

          {/* Anime status badge (Airing / Completed / Upcoming) */}
          <span className="absolute bottom-2 left-2 text-[10px] font-semibold bg-black/60 text-white px-1.5 py-0.5 rounded backdrop-blur-sm">
            {anime.status}
          </span>
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

      {/* ── Action bar (auth only) — not part of the link ── */}
      {isAuthenticated && (
        <div
          className="flex items-center gap-2 px-3 py-2 border-t border-slate-100 bg-slate-50"
        >
          {/* Favorite toggle with inline flash */}
          <button
            type="button"
            onClick={(e) => { e.stopPropagation(); favMutation.mutate() }}
            disabled={favMutation.isPending}
            className={`
              relative flex items-center gap-1 text-xs font-medium px-2 py-1 rounded-md
              transition-all duration-200 disabled:opacity-50 min-w-[64px] justify-center
              ${favFlash.flash
                ? 'bg-green-500 text-white scale-105'
                : isFavorited
                  ? 'text-red-500 bg-red-50 hover:bg-red-100'
                  : 'text-slate-500 hover:text-red-500 hover:bg-red-50'
              }
            `}
            aria-label={isFavorited ? 'Remove from favorites' : 'Add to favorites'}
          >
            {favFlash.flash
              ? favFlash.flash
              : (isFavorited ? '♥ Saved' : '♡ Save')
            }
          </button>

          {/* Watchlist dropdown with inline flash */}
          <div className="flex-1 relative">
            {watchFlash.flash ? (
              <div className="w-full text-xs font-medium text-center py-1 rounded-md bg-green-500 text-white transition-all duration-200 scale-105">
                {watchFlash.flash}
              </div>
            ) : (
              <select
                value={watchlistStatus ?? ''}
                disabled={watchlistMutation.isPending}
                onClick={(e) => e.stopPropagation()}
                onChange={(e) => {
                  const val = e.target.value
                  watchlistMutation.mutate(val ? (val as WatchlistStatus) : null)
                }}
                className="w-full text-xs border border-slate-200 rounded-md px-2 py-1 bg-white text-slate-600 focus:outline-none focus:ring-1 focus:ring-brand-500 disabled:opacity-50"
                aria-label="Watchlist status"
              >
                <option value="">{watchlistStatus ? '✕ Remove' : '+ Watchlist'}</option>
                {WATCHLIST_STATUSES.map((s) => (
                  <option key={s} value={s}>{WATCHLIST_STATUS_LABELS[s]}</option>
                ))}
              </select>
            )}
          </div>
        </div>
      )}
    </div>
  )
}
