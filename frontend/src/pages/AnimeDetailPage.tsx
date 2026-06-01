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

function formatYear(date: string): string | null {
  const year = new Date(date).getUTCFullYear()
  return Number.isFinite(year) && year > 1900 ? String(year) : null
}

// ── Star rating input ────────────────────────────────────────────────────────
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
          className={`text-xl leading-none transition-colors ${
            n <= (hovered || value) ? 'text-yellow-400' : 'text-slate-300'
          }`}
          aria-label={`${n} star${n !== 1 ? 's' : ''}`}
          aria-pressed={value === n}
        >
          ★
        </button>
      ))}
      <span className="ml-2 text-sm text-slate-500 tabular-nums">
        {(hovered || value) > 0 ? `${hovered || value}/10` : 'Select'}
      </span>
    </div>
  )
}

// ── Review form ──────────────────────────────────────────────────────────────
interface ReviewFormProps {
  animeId: number
  initialContent?: string
  initialRating?: number
  reviewId?: number
  onDone: () => void
  onCancel?: () => void
}

function ReviewForm({
  animeId, initialContent = '', initialRating = 0,
  reviewId, onDone, onCancel,
}: ReviewFormProps) {
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

  const charLeft = 500 - content.length

  return (
    <form
      onSubmit={(e) => { e.preventDefault(); mutation.mutate() }}
      className="flex flex-col gap-3 bg-slate-50 border border-slate-200 rounded-xl p-4"
    >
      <p className="text-sm font-semibold text-slate-700">
        {isEdit ? 'Edit your review' : 'Write a review'}
      </p>

      <StarRating value={rating} onChange={setRating} />

      <div className="flex flex-col gap-1">
        <textarea
          value={content}
          onChange={(e) => setContent(e.target.value)}
          maxLength={500}
          rows={4}
          placeholder="Share your thoughts about this anime…"
          className="w-full rounded-lg border border-slate-300 px-3 py-2 text-sm resize-none focus:outline-none focus:ring-2 focus:ring-brand-500"
          required
        />
        <p className={`text-xs text-right ${charLeft < 50 ? 'text-orange-500' : 'text-slate-400'}`}>
          {charLeft} characters left
        </p>
      </div>

      {mutation.isError && (
        <p className="text-sm text-red-600" role="alert">
          Failed to submit. Please try again.
        </p>
      )}

      <div className="flex items-center gap-2">
        <Button
          type="submit"
          disabled={mutation.isPending || rating === 0 || content.trim().length === 0}
        >
          {mutation.isPending ? 'Saving…' : isEdit ? 'Save changes' : 'Submit review'}
        </Button>
        {onCancel && (
          <Button type="button" variant="secondary" onClick={onCancel}>
            Cancel
          </Button>
        )}
      </div>
    </form>
  )
}

// ── Main page ────────────────────────────────────────────────────────────────
export function AnimeDetailPage() {
  const { id }                      = useParams<{ id: string }>()
  const animeId                     = id ? Number.parseInt(id, 10) : NaN
  const { isAuthenticated, userId } = useAuth()
  const queryClient                 = useQueryClient()
  const [editingReviewId, setEditingReviewId] = useState<number | null>(null)

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

  const myWatchlistEntry = watchlist?.entries.find((e) => e.anime_id === animeId)

  const watchlistMutation = useMutation<void, Error, WatchlistStatus | null>({
    mutationFn: async (status) => {
      if (status) await upsertWatchlistEntry(animeId, status)
      else await removeWatchlistEntry(animeId)
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['watchlist'] }),
  })

  const isFavorited = favorites?.some((f) => f.anime_id === animeId) ?? false

  const addMutation = useMutation({
    mutationFn: () => addFavorite(animeId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['favorites'] }),
  })
  const removeMutation = useMutation({
    mutationFn: () => removeFavorite(animeId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['favorites'] }),
  })
  const deleteMutation = useMutation({
    mutationFn: (reviewId: number) => deleteReview(reviewId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['anime', animeId] }),
  })

  if (!Number.isFinite(animeId)) {
    return <p className="text-red-600" role="alert">Invalid anime id.</p>
  }
  if (isLoading) return <p className="text-slate-500">Loading…</p>
  if (isError || !anime) {
    return (
      <div role="alert" data-testid="anime-detail-error">
        <p className="text-red-600">
          Could not load this anime. {(error as Error | undefined)?.message ?? ''}
        </p>
        <Link to="/" className="text-brand-700 underline">Back to catalog</Link>
      </div>
    )
  }

  const startYear = formatYear(anime.start_date)
  const endYear   = formatYear(anime.end_date)
  const yearRange = startYear
    ? (endYear && endYear !== startYear ? `${startYear} – ${endYear}` : startYear)
    : null

  // Separate current user's review from others
  const myReview     = isAuthenticated && userId
    ? anime.reviews?.find((r) => (r.userId ?? r.user_id) === userId)
    : undefined
  const otherReviews = anime.reviews?.filter((r) => (r.userId ?? r.user_id) !== userId) ?? []

  return (
    <article className="flex flex-col gap-6" data-testid="anime-detail">
      <header className="flex flex-col gap-2">
        <Link to="/" className="text-sm text-brand-700 underline w-fit">← Catalog</Link>

        <div className="flex flex-col sm:flex-row gap-6">
          {anime.cover_url && (
            <div className="shrink-0 w-full sm:w-40 md:w-52">
              <img
                src={anime.cover_url}
                alt={`${anime.title} cover`}
                className="w-full rounded-lg shadow-md object-cover aspect-[3/4]"
              />
            </div>
          )}

          <div className="flex flex-col gap-3 min-w-0">
            <div className="flex items-start justify-between gap-4">
              <h1 className="text-3xl font-bold" data-testid="anime-detail-title">
                {anime.title}
              </h1>
              {isAuthenticated && (
                <div className="flex items-center gap-2 shrink-0">
                  {/* Watchlist dropdown */}
                  <select
                    value={myWatchlistEntry?.status ?? ''}
                    disabled={watchlistMutation.isPending}
                    onChange={(e) => {
                      const val = e.target.value
                      watchlistMutation.mutate(val ? (val as WatchlistStatus) : null)
                    }}
                    className="h-9 rounded-lg border border-slate-300 bg-white px-2 text-sm text-slate-700 focus:outline-none focus:ring-2 focus:ring-brand-500 disabled:opacity-50"
                    aria-label="Add to watchlist"
                    data-testid="watchlist-select"
                  >
                    <option value="">+ Watchlist</option>
                    {WATCHLIST_STATUSES.map((s) => (
                      <option key={s} value={s}>{WATCHLIST_STATUS_LABELS[s]}</option>
                    ))}
                  </select>

                  {/* Favorite */}
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
              <span className="inline-flex items-center gap-1 font-semibold text-brand-700">
                ★ {anime.rating.toFixed(1)}
              </span>
              <span>·</span><span>{anime.status}</span>
              <span>·</span><span>{anime.episodes} episodes</span>
              {yearRange && <><span>·</span><span>{yearRange}</span></>}
            </div>
          </div>
        </div>
      </header>

      {anime.description && (
        <section>
          <h2 className="text-lg font-semibold mb-1">Synopsis</h2>
          <p className="text-slate-700 whitespace-pre-line">{anime.description}</p>
        </section>
      )}

      {(anime.genres?.length || anime.tags?.length) ? (
        <section className="flex flex-wrap gap-2">
          {anime.genres?.map((g) => (
            <span key={`genre-${g.id}`} className="bg-brand-100 text-brand-700 text-xs font-semibold px-2 py-1 rounded" data-testid="anime-detail-genre">
              {g.name}
            </span>
          ))}
          {anime.tags?.map((t) => (
            <span key={`tag-${t.id}`} className="bg-slate-200 text-slate-700 text-xs px-2 py-1 rounded" data-testid="anime-detail-tag">
              #{t.name}
            </span>
          ))}
        </section>
      ) : null}

      {/* ── Reviews ── */}
      <section className="flex flex-col gap-4">
        <div className="flex items-center justify-between">
          <h2 className="text-lg font-semibold">
            Reviews
            {anime.reviews && anime.reviews.length > 0 && (
              <span className="ml-2 text-sm font-normal text-slate-500">
                ({anime.reviews.length})
              </span>
            )}
          </h2>
        </div>

        {/* Current user's review */}
        {isAuthenticated && (
          myReview ? (
            editingReviewId === myReview.id ? (
              <ReviewForm
                animeId={animeId}
                reviewId={myReview.id}
                initialContent={myReview.content}
                initialRating={myReview.rating}
                onDone={() => setEditingReviewId(null)}
                onCancel={() => setEditingReviewId(null)}
              />
            ) : (
              <div className="bg-brand-50 border border-brand-200 rounded-xl p-4" data-testid="my-review">
                <div className="flex items-center justify-between mb-2">
                  <div className="flex items-center gap-2">
                    <span className="text-xs font-semibold bg-brand-600 text-white px-2 py-0.5 rounded-full">
                      Your review
                    </span>
                    <span className="text-sm font-semibold text-brand-700">
                      {'★'.repeat(Math.round(myReview.rating / 2))} {myReview.rating}/10
                    </span>
                  </div>
                  <div className="flex items-center gap-3">
                    <button
                      type="button"
                      onClick={() => setEditingReviewId(myReview.id)}
                      className="text-xs text-slate-500 hover:text-brand-700 transition-colors"
                    >
                      Edit
                    </button>
                    <button
                      type="button"
                      onClick={() => { if (window.confirm('Delete your review?')) deleteMutation.mutate(myReview.id) }}
                      disabled={deleteMutation.isPending}
                      className="text-xs text-slate-500 hover:text-red-600 transition-colors disabled:opacity-50"
                    >
                      Delete
                    </button>
                  </div>
                </div>
                <p className="text-slate-700 whitespace-pre-line text-sm">{myReview.content}</p>
                <time className="text-xs text-slate-400 mt-2 block" dateTime={myReview.created_at}>
                  {new Date(myReview.created_at).toLocaleDateString()}
                </time>
              </div>
            )
          ) : (
            <ReviewForm animeId={animeId} onDone={() => {}} />
          )
        )}

        {!isAuthenticated && (
          <p className="text-sm text-slate-500">
            <Link to="/login" className="text-brand-700 underline">Sign in</Link> to write a review.
          </p>
        )}

        {/* Other users' reviews */}
        {otherReviews.length === 0 && !myReview && (
          <p className="text-slate-500" data-testid="anime-detail-no-reviews">
            No reviews yet. Be the first!
          </p>
        )}

        {otherReviews.length > 0 && (
          <ul className="flex flex-col gap-3" data-testid="anime-detail-reviews">
            {otherReviews.map((review) => (
              <li key={review.id} className="bg-white border border-slate-200 rounded-xl p-4">
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-semibold text-brand-700">
                    ★ {review.rating}/10
                  </span>
                  <time className="text-xs text-slate-500" dateTime={review.created_at}>
                    {new Date(review.created_at).toLocaleDateString()}
                  </time>
                </div>
                <p className="text-slate-700 whitespace-pre-line text-sm">{review.content}</p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </article>
  )
}
