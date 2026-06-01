import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { Link } from 'react-router-dom'
import { listFavorites, removeFavorite } from '../lib/api/favorites'
import { Button } from '../components/Button'
import { ToastContainer, useToast } from '../components/Toast'

export function FavoritesPage() {
  const queryClient = useQueryClient()
  const { toasts, push: pushToast, remove: removeToast } = useToast()

  const { data, isLoading, isError, error } = useQuery({
    queryKey: ['favorites'],
    queryFn: listFavorites,
  })

  const removeMutation = useMutation({
    mutationFn: (animeId: number) => removeFavorite(animeId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['favorites'] })
      pushToast('Removed from favorites', 'info')
    },
  })

  return (
    <>
    <ToastContainer toasts={toasts} onDone={removeToast} />
    <div className="flex flex-col gap-6">
      <header className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold">Favorites</h1>
        <p className="text-slate-600">Anime you've saved to watch or revisit.</p>
      </header>

      {isLoading && <p className="text-slate-500">Loading favorites...</p>}

      {isError && (
        <p className="text-red-600" role="alert">
          Could not load favorites. {(error as Error).message}
        </p>
      )}

      {removeMutation.isError && (
        <p className="text-red-600" role="alert">
          Failed to remove. {(removeMutation.error as Error).message}
        </p>
      )}

      {data && data.length === 0 && (
        <div className="flex flex-col items-center gap-4 py-16 text-center">
          <span className="text-5xl">🎌</span>
          <p className="text-slate-500" data-testid="favorites-empty">
            No favorites yet.{' '}
            <Link to="/" className="text-brand-700 underline">
              Browse the catalog
            </Link>{' '}
            and add some!
          </p>
        </div>
      )}

      {data && data.length > 0 && (
        <ul
          className="grid gap-4 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3"
          data-testid="favorites-list"
        >
          {data.map((fav) => {
            const anime = fav.anime
            if (!anime) return null
            return (
              <li key={fav.id} className="relative group">
                <Link
                  to={`/animes/${anime.id}`}
                  className="block bg-white border border-slate-200 rounded-lg p-4 hover:shadow-md hover:border-brand-300 transition"
                  data-testid={`favorite-card-${anime.id}`}
                >
                  <div className="flex items-start justify-between gap-2">
                    <h3 className="font-semibold text-slate-900 line-clamp-2 pr-6">
                      {anime.title}
                    </h3>
                    <span className="shrink-0 inline-block bg-brand-100 text-brand-700 text-xs font-semibold px-2 py-0.5 rounded">
                      ★ {anime.rating.toFixed(1)}
                    </span>
                  </div>
                  <p className="text-xs text-slate-500 mt-1">
                    {anime.status} · {anime.episodes} ep
                  </p>
                  {anime.description && (
                    <p className="text-sm text-slate-600 mt-2 line-clamp-3">
                      {anime.description}
                    </p>
                  )}
                </Link>
                <Button
                  variant="danger"
                  className="absolute top-3 right-3 !px-2 !py-1 text-xs opacity-0 group-hover:opacity-100 focus:opacity-100 transition-opacity"
                  aria-label={`Remove ${anime.title} from favorites`}
                  disabled={removeMutation.isPending}
                  onClick={(e) => {
                    e.preventDefault()
                    removeMutation.mutate(anime.id)
                  }}
                  data-testid={`remove-favorite-${anime.id}`}
                >
                  ✕ Remove
                </Button>
              </li>
            )
          })}
        </ul>
      )}
    </div>
    </>
  )
}
