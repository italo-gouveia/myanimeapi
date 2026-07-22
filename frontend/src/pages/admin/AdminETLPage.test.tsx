import { render, screen, fireEvent } from '@testing-library/react'
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { AdminETLPage } from './AdminETLPage'

// Mock @tanstack/react-query so we control mutate/isPending/isError without a real QueryClient
const mockMutate = vi.fn()
let mockIsPending = false
let mockIsError = false
let mockError: Error | null = null

vi.mock('@tanstack/react-query', () => ({
  useMutation: (_opts: { mutationFn: () => unknown; onSuccess?: (data: unknown) => void }) => ({
    mutate: mockMutate,
    isPending: mockIsPending,
    isError: mockIsError,
    error: mockError,
  }),
}))

// Mock the admin API module — we test UI behaviour, not HTTP
vi.mock('../../lib/api/admin', () => ({
  triggerETLSync: vi.fn(),
}))

beforeEach(() => {
  mockMutate.mockReset()
  mockIsPending = false
  mockIsError = false
  mockError = null
})

describe('AdminETLPage', () => {
  it('renders the ETL sync UI elements (heading, input, button)', () => {
    render(<AdminETLPage />)

    expect(screen.getByRole('heading', { name: /etl sync/i })).toBeInTheDocument()

    const input = screen.getByLabelText(/pages to sync/i)
    expect(input).toBeInTheDocument()
    expect(input).toHaveValue(5) // default value

    expect(screen.getByRole('button', { name: /start sync/i })).toBeInTheDocument()
  })

  it('renders the label text "Pages to sync" with the spacing fix applied', () => {
    render(<AdminETLPage />)

    // The label contains "Pages to sync" followed by a hint span.
    // The {' '} fix ensures there is a space between the label text and the span
    // so the accessible name includes both parts correctly.
    const label = screen.getByText(/pages to sync/i)
    expect(label).toBeInTheDocument()
    // The hint "(1 page = 25 anime)" is also visible in the document
    expect(screen.getByText(/1 page = 25 anime/i)).toBeInTheDocument()
  })

  it('shows admin-only ETL content: Jikan API description and limit note', () => {
    render(<AdminETLPage />)

    // Description paragraph visible only on the admin ETL page
    expect(screen.getByText(/jikan api/i)).toBeInTheDocument()
    expect(screen.getByText(/max 20 pages/i)).toBeInTheDocument()
    // The Jikan link should point to the correct URL
    const jikanLink = screen.getByRole('link', { name: /jikan api/i })
    expect(jikanLink).toHaveAttribute('href', 'https://jikan.moe')
  })

  it('updates the pages count when the user changes the input', async () => {
    render(<AdminETLPage />)

    const input = screen.getByLabelText(/pages to sync/i)
    // Use fireEvent.change to set an exact value on the controlled number input —
    // userEvent.type appends characters and triggers intermediate clamped values.
    fireEvent.change(input, { target: { value: '8' } })

    expect(input).toHaveValue(8)
    // The helper text should reflect the new value (8 * 25 = 200)
    expect(screen.getByText(/200 anime/i)).toBeInTheDocument()
  })
})
