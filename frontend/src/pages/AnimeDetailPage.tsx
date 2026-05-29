import { Link, useParams } from 'react-router-dom'
import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { getAnime } from '../lib/api/animes'
import { listFavorites, addFavorite, removeFavorite } from '../lib/api/favorites'
import { useAuth } from '../lib/auth/AuthProvider'
import { Button } from '../components/Button'

function formatYear(date: string): string | null {
  const year = new Date(date).getUTCFullYear()
  return Number.isFinite(year) && year > 1900 ? String(year) : null
}

export function AnimeDetailPage() {
  const { id } = useParams<{ id: string }>()
  const animeId = id ? Number.parseInt(id, 10) : NaN
  const { isAuthenticated } = useAuth()
  const queryClient = useQueryClient()

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

  const isFavorited = favorites?.some((f) => f.anime_id === animeId) ?? false

  const addMutation = useMutation({
    mutationFn: () => addFavorite(animeId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['favorites'] }),
  })

  const removeMutation = useMutation({
    mutationFn: () => removeFavorite(animeId),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ['favorites'] }),
  })

  if (!Number.isFinite(animeId)) {
    return (
      <p className="text-red-600" role="alert">
        Invalid anime id.
      </p>
    )
  }

  if (isLoading) {
    return <p className="text-slate-500">Loading...</p>
  }

  if (isError || !anime) {
    return (
      <div role="alert" data-testid="anime-detail-error">
        <p className="text-red-600">
          Could not load this anime. {(error as Error | undefined)?.message ?? ''}
        </p>
        <Link to="/" className="text-brand-700 underline">
          Back to catalog
        </Link>
      </div>
    )
  }

  const startYear = formatYear(anime.start_date)
  const endYear = formatYear(anime.end_date)
  const yearRange = startYear ? (endYear && endYear !== startYear ? `${startYear} – ${endYear}` : startYear) : null

  return (
    <article className="flex flex-col gap-6" data-testid="anime-detail">
      <header className="flex flex-col gap-2">
        <Link to="/" className="text-sm text-brand-700 underline w-fit">
          ← Catalog
        </Link>
        <div className="flex items-start justify-between gap-4">
          <h1 className="text-3xl font-bold" data-testid="anime-detail-title">
            {anime.title}
          </h1>
          {isAuthenticated && (
            <Button
              variant={isFavorited ? 'secondary' : 'primary'}
              className="shrink-0"
              disabled={addMutation.isPending || removeMutation.isPending}
              onClick={() =>
                isFavorited ? removeMutation.mutate() : addMutation.mutate()
              }
              data-testid="favorite-toggle"
              aria-pressed={isFavorited}
            >
              {isFavorited ? '★ Saved' : '☆ Save'}
            </Button>
          )}
        </div>
        <div className="flex items-center gap-3 text-sm text-slate-600">
          <span className="inline-flex items-center gap-1 font-semibold text-brand-700">
            ★ {anime.rating.toFixed(1)}
          </span>
          <span>·</span>
          <span>{anime.status}</span>
          <span>·</span>
          <span>{anime.episodes} episodes</span>
          {yearRange && (
            <>
              <span>·</span>
              <span>{yearRange}</span>
            </>
          )}
        </div>
      </header>

      {anime.description && (
        <section>
          <h2 className="text-lg font-semibold mb-1">Synopsis</h2>
          <p className="text-slate-700 whitespace-pre-line">{anime.description}</p>
        </section>
      )}

      {(anime.genres?.length || anime.tags?.length) && (
        <section className="flex flex-wrap gap-2">
          {anime.genres?.map((g) => (
            <span
              key={`genre-${g.id}`}
              className="bg-brand-100 text-brand-700 text-xs font-semibold px-2 py-1 rounded"
              data-testid="anime-detail-genre"
            >
              {g.name}
            </span>
          ))}
          {anime.tags?.map((t) => (
            <span
              key={`tag-${t.id}`}
              className="bg-slate-200 text-slate-700 text-xs px-2 py-1 rounded"
              data-testid="anime-detail-tag"
            >
              #{t.name}
            </span>
          ))}
        </section>
      )}

      <section>
        <h2 className="text-lg font-semibold mb-3">Reviews</h2>
        {!anime.reviews || anime.reviews.length === 0 ? (
          <p className="text-slate-500" data-testid="anime-detail-no-reviews">
            No reviews yet.
          </p>
        ) : (
          <ul className="flex flex-col gap-3" data-testid="anime-detail-reviews">
            {anime.reviews.map((review) => (
              <li
                key={review.id}
                className="bg-white border border-slate-200 rounded-lg p-4"
              >
                <div className="flex items-center justify-between mb-2">
                  <span className="text-sm font-semibold text-brand-700">
                    ★ {review.rating}
                  </span>
                  <time className="text-xs text-slate-500" dateTime={review.created_at}>
                    {new Date(review.created_at).toLocaleDateString()}
                  </time>
                </div>
                <p className="text-slate-700 whitespace-pre-line">{review.content}</p>
              </li>
            ))}
          </ul>
        )}
      </section>
    </article>
  )
}
