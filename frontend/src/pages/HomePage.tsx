import { useState, useEffect, useMemo } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import { listAnimes, searchAnimes } from '../lib/api/animes'
import { AnimeCard } from '../components/AnimeCard'
import { Button } from '../components/Button'
import { Input } from '../components/Input'

const PAGE_SIZE = 12

function useDebouncedValue<T>(value: T, delay = 300): T {
  const [debounced, setDebounced] = useState(value)
  useEffect(() => {
    const id = window.setTimeout(() => setDebounced(value), delay)
    return () => window.clearTimeout(id)
  }, [value, delay])
  return debounced
}

export function HomePage() {
  const [search, setSearch] = useState('')
  const [page, setPage] = useState(1)
  const debouncedSearch = useDebouncedValue(search.trim(), 350)

  // Reset to page 1 whenever the active search term changes.
  useEffect(() => {
    setPage(1)
  }, [debouncedSearch])

  const queryKey = useMemo(
    () => ['animes', { search: debouncedSearch, page }] as const,
    [debouncedSearch, page],
  )

  const { data, isLoading, isError, error, isFetching } = useQuery({
    queryKey,
    queryFn: () =>
      debouncedSearch
        ? searchAnimes(debouncedSearch, { page, limit: PAGE_SIZE })
        : listAnimes({ page, limit: PAGE_SIZE, sort_by: 'created_at', order: 'desc' }),
    placeholderData: keepPreviousData,
  })

  const totalPages = data ? Math.max(1, Math.ceil(data.total / PAGE_SIZE)) : 1

  return (
    <div className="flex flex-col gap-6">
      <header className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold">Catalog</h1>
        <p className="text-slate-600">Browse and search anime.</p>
      </header>

      <Input
        label="Search by title"
        placeholder="e.g. naruto"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
        data-testid="catalog-search"
      />

      {isLoading && <p className="text-slate-500">Loading...</p>}

      {isError && (
        <p className="text-red-600" role="alert">
          Could not load the catalog. {(error as Error).message}
        </p>
      )}

      {data && data.data.length === 0 && !isLoading && (
        <p className="text-slate-500" data-testid="catalog-empty">
          {debouncedSearch
            ? `No anime matched "${debouncedSearch}".`
            : 'No anime in the catalog yet.'}
        </p>
      )}

      {data && data.data.length > 0 && (
        <ul
          className="grid gap-3 grid-cols-1 sm:grid-cols-2 lg:grid-cols-3"
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
            Previous
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
            Next
          </Button>
        </div>
      )}
    </div>
  )
}
