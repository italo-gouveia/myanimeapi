import { useState } from 'react'
import { useMutation } from '@tanstack/react-query'
import { triggerETLSync, type ETLSyncResult } from '../../lib/api/admin'
import { Button } from '../../components/Button'

export function AdminETLPage() {
  const [pages, setPages] = useState(5)
  const [result, setResult] = useState<ETLSyncResult | null>(null)

  const { mutate, isPending, isError, error } = useMutation({
    mutationFn: () => triggerETLSync(pages),
    onSuccess: (data) => setResult(data),
  })

  return (
    <div className="flex flex-col gap-6">
      <h1 className="text-2xl font-bold">ETL Sync</h1>
      <p className="text-slate-500 text-sm max-w-lg">
        Pull anime data from the{' '}
        <a
          href="https://jikan.moe"
          target="_blank"
          rel="noreferrer"
          className="text-brand-700 underline"
        >
          Jikan API
        </a>{' '}
        (MyAnimeList public data). Each page fetches 25 titles. Jikan enforces a ~3 req/s limit, so
        syncing many pages takes a moment.
      </p>

      <div className="bg-white border border-slate-200 rounded-xl p-6 flex flex-col gap-5 max-w-sm">
        <div className="flex flex-col gap-1">
          <label htmlFor="pages-input" className="text-sm font-medium text-slate-700">
            Pages to sync{' '}
            <span className="ml-1 text-xs font-normal text-slate-400">(1 page = 25 anime)</span>
          </label>
          <input
            id="pages-input"
            type="number"
            min={1}
            max={20}
            value={pages}
            onChange={(e) => setPages(Math.min(20, Math.max(1, Number(e.target.value))))}
            className="border border-slate-300 rounded-lg px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-brand-500 w-24"
          />
          <p className="text-xs text-slate-400">≈ {pages * 25} anime · max 20 pages</p>
        </div>

        <Button onClick={() => mutate()} disabled={isPending} className="w-full">
          {isPending ? 'Syncing…' : '🔄 Start sync'}
        </Button>

        {isError && (
          <p className="text-sm text-red-600" role="alert">
            Sync failed: {(error as Error).message}
          </p>
        )}
      </div>

      {result && (
        <div
          className="bg-green-50 border border-green-200 rounded-xl p-6 max-w-sm"
          data-testid="etl-result"
        >
          <p className="text-sm font-semibold text-green-800 mb-3">Sync complete ✓</p>
          <dl className="grid grid-cols-2 gap-x-4 gap-y-2 text-sm">
            <dt className="text-slate-500">Pages processed</dt>
            <dd className="font-medium text-slate-800">{result.pages}</dd>
            <dt className="text-slate-500">Imported</dt>
            <dd className="font-medium text-green-700">{result.imported}</dd>
            <dt className="text-slate-500">Updated</dt>
            <dd className="font-medium text-blue-700">{result.updated}</dd>
            <dt className="text-slate-500">Skipped</dt>
            <dd className="font-medium text-slate-600">{result.skipped}</dd>
            {result.errors > 0 && (
              <>
                <dt className="text-slate-500">Errors</dt>
                <dd className="font-medium text-red-600">{result.errors}</dd>
              </>
            )}
          </dl>
        </div>
      )}
    </div>
  )
}
