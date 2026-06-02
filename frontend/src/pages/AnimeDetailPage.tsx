import { useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getAnime } from '../lib/api/animes'
import { listFavorites, addFavorite, removeFavorite } from '../lib/api/favorites'
import { createReview, updateReview, deleteReview } from '../lib/api/reviews'
import {
  getWatchlist, upsertWatchlistEntry, removeWatchlistEntry,
  WATCHLIST_STATUS_LABELS, WATCHLIST_STATUSES, type WatchlistStatus,
} from '../lib/api/watchlist'
import { useAuth } from '../lib/auth/AuthProvider'
import { Button } from '../components/Button'
import { ToastContainer, useToast } from '../components/Toast'
import type { Review } from '../lib/api/types'

function formatYear(date: string): string | null {
  const year = new Date(date).getUTCFullYear()
  return Number.isFinite(year) && year > 1900 ? String(year) : null
}

// ── Helpers ──────────────────────────────────────────────────────────────────
function initials(name?: string): string {
  if (!name) return '?'
  return name.split(/[\s_-]/).map((w) => w[0]?.toUpperCase() ?? '').slice(0, 2).join('')
}

function starsFromRating(rating: number): string {
  const full = Math.round(rating / 2)
  return '★'.repeat(full) + '☆'.repeat(5 - full)
}

// ── Avatar ────────────────────────────────────────────────────────────────────
function Avatar({ name, size = 'md' }: { name?: string; size?: 'sm' | 'md' }) {
  const sz = size === 'sm' ? 'w-7 h-7 text-xs' : 'w-9 h-9 text-sm'
  return (
    <div className={`${sz} rounded-full bg-brand-600 text-white flex items-center justify-center font-semibold shrink-0`}>
      {initials(name)}
    </div>
  )
}

// ── Star rating input ─────────────────────────────────────────────────────────
function StarRating({ value, onChange }: { value: number; onChange: (n: number) => void }) {
  const [hovered, setHovered] = useState(0)
  return (
    <div className="flex items-center gap-0.5" role="group" aria-label="Rating">
      {Array.from({ length: 10 }, (_, i) => i + 1).map((n) => (
        <button
          key={n}
          type="button"
          onClick={() => onChange(n)}
          onMouseEnter={() => setHovered(n)}
          onMouseLeave={() => setHovered(0)}
          className={`text-2xl leading-none transition-colors ${n <= (hovered || value) ? 'text-yellow-400' : 'text-slate-200'}`}
          aria-label={`${n} stars`}
        >★</button>
      ))}
      <span className="ml-2 text-sm text-slate-500 tabular-nums w-12">
        {(hovered || value) > 0 ? `${hovered || value}/10` : ''}
      </span>
    </div>
  )
}

// ── Review form ───────────────────────────────────────────────────────────────
interface ReviewFormProps {
  animeId: number
  initialContent?: string
  initialRating?: number
  reviewId?: number
  onDone: () => void
  onCancel?: () => void
}

function ReviewForm({ animeId, initialContent = '', initialRating = 0, reviewId, onDone, onCancel }: ReviewFormProps) {
  const [content, setContent] = useState(initialContent)
  const [rating, setRating]   = useState(initialRating)
  const queryClient           = useQueryClient()
  const isEdit                = Boolean(reviewId)

  const mutation = useMutation({
    mutationFn: () =>
      isEdit
        ? updateReview(reviewId!, { content, rating })
        : createReview({ animeId, content, rating }),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['anime', animeId] })
      onDone()
    },
  })

  const errorMsg = (() => {
    if (!mutation.isError) return null
    const status = (mutation.error as { response?: { status?: number } })?.response?.status
    if (status === 409) return 'You already have a review for this anime — scroll up to find it.'
    return 'Failed to submit. Please try again.'
  })()

  const charLeft = 500 - content.length
  const canSubmit = rating > 0 && content.trim().length > 0 && !mutation.isPending

  return (
    <form onSubmit={(e) => { e.preventDefault(); mutation.mutate() }} className="flex flex-col gap-4">
      <StarRating value={rating} onChange={setRating} />

      <div className="flex flex-col gap-1">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          maxLength={500}
          rows={4}
          placeholder="What did you think? Share your honest opinion…"
          className="w-full rounded-xl border border-slate-200 px-4 py-3 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-brand-500 bg-white"
          required
        />
        <p className={`text-xs text-right tabular-nums ${charLeft < 50 ? 'text-orange-500' : 'text-slate-400'}`}>
          {charLeft} / 500
        </p>
      </div>

      {errorMsg && <p className="text-sm text-red-600" role="alert">{errorMsg}</p>}

      <div className="flex gap-2">
        <Button type="submit" disabled={!canSubmit}>
          {mutation.isPending ? 'Saving…' : isEdit ? 'Save changes' : 'Post review'}
        </Button>
        {onCancel && (
          <Button type="button" variant="secondary" onClick={onCancel}>Cancel</Button>
        )}
      </div>
    </form>
  )
}

// ── Review card ───────────────────────────────────────────────────────────────
interface ReviewCardProps {
  review: Review
  isOwn?: boolean
  onEdit?: () => void
  onDelete?: () => void
  isDeleting?: boolean
}

function ReviewCard({ review, isOwn, onEdit, onDelete, isDeleting }: ReviewCardProps) {
  const displayName = review.username || `User #${review.userId ?? review.user_id ?? '?'}`

  return (
    <div className={`rounded-xl p-4 flex flex-col gap-3 ${isOwn ? 'bg-brand-50 border border-brand-200' : 'bg-white border border-slate-200'}`}>
      {/* Header */}
      <div className="flex items-start justify-between gap-3">
        <div className="flex items-center gap-2.5 min-w-0">
          <Avatar name={displayName} size="sm" />
          <div className="min-w-0">
            <div className="flex items-center gap-2 flex-wrap">
              <span className="text-sm font-semibold text-slate-800 truncate">{displayName}</span>
              {isOwn && (
                <span className="text-[10px] font-bold bg-brand-600 text-white px-1.5 py-0.5 rounded-full uppercase tracking-wide">
                  You
                </span>
              )}
            </div>
            <div className="flex items-center gap-1.5 mt-0.5">
              <span className="text-yellow-400 text-xs leading-none">{starsFromRating(review.rating)}</span>
              <span className="text-xs text-slate-500 font-medium">{review.rating}/10</span>
            </div>
          </div>
        </div>

        <div className="flex items-center gap-3 shrink-0">
          <time className="text-xs text-slate-400" dateTime={review.created_at}>
            {new Date(review.created_at).toLocaleDateString('en-US', { month: 'short', day: 'numeric', year: 'numeric' })}
          </time>
          {isOwn && (
            <div className="flex gap-2">
              {onEdit && (
                <button type="button" onClick={onEdit} className="text-xs text-slate-400 hover:text-brand-700 transition-colors">
                  Edit
                </button>
              )}
              {onDelete && (
                <button
                  type="button"
                  onClick={() => { if (window.confirm('Delete your review?')) onDelete() }}
                  disabled={isDeleting}
                  className="text-xs text-slate-400 hover:text-red-600 transition-colors disabled:opacity-40"
                >
                  Delete
                </button>
              )}
            </div>
          )}
        </div>
      </div>

      {/* Content */}
      <p className="text-sm text-slate-700 whitespace-pre-line leading-relaxed">{review.content}</p>
    </div>
  )
}

// ── Main page ─────────────────────────────────────────────────────────────────
export function AnimeDetailPage() {
  const { id }                      = useParams<{ id: string }>()
  const animeId                     = id ? Number.parseInt(id, 10) : NaN
  const { isAuthenticated, userId } = useAuth()
  const queryClient                 = useQueryClient()
  const [editingReviewId, setEditingReviewId] = useState<number | null>(null)
  const [showForm, setShowForm]               = useState(false)
  const { toasts, push: pushToast, remove: removeToast } = useToast()

  const { data: anime, isLoading, isError, error } = useQuery({
    queryKey: ['anime', animeId],
    queryFn: () => getAnime(animeId),
    enabled: Number.isFinite(animeId),
  })

  const { data: favorites } = useQuery({
    queryKey: ['favorites'],
    queryFn: listFavorites,
    enabled: isAuthenticated,
  })

  const { data: watchlist } = useQuery({
    queryKey: ['watchlist'],
    queryFn: getWatchlist,
    enabled: isAuthenticated,
  })

  const isFavorited       = favorites?.some((f) => f.anime_id === animeId) ?? false
  const myWatchlistEntry  = watchlist?.entries.find((e) => e.anime_id === animeId)

  const addMutation    = useMutation({ mutationFn: () => addFavorite(animeId),    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['favorites'] }) })
  const removeMutation = useMutation({ mutationFn: () => removeFavorite(animeId), onSuccess: () => queryClient.invalidateQueries({ queryKey: ['favorites'] }) })

  const watchlistMutation = useMutation<void, Error, WatchlistStatus | null>({
    mutationFn: async (status) => {
      if (status) await upsertWatchlistEntry(animeId, status)
      else await removeWatchlistEntry(animeId)
    },
    onSuccess: (_, status) => {
      queryClient.invalidateQueries({ queryKey: ['watchlist'] })
      if (status) pushToast(`Added to ${WATCHLIST_STATUS_LABELS[status]}`)
      else pushToast('Removed from watchlist', 'info')
    },
  })

  const deleteMutation = useMutation({
    mutationFn: async (reviewId: number) => {
      try {
        await deleteReview(reviewId)
      } catch (err: unknown) {
        // 404 = already gone — treat as success (idempotent)
        if ((err as { response?: { status?: number } })?.response?.status === 404) return
        throw err
      }
    },
    // Always refresh the anime detail so stale cards disappear
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['anime', animeId] }),
    onError:   () => queryClient.invalidateQueries({ queryKey: ['anime', animeId] }),
  })

  if (!Number.isFinite(animeId)) return <p className="text-red-600" role="alert">Invalid anime id.</p>
  if (isLoading) return <p className="text-slate-500">Loading…</p>
  if (isError || !anime) {
    return (
      <div role="alert" data-testid="anime-detail-error">
        <p className="text-red-600">Could not load this anime. {(error as Error | undefined)?.message ?? ''}</p>
        <Link to="/" className="text-brand-700 underline">Back to catalog</Link>
      </div>
    )
  }

  const startYear = formatYear(anime.start_date)
  const endYear   = formatYear(anime.end_date)
  const yearRange = startYear ? (endYear && endYear !== startYear ? `${startYear} – ${endYear}` : startYear) : null

  const myReview     = isAuthenticated && userId
    ? anime.reviews?.find((r) => (r.userId ?? r.user_id) === userId)
    : undefined
  const otherReviews = (anime.reviews ?? []).filter((r) => (r.userId ?? r.user_id) !== userId)
  const totalReviews = (anime.reviews ?? []).length

  return (
    <>
    <ToastContainer toasts={toasts} onDone={removeToast} />
    <article className="flex flex-col gap-6" data-testid="anime-detail">
      {/* ── Header ── */}
      <header className="flex flex-col gap-2">
        <Link to="/" className="text-sm text-brand-700 underline w-fit">← Catalog</Link>

        <div className="flex flex-col sm:flex-row gap-6">
          {anime.cover_url && (
            <div className="shrink-0 w-full sm:w-40 md:w-52">
              <img src={anime.cover_url} alt={`${anime.title} cover`} className="w-full rounded-lg shadow-md object-cover aspect-[3/4]" />
            </div>
          )}
          <div className="flex flex-col gap-3 min-w-0">
            <div className="flex items-start justify-between gap-4">
              <h1 className="text-3xl font-bold" data-testid="anime-detail-title">{anime.title}</h1>
              {isAuthenticated && (
                <div className="flex items-center gap-2 shrink-0">
                  <select
                    value={myWatchlistEntry?.status ?? ''}
                    disabled={watchlistMutation.isPending}
                    onChange={(e) => watchlistMutation.mutate(e.target.value ? (e.target.value as WatchlistStatus) : null)}
                    className="h-9 rounded-lg border border-slate-300 bg-white px-2 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-brand-500 disabled:opacity-50"
                    aria-label="Add to watchlist"
                    data-testid="watchlist-select"
                  >
                    <option value="">+ Watchlist</option>
                    {WATCHLIST_STATUSES.map((s) => (
                      <option key={s} value={s}>{WATCHLIST_STATUS_LABELS[s]}</option>
                    ))}
                  </select>
                  <Button
                    variant={isFavorited ? 'secondary' : 'primary'}
                    disabled={addMutation.isPending || removeMutation.isPending}
                    onClick={() => isFavorited ? removeMutation.mutate() : addMutation.mutate()}
                    data-testid="favorite-toggle"
                    aria-pressed={isFavorited}
                  >
                    {isFavorited ? '★ Saved' : '☆ Save'}
                  </Button>
                </div>
              )}
            </div>
            <div className="flex items-center gap-3 text-sm text-slate-600">
              <span className="inline-flex items-center gap-1 font-semibold text-brand-700">★ {anime.rating.toFixed(1)}</span>
              <span>·</span><span>{anime.status}</span>
              <span>·</span><span>{anime.episodes} episodes</span>
              {yearRange && <><span>·</span><span>{yearRange}</span></>}
            </div>
          </div>
        </div>
      </header>

      {/* ── Synopsis ── */}
      {anime.description && (
        <section>
          <h2 className="text-lg font-semibold mb-1">Synopsis</h2>
          <p className="text-slate-700 whitespace-pre-line">{anime.description}</p>
        </section>
      )}

      {/* ── Tags/Genres ── */}
      {(anime.genres?.length || anime.tags?.length) ? (
        <section className="flex flex-wrap gap-2">
          {anime.genres?.map((g) => (
            <span key={`genre-${g.id}`} className="bg-brand-100 text-brand-700 text-xs font-semibold px-2 py-1 rounded" data-testid="anime-detail-genre">{g.name}</span>
          ))}
          {anime.tags?.map((t) => (
            <span key={`tag-${t.id}`} className="bg-slate-200 text-slate-700 text-xs px-2 py-1 rounded" data-testid="anime-detail-tag">#{t.name}</span>
          ))}
        </section>
      ) : null}

      {/* ── Reviews ── */}
      <section className="flex flex-col gap-5">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold">
            Reviews
            {totalReviews > 0 && <span className="ml-2 text-sm font-normal text-slate-400">({totalReviews})</span>}
          </h2>
        </div>

        {/* My review or write form */}
        {isAuthenticated ? (
          myReview ? (
            editingReviewId === myReview.id ? (
              /* Editing */
              <div className="flex flex-col gap-3">
                <p className="text-sm font-semibold text-slate-600">Edit your review</p>
                <ReviewForm
                  animeId={animeId}
                  reviewId={myReview.id}
                  initialContent={myReview.content}
                  initialRating={myReview.rating}
                  onDone={() => setEditingReviewId(null)}
                  onCancel={() => setEditingReviewId(null)}
                />
              </div>
            ) : (
              /* My review card */
              <ReviewCard
                review={myReview}
                isOwn
                onEdit={() => setEditingReviewId(myReview.id)}
                onDelete={() => deleteMutation.mutate(myReview.id)}
                isDeleting={deleteMutation.isPending}
              />
            )
          ) : showForm ? (
            /* Write form */
            <div className="flex flex-col gap-3 bg-slate-50 border border-slate-200 rounded-xl p-4">
              <div className="flex items-center justify-between">
                <p className="text-sm font-semibold text-slate-700">Write a review</p>
                <button type="button" onClick={() => setShowForm(false)} className="text-slate-400 hover:text-slate-600 text-lg leading-none">✕</button>
              </div>
              <ReviewForm animeId={animeId} onDone={() => setShowForm(false)} onCancel={() => setShowForm(false)} />
            </div>
          ) : (
            /* CTA to open form */
            <button
              type="button"
              onClick={() => setShowForm(true)}
              className="flex items-center gap-3 w-full text-left p-4 rounded-xl border border-dashed border-slate-300 hover:border-brand-400 hover:bg-brand-50 transition-colors group"
              data-testid="write-review-cta"
            >
              <div className="w-9 h-9 rounded-full border-2 border-dashed border-slate-300 group-hover:border-brand-400 flex items-center justify-center text-slate-400 group-hover:text-brand-600 transition-colors text-lg">
                +
              </div>
              <div>
                <p className="text-sm font-semibold text-slate-700 group-hover:text-brand-700 transition-colors">Write a review</p>
                <p className="text-xs text-slate-400">Share your thoughts about this anime</p>
              </div>
            </button>
          )
        ) : (
          /* Not logged in */
          <div className="flex items-center gap-3 p-4 rounded-xl bg-slate-50 border border-slate-200">
            <div className="w-9 h-9 rounded-full bg-slate-200 flex items-center justify-center text-slate-400 text-lg shrink-0">✍</div>
            <p className="text-sm text-slate-600">
              <Link to="/login" className="text-brand-700 font-semibold hover:underline">Sign in</Link> to write a review and share your opinion.
            </p>
          </div>
        )}

        {/* Divider if there are other reviews */}
        {otherReviews.length > 0 && myReview && (
          <hr className="border-slate-200" />
        )}

        {/* Other reviews feed */}
        {otherReviews.length > 0 && (
          <ul className="flex flex-col gap-3" data-testid="anime-detail-reviews">
            {otherReviews.map((review) => (
              <li key={review.id}>
                <ReviewCard review={review} />
              </li>
            ))}
          </ul>
        )}

        {/* Empty state */}
        {totalReviews === 0 && !showForm && (
          <p className="text-sm text-slate-400 text-center py-4" data-testid="anime-detail-no-reviews">
            No reviews yet — be the first to share your thoughts!
          </p>
        )}
      </section>
    </article>
    </>
  )
}
