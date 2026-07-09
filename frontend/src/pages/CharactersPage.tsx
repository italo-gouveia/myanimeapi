import { useState } from 'react'
import { useQuery, keepPreviousData } from '@tanstack/react-query'
import { listCharacters, searchCharacters } from '../lib/api/characters'
import { CharacterCard } from '../components/CharacterCard'

const PAGE_SIZE = 20

export function CharactersPage() {
  const [page, setPage]     = useState(1)
  const [search, setSearch] = useState('')
  const [query, setQuery]   = useState('')

  const queryKey = query
    ? ['characters', 'search', query, page]
    : ['characters', page]

  const { data, isLoading, isError } = useQuery({
    queryKey,
    queryFn: () =>
      query
        ? searchCharacters(query, { page, limit: PAGE_SIZE })
        : listCharacters({ page, limit: PAGE_SIZE }),
    placeholderData: keepPreviousData,
  })

  const totalPages = data ? Math.ceil(data.total / PAGE_SIZE) : 1

  function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    setQuery(search.trim())
    setPage(1)
  }

  function handleClear() {
    setSearch('')
    setQuery('')
    setPage(1)
  }

  return (
    <div className="flex flex-col gap-6">
      <div className="flex flex-col gap-1">
        <h1 className="text-2xl font-bold">Characters</h1>
        <p className="text-sm text-slate-500">Browse all anime characters</p>
      </div>

      {/* Search */}
      <form onSubmit={handleSearch} className="flex gap-2">
        <input
          type="search"
          placeholder="Search by name…"
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          aria-label="Search characters by name"
          className="flex-1 px-3 py-2 rounded-lg border border-slate-300 bg-white text-slate-900 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500"
        />
        <button
          type="submit"
          className="px-4 py-2 rounded-lg bg-brand-600 text-white text-sm font-medium hover:bg-brand-700 transition-colors"
        >
          Search
        </button>
        {query && (
          <button
            type="button"
            onClick={handleClear}
            className="px-4 py-2 rounded-lg border border-slate-300 text-sm text-slate-600 hover:bg-slate-50 transition-colors"
          >
            Clear
          </button>
        )}
      </form>

      {/* Results info */}
      {data && (
        <p className="text-sm text-slate-500">
          {data.total === 0
            ? 'No characters found.'
            : `${data.total} character${data.total !== 1 ? 's' : ''}${query ? ` matching "${query}"` : ''}`}
        </p>
      )}

      {/* Loading / Error / Grid */}
      {isLoading && <p className="text-slate-400">Loading…</p>}

      {isError && (
        <p className="text-red-600" role="alert">Could not load characters. Please try again.</p>
      )}

      {data && data.data.length > 0 && (
        <ul
          className="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 gap-4"
          data-testid="characters-grid"
        >
          {data.data.map((char) => (
            <li key={char.id}>
              <CharacterCard character={char} />
            </li>
          ))}
        </ul>
      )}

      {/* Pagination */}
      {totalPages > 1 && (
        <nav className="flex items-center justify-center gap-2" aria-label="Pagination">
          <button
            type="button"
            onClick={() => setPage((p) => Math.max(1, p - 1))}
            disabled={page <= 1}
            className="px-3 py-1.5 rounded-lg border border-slate-300 text-sm text-slate-600 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            ← Prev
          </button>
          <span className="text-sm text-slate-500 tabular-nums">
            {page} / {totalPages}
          </span>
          <button
            type="button"
            onClick={() => setPage((p) => Math.min(totalPages, p + 1))}
            disabled={page >= totalPages}
            className="px-3 py-1.5 rounded-lg border border-slate-300 text-sm text-slate-600 hover:bg-slate-50 disabled:opacity-40 disabled:cursor-not-allowed transition-colors"
          >
            Next →
          </button>
        </nav>
      )}
    </div>
  )
}
