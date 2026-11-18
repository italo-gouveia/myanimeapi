interface TableSkeletonProps {
  rows?: number
  cols: number
}

export function TableSkeleton({ rows = 5, cols }: Readonly<TableSkeletonProps>) {
  const rowKeys = Array.from({ length: rows }, (_, i) => `skel-row-${i}`)
  const colKeys = Array.from({ length: cols }, (_, j) => `skel-col-${j}`)
  return (
    <>
      {rowKeys.map((rowKey) => (
        <tr key={rowKey}>
          {colKeys.map((colKey) => (
            <td key={colKey} className="px-4 py-3">
              <div className="h-4 bg-slate-100 rounded animate-pulse" />
            </td>
          ))}
        </tr>
      ))}
    </>
  )
}

interface PaginationProps {
  page: number
  totalPages: number
  onPrev: () => void
  onNext: () => void
}

export function Pagination({ page, totalPages, onPrev, onNext }: Readonly<PaginationProps>) {
  if (totalPages <= 1) return null
  return (
    <div className="flex items-center justify-between">
      <button
        type="button"
        disabled={page === 1}
        onClick={onPrev}
        className="px-3 py-1.5 rounded border border-slate-200 text-sm disabled:opacity-40 hover:bg-slate-50"
      >
        ← Previous
      </button>
      <span className="text-sm text-slate-500">
        Page {page} of {totalPages}
      </span>
      <button
        type="button"
        disabled={page === totalPages}
        onClick={onNext}
        className="px-3 py-1.5 rounded border border-slate-200 text-sm disabled:opacity-40 hover:bg-slate-50"
      >
        Next →
      </button>
    </div>
  )
}
