import { useState, useEffect, useMemo, useRef } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import { listAnimes, searchAnimes, type ListAnimesParams } from '../lib/api/animes'
import { listGenres } from '../lib/api/genres'
import { listFavorites } from '../lib/api/favorites'
import { getWatchlist } from '../lib/api/watchlist'
import { AnimeCard } from '../components/AnimeCard'
import { Button } from '../components/Button'
import { Input } from '../components/Input'
import { ToastContainer, useToast } from '../components/Toast'
import { useAuth } from '../lib/auth/AuthProvider'

const PAGE_SIZE = 12

const SORT_OPTIONS: { label: string; value: ListAnimesParams['sort_by']; order: 'asc' | 'desc' }[] = [
  { label: 'Newest first',  value: 'created_at', order: 'desc' },
  { label: 'Highest rated', value: 'rating',     order: 'desc' },
  { label: 'Most episodes', value: 'episodes',   order: 'desc' },
  { label: 'Title A → Z',  value: 'title',      order: 'asc'  },
  { label: 'Title Z → A',  value: 'title',      order: 'desc' },
  { label: 'Oldest first',  value: 'created_at', order: 'asc'  },
]

const ALL_STATUSES = ['Airing', 'Completed', 'Upcoming'] as const

function useDebouncedValue<T>(value: T, delay = 350): T {
  const [debounced, setDebounced] = useState(value)
  useEffect(() => {
    const id = window.setTimeout(() => setDebounced(value), delay)
    return () => window.clearTimeout(id)
  }, [value, delay])
  return debounced
}

function toggle<T>(arr: T[], item: T): T[] {
  return arr.includes(item) ? arr.filter((x) => x !== item) : [...arr, item]
}

// ── Genre dropdown with checkboxes ──────────────────────────────────────────
interface GenreDropdownProps {
  genres: { id: number; name: string }[]
  selected: string[]
  onChange: (selected: string[]) => void
}

function GenreDropdown({ genres, selected, onChange }: GenreDropdownProps) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  // Close on outside click
  useEffect(() => {
    function onDown(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false)
    }
    document.addEventListener('mousedown', onDown)
    return () => document.removeEventListener('mousedown', onDown)
  }, [])

  return (
    <div ref={ref} className="relative">
      <button
        type="button"
        onClick={() => setOpen((o) => !o)}
        className={
          'flex items-center gap-1.5 h-9 px-3 rounded-lg border text-sm font-medium transition-colors ' +
          (selected.length > 0
            ? 'border-brand-500 bg-brand-50 text-brand-700'
            : 'border-slate-300 bg-white text-slate-700 hover:border-slate-400')
        }
        aria-haspopup="listbox"
        aria-expanded={open}
        data-testid="catalog-genre-btn"
      >
        Genres
        {selected.length > 0 && (
          <span className="inline-flex items-center justify-center w-5 h-5 rounded-full bg-brand-600 text-white text-[11px] font-bold">
            {selected.length}
          </span>
        )}
        <span className="text-slate-400">{open ? '▲' : '▼'}</span>
      </button>

      {open && (
        <div
          className="absolute z-20 mt-1 w-52 max-h-72 overflow-y-auto rounded-xl border border-slate-200 bg-white shadow-lg py-1"
          role="listbox"
          aria-multiselectable="true"
        >
          {selected.length > 0 && (
            <button
              type="button"
              onClick={() => onChange([])}
              className="w-full text-left px-3 py-1.5 text-xs text-brand-600 font-medium hover:bg-slate-50"
            >
              ✕ Clear genres
            </button>
          )}
          {genres.map((g) => {
            const checked = selected.includes(g.name)
            return (
              <label
                key={g.id}
                className="flex items-center gap-2 px-3 py-1.5 text-sm text-slate-700 hover:bg-slate-50 cursor-pointer"
              >
                <input
                  type="checkbox"
                  className="accent-brand-600 w-3.5 h-3.5"
                  checked={checked}
                  onChange={() => onChange(toggle(selected, g.name))}
                  aria-label={g.name}
                />
                {g.name}
              </label>
            )
          })}
        </div>
      )}
    </div>
  )
}

// ── Main page ────────────────────────────────────────────────────────────────
const selectClass =
  'h-9 rounded-lg border border-slate-300 bg-white px-3 text-sm text-slate-700 ' +
  'focus:outline-none focus:ring-2 focus:ring-brand-500 transition-colors'

export function HomePage() {
  const [search, setSearch]       = useState('')
  const [sortIdx, setSortIdx]     = useState(0)
  const [statuses, setStatuses]   = useState<string[]>([])
  const [genres, setGenres]       = useState<string[]>([])
  const [page, setPage]           = useState(1)
  const { isAuthenticated }       = useAuth()
  const { toasts, push: pushToast, remove: removeToast } = useToast()

  const debouncedSearch = useDebouncedValue(search.trim())
  const sort = SORT_OPTIONS[sortIdx]

  useEffect(() => { setPage(1) }, [debouncedSearch, sortIdx, statuses, genres])

  const queryKey = useMemo(
    () => ['animes', { search: debouncedSearch, sortIdx, statuses, genres, page }] as const,
    [debouncedSearch, sortIdx, statuses, genres, page],
  )

  const { data, isLoading, isError, error, isFetching } = useQuery({
    queryKey,
    queryFn: () =>
      debouncedSearch
        ? searchAnimes(debouncedSearch, {
            page, limit: PAGE_SIZE,
            sort_by: sort.value, order: sort.order,
            genres: genres.length ? genres : undefined,
            statuses: statuses.length ? statuses : undefined,
          })
        : listAnimes({
            page, limit: PAGE_SIZE,
            sort_by: sort.value, order: sort.order,
            statuses: statuses.length ? statuses : undefined,
            genres: genres.length ? genres : undefined,
          }),
    placeholderData: keepPreviousData,
  })

  const { data: genreList } = useQuery({
    queryKey: ['genres'],
    queryFn: listGenres,
    staleTime: Infinity,
  })

  const { data: favorites } = useQuery({
    queryKey: ['favorites'],
    queryFn: listFavorites,
    enabled: isAuthenticated,
    staleTime: 30_000,
  })

  const { data: watchlist } = useQuery({
    queryKey: ['watchlist'],
    queryFn: getWatchlist,
    enabled: isAuthenticated,
    staleTime: 30_000,
  })

  const favoritedIds = useMemo(
    () => new Set((favorites ?? []).map((f) => f.anime_id)),
    [favorites],
  )

  const watchlistMap = useMemo(() => {
    const m = new Map<number, string>()
    for (const e of watchlist?.entries ?? []) m.set(e.anime_id, e.status)
    return m
  }, [watchlist])

  const totalPages    = data ? Math.max(1, Math.ceil(data.total / PAGE_SIZE)) : 1
  const activeFilters = statuses.length + genres.length

  const clearFilters = () => {
    setSearch('')
    setSortIdx(0)
    setStatuses([])
    setGenres([])
  }

  return (
    <>
    <ToastContainer toasts={toasts} onDone={removeToast} />
    <div className="flex flex-col gap-6">
      <header className="flex flex-col gap-1">
        <h1 className="text-3xl font-bold">Catalog</h1>
        <p className="text-slate-500 text-sm">
          {data ? `${data.total.toLocaleString()} titles` : 'Browse and search anime.'}
        </p>
      </header>

      {/* Search */}
      <Input
        label="Search by title"
        placeholder="e.g. Naruto, Attack on Titan…"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        data-testid="catalog-search"
      />

      {/* Filter bar */}
      <div className="flex flex-wrap items-center gap-2">
        {/* Sort */}
        <select
          value={sortIdx}
          onChange={(e) => setSortIdx(Number(e.target.value))}
          className={selectClass}
          aria-label="Sort by"
          data-testid="catalog-sort"
        >
          {SORT_OPTIONS.map((opt, i) => (
            <option key={i} value={i}>{opt.label}</option>
          ))}
        </select>

        {/* Status chips — multi-select, OR logic */}
        <div className="flex items-center gap-1.5" role="group" aria-label="Filter by status">
          {ALL_STATUSES.map((s) => {
            const active = statuses.includes(s)
            return (
              <button
                key={s}
                type="button"
                onClick={() => setStatuses((prev) => toggle(prev, s))}
                className={
                  'h-9 px-3 rounded-lg border text-sm font-medium transition-colors ' +
                  (active
                    ? 'border-brand-500 bg-brand-600 text-white'
                    : 'border-slate-300 bg-white text-slate-600 hover:border-slate-400 hover:bg-slate-50')
                }
                aria-pressed={active}
                data-testid={`catalog-status-${s.toLowerCase()}`}
              >
                {s}
              </button>
            )
          })}
        </div>

        {/* Genre multi-select dropdown */}
        {genreList && (
          <GenreDropdown
            genres={genreList}
            selected={genres}
            onChange={setGenres}
          />
        )}

        {/* Clear all */}
        {(activeFilters > 0 || search) && (
          <button
            type="button"
            onClick={clearFilters}
            className="flex items-center gap-1 h-9 px-3 rounded-lg bg-slate-100 text-sm font-medium text-slate-600 hover:bg-slate-200 transition-colors"
            data-testid="catalog-clear"
          >
            ✕ Clear
            {activeFilters > 0 && (
              <span className="inline-flex items-center justify-center w-5 h-5 rounded-full bg-slate-400 text-white text-[11px] font-bold">
                {activeFilters}
              </span>
            )}
          </button>
        )}
      </div>

      {/* Active filter chips summary */}
      {(statuses.length > 0 || genres.length > 0) && (
        <div className="flex flex-wrap gap-1.5 -mt-2">
          {statuses.map((s) => (
            <span
              key={s}
              className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-brand-100 text-brand-700 text-xs font-medium"
            >
              {s}
              <button type="button" onClick={() => setStatuses((p) => p.filter((x) => x !== s))} aria-label={`Remove ${s}`}>✕</button>
            </span>
          ))}
          {genres.map((g) => (
            <span
              key={g}
              className="inline-flex items-center gap-1 px-2 py-0.5 rounded-full bg-slate-100 text-slate-700 text-xs font-medium"
            >
              {g}
              <button type="button" onClick={() => setGenres((p) => p.filter((x) => x !== g))} aria-label={`Remove ${g}`}>✕</button>
            </span>
          ))}
        </div>
      )}

      {isLoading && <p className="text-slate-500">Loading…</p>}
      {isError && (
        <p className="text-red-600" role="alert">
          Could not load the catalog. {(error as Error).message}
        </p>
      )}
      {data && data.data.length === 0 && !isLoading && (
        <p className="text-slate-500" data-testid="catalog-empty">
          {debouncedSearch
            ? `No anime matched "${debouncedSearch}".`
            : 'No anime found for the selected filters.'}
        </p>
      )}

      {data && data.data.length > 0 && (
        <ul
          className={`grid gap-3 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 transition-opacity ${isFetching ? 'opacity-60' : 'opacity-100'}`}
          data-testid="catalog-list"
          aria-busy={isFetching}
        >
          {data.data.map((anime) => (
            <li key={anime.id}>
              <AnimeCard
                anime={anime}
                isAuthenticated={isAuthenticated}
                isFavorited={favoritedIds.has(anime.id)}
                watchlistStatus={(watchlistMap.get(anime.id) ?? null) as import('../lib/api/watchlist').WatchlistStatus | null}
                onAction={(msg) => pushToast(msg)}
              />
            </li>
          ))}
        </ul>
      )}

      {data && data.total > PAGE_SIZE && (
        <div className="flex items-center justify-between gap-3 pt-2">
          <Button variant="secondary" disabled={page <= 1 || isFetching} onClick={() => setPage((p) => Math.max(1, p - 1))} data-testid="catalog-prev">
            ← Previous
          </Button>
          <span className="text-sm text-slate-600" data-testid="catalog-page-label">
            Page {page} of {totalPages}
          </span>
          <Button variant="secondary" disabled={page >= totalPages || isFetching} onClick={() => setPage((p) => p + 1)} data-testid="catalog-next">
            Next →
          </Button>
        </div>
      )}
    </div>
    </>
  )
}
