import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { TableSkeleton, Pagination } from './components'

describe('TableSkeleton', () => {
  it('renders default 5 rows with given cols', () => {
    render(
      <table>
        <tbody>
          <TableSkeleton cols={3} />
        </tbody>
      </table>,
    )
    const rows = screen.getAllByRole('row')
    expect(rows).toHaveLength(5)
  })

  it('renders custom number of rows', () => {
    render(
      <table>
        <tbody>
          <TableSkeleton rows={2} cols={4} />
        </tbody>
      </table>,
    )
    expect(screen.getAllByRole('row')).toHaveLength(2)
  })

  it('renders correct number of cells per row', () => {
    render(
      <table>
        <tbody>
          <TableSkeleton rows={1} cols={3} />
        </tbody>
      </table>,
    )
    expect(screen.getAllByRole('cell')).toHaveLength(3)
  })
})

describe('Pagination', () => {
  it('returns null when totalPages <= 1', () => {
    const { container } = render(
      <Pagination page={1} totalPages={1} onPrev={vi.fn()} onNext={vi.fn()} />,
    )
    expect(container.firstChild).toBeNull()
  })

  it('renders page info when totalPages > 1', () => {
    render(
      <Pagination page={2} totalPages={5} onPrev={vi.fn()} onNext={vi.fn()} />,
    )
    expect(screen.getByText('Page 2 of 5')).toBeInTheDocument()
  })

  it('disables Prev button on first page', () => {
    render(
      <Pagination page={1} totalPages={3} onPrev={vi.fn()} onNext={vi.fn()} />,
    )
    expect(screen.getByRole('button', { name: /previous/i })).toBeDisabled()
  })

  it('disables Next button on last page', () => {
    render(
      <Pagination page={3} totalPages={3} onPrev={vi.fn()} onNext={vi.fn()} />,
    )
    expect(screen.getByRole('button', { name: /next/i })).toBeDisabled()
  })

  it('calls onPrev when Previous clicked', async () => {
    const onPrev = vi.fn()
    render(<Pagination page={2} totalPages={3} onPrev={onPrev} onNext={vi.fn()} />)
    await userEvent.click(screen.getByRole('button', { name: /previous/i }))
    expect(onPrev).toHaveBeenCalledOnce()
  })

  it('calls onNext when Next clicked', async () => {
    const onNext = vi.fn()
    render(<Pagination page={1} totalPages={3} onPrev={vi.fn()} onNext={onNext} />)
    await userEvent.click(screen.getByRole('button', { name: /next/i }))
    expect(onNext).toHaveBeenCalledOnce()
  })
})
