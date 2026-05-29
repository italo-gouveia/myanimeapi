import { Link } from 'react-router-dom'
import type { Anime } from '../lib/api/types'

export function AnimeCard({ anime }: { anime: Anime }) {
  return (
    <Link
      to={`/animes/${anime.id}`}
      className="block bg-white border border-slate-200 rounded-lg p-4 hover:shadow-md hover:border-brand-300 transition"
      data-testid={`anime-card-${anime.id}`}
    >
      <div className="flex items-start justify-between gap-2">
        <h3 className="font-semibold text-slate-900 line-clamp-2">{anime.title}</h3>
        <span
          className="shrink-0 inline-block bg-brand-100 text-brand-700 text-xs font-semibold px-2 py-0.5 rounded"
          aria-label={`Rating ${anime.rating}`}
        >
          ★ {anime.rating.toFixed(1)}
        </span>
      </div>
      <p className="text-xs text-slate-500 mt-1">
        {anime.status} · {anime.episodes} ep
      </p>
      {anime.description && (
        <p className="text-sm text-slate-600 mt-2 line-clamp-3">{anime.description}</p>
      )}
    </Link>
  )
}
