import { useState, useEffect, useMemo } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import { listAnimes, searchAnimes, type ListAnimesParams } from '../lib/api/animes'
import { listGenres } from '../lib/api/genres'
import { AnimeCard } from '../components/AnimeCard'
import { Button } from '../components/Button'
import { Input } from '../components/Input'

const PAGE_SIZE = 12

const SORT_OPTIONS: { label: string; value: ListAnimesParams['sort_by']; order: 'asc' | 'desc' }[] = [
  { label: 'Newest first',   value: 'created_at',  order: 'desc' },
  { label: 'Highest rated',  value: 'rating',      order: 'desc' },
  { label: 'Most episodes',  value: 'episodes',    order: 'desc' },
  { label: 'Title A → Z',   value: 'title',       order: 'asc'  },
  { label: 'Title Z → A',   value: 'title',       order: 'desc' },
  { label: 'Oldest first',   value: 'created_at',  order: 'asc'  },
]

const STATUS_OPTIONS = [
  { label: 'All statuses', value: '' },
  { label: 'Airing',       value: 'Airing' },
  { label: 'Completed',    value: 'Completed' },
  { label: 'Upcoming',     value: 'Upcoming' },
]

function useDebouncedValue<T>(value: T, delay = 350): T {
  const [debounced, setDebounced] = useState(value)
  useEffect(() => {
    const id = window.setTimeout(() => setDebounced(value), delay)
    return () => window.clearTimeout(id)
  }, [value, delay])
  return debounced
}

const selectClass =
  'h-9 rounded-lg border border-slate-300 bg-white px-3 text-sm text-slate-700 ' +
  'focus:outline-none focus:ring-2 focus:ring-brand-500 transition-colors'

export function HomePage() {
  const [search, setSearch]       = useState('')
  const [sortIdx, setSortIdx]     = useState(0)
  const [status, setStatus]       = useState('')
  const [genre, setGenre]         = useState('')
  const [page, setPage]           = useState(1)

  const debouncedSearch = useDebouncedValue(search.trim())
  const sort = SORT_OPTIONS[sortIdx]

  // Reset to page 1 when any filter changes.
  useEffect(() => { setPage(1) }, [debouncedSearch, sortIdx, status, genre])

  const queryKey = useMemo(
    () => ['animes', { search: debouncedSearch, sortIdx, status, genre, page }] as const,
    [debouncedSearch, sortIdx, status, genre, page],
  )

  const { data, isLoading, isError, error, isFetching } = useQuery({
    queryKey,
    queryFn: () =>
      debouncedSearch
        ? searchAnimes(debouncedSearch, { page, limit: PAGE_SIZE, sort_by: sort.value, order: sort.order })
        : listAnimes({ page, limit: PAGE_SIZE, sort_by: sort.value, order: sort.order, status: status || undefined, genre: genre || undefined }),
    placeholderData: keepPreviousData,
  })

  // Genres for the dropdown — loaded once, cached.
  const { data: genres } = useQuery({
    queryKey: ['genres'],
    queryFn: listGenres,
    staleTime: Infinity,
  })

  const totalPages  = data ? Math.max(1, Math.ceil(data.total / PAGE_SIZE)) : 1
  const activeFilters = [status, genre].filter(Boolean).length

  const clearFilters = () => {
    setSearch('')
    setSortIdx(0)
    setStatus('')
    setGenre('')
  }

  return (
    <div className="flex flex-col gap-6">
      <header className="flex flex-col gap-1">
        <h1 className="text-3xl font-bold">Catalog</h1>
        <p className="text-slate-500 text-sm">
          {data ? `${data.total.toLocaleString()} titles` : 'Browse and search anime.'}
        </p>
      </header>

      {/* Search + filters row */}
      <div className="flex flex-col gap-3">
        <Input
          label="Search by title"
          placeholder="e.g. Naruto, Attack on Titan…"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          data-testid="catalog-search"
        />

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

          {/* Status */}
          <select
            value={status}
            onChange={(e) => setStatus(e.target.value)}
            className={selectClass}
            aria-label="Filter by status"
            data-testid="catalog-status"
          >
            {STATUS_OPTIONS.map((opt) => (
              <option key={opt.value} value={opt.value}>{opt.label}</option>
            ))}
          </select>

          {/* Genre */}
          <select
            value={genre}
            onChange={(e) => setGenre(e.target.value)}
            className={selectClass}
            aria-label="Filter by genre"
            data-testid="catalog-genre"
          >
            <option value="">All genres</option>
            {genres?.map((g) => (
              <option key={g.id} value={g.name}>{g.name}</option>
            ))}
          </select>

          {/* Clear filters badge */}
          {(activeFilters > 0 || search) && (
            <button
              type="button"
              onClick={clearFilters}
              className="flex items-center gap-1 px-2.5 py-1.5 rounded-lg bg-slate-100 text-xs font-medium text-slate-600 hover:bg-slate-200 transition-colors"
              data-testid="catalog-clear"
            >
              ✕ Clear
              {activeFilters > 0 && (
                <span className="inline-flex items-center justify-center w-4 h-4 rounded-full bg-brand-600 text-white text-[10px]">
                  {activeFilters}
                </span>
              )}
            </button>
          )}
        </div>
      </div>

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
              <AnimeCard anime={anime} />
            </li>
          ))}
        </ul>
      )}

      {data && data.total > PAGE_SIZE && (
        <div className="flex items-center justify-between gap-3 pt-2">
          <Button
            variant="secondary"
            disabled={page <= 1 || isFetching}
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            data-testid="catalog-prev"
          >
            ← Previous
          </Button>
          <span className="text-sm text-slate-600" data-testid="catalog-page-label">
            Page {page} of {totalPages}
          </span>
          <Button
            variant="secondary"
            disabled={page >= totalPages || isFetching}
            onClick={() => setPage((p) => p + 1)}
            data-testid="catalog-next"
          >
            Next →
          </Button>
        </div>
      )}
    </div>
  )
}
