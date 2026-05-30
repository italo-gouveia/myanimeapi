import { Link } from 'react-router-dom'
import type { Anime } from '../lib/api/types'

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

export function AnimeCard({ anime }: { anime: Anime }) {
  return (
    <Link
      to={`/animes/${anime.id}`}
      className="group flex flex-col bg-white border border-slate-200 rounded-lg overflow-hidden hover:shadow-md hover:border-brand-300 transition"
      data-testid={`anime-card-${anime.id}`}
    >
      {/* Cover image */}
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
        {/* Rating badge overlay */}
        <span
          className="absolute top-2 right-2 inline-block bg-black/60 text-yellow-400 text-xs font-bold px-2 py-0.5 rounded backdrop-blur-sm"
          aria-label={`Rating ${anime.rating}`}
        >
          ★ {anime.rating.toFixed(1)}
        </span>
      </div>

      {/* Info */}
      <div className="p-3 flex flex-col gap-1">
        <h3 className="font-semibold text-slate-900 line-clamp-2 text-sm leading-snug">
          {anime.title}
        </h3>
        <p className="text-xs text-slate-500">
          {anime.status} · {anime.episodes} ep
        </p>
      </div>
    </Link>
  )
}
