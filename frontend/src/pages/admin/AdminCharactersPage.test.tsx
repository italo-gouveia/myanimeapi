/* eslint-disable @typescript-eslint/no-explicit-any */
import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent, act } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { AdminCharactersPage } from './AdminCharactersPage'

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
  useMutation: vi.fn(),
  useQueryClient: () => ({ invalidateQueries: vi.fn() }),
}))

vi.mock('../../lib/api/characters', () => ({
  listCharacters: vi.fn(),
  createCharacter: vi.fn(),
  updateCharacter: vi.fn(),
  deleteCharacter: vi.fn(),
}))

import { useQuery, useMutation } from '@tanstack/react-query'

const mockChar = { id: 1, name: 'Naruto Uzumaki', voice_actor: 'Junko Takeuchi', description: 'The main hero', image_url: '' }

function setupMutation(overrides: Record<string, any> = {}) {
  vi.mocked(useMutation).mockReturnValue({
    mutate: vi.fn(),
    isPending: false,
    isError: false,
    error: null,
    ...overrides,
  } as any)
}

function renderPage() {
  return render(
    <MemoryRouter>
      <AdminCharactersPage />
    </MemoryRouter>,
  )
}

describe('AdminCharactersPage', () => {
  beforeEach(() => {
    vi.clearAllMocks()
    setupMutation()
  })

  it('shows loading skeleton', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: true, isError: false, data: undefined } as any)
    renderPage()
    const cells = document.querySelectorAll('td .animate-pulse')
    expect(cells.length).toBeGreaterThan(0)
  })

  it('shows error state', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: true, data: undefined } as any)
    renderPage()
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/failed to load characters/i)).toBeInTheDocument()
  })

  it('renders character rows', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [mockChar], total: 1, page: 1, limit: 20 },
    } as any)
    renderPage()
    expect(screen.getByText('Naruto Uzumaki')).toBeInTheDocument()
    expect(screen.getByText('Junko Takeuchi')).toBeInTheDocument()
  })

  it('shows total count', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [mockChar], total: 42, page: 1, limit: 20 },
    } as any)
    renderPage()
    expect(screen.getByText('42 total')).toBeInTheDocument()
  })

  it('opens create modal on button click', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [], total: 0, page: 1, limit: 20 },
    } as any)
    renderPage()
    fireEvent.click(screen.getByRole('button', { name: /new character/i }))
    expect(screen.getByRole('dialog')).toBeInTheDocument()
    expect(screen.getByText('New character')).toBeInTheDocument()
  })

  it('closes create modal on cancel', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [], total: 0, page: 1, limit: 20 },
    } as any)
    renderPage()
    fireEvent.click(screen.getByRole('button', { name: /new character/i }))
    fireEvent.click(screen.getByRole('button', { name: /cancel/i }))
    expect(screen.queryByRole('dialog')).not.toBeInTheDocument()
  })

  it('opens edit modal and pre-fills name', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [mockChar], total: 1, page: 1, limit: 20 },
    } as any)
    renderPage()
    fireEvent.click(screen.getByRole('button', { name: /edit/i }))
    expect(screen.getByRole('dialog')).toBeInTheDocument()
    expect(screen.getByDisplayValue('Naruto Uzumaki')).toBeInTheDocument()
  })

  it('opens delete confirmation modal', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [mockChar], total: 1, page: 1, limit: 20 },
    } as any)
    renderPage()
    // click the row-level Delete button (first one in the table)
    fireEvent.click(screen.getAllByRole('button', { name: /delete/i })[0])
    const dialog = screen.getByRole('dialog')
    expect(dialog).toBeInTheDocument()
    expect(dialog).toHaveTextContent(/delete character/i)
    expect(dialog).toHaveTextContent(/naruto uzumaki/i)
  })

  it('calls createCharacter mutate on form submit', () => {
    const mutate = vi.fn()
    vi.mocked(useMutation).mockReturnValue({ mutate, isPending: false, isError: false, error: null } as any)
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [], total: 0, page: 1, limit: 20 },
    } as any)
    renderPage()
    fireEvent.click(screen.getByRole('button', { name: /new character/i }))
    fireEvent.change(screen.getByLabelText(/name/i), { target: { value: 'Sasuke' } })
    fireEvent.click(screen.getByRole('button', { name: /create/i }))
    expect(mutate).toHaveBeenCalledWith(expect.objectContaining({ name: 'Sasuke' }))
  })

  it('calls deleteCharacter mutate on confirm', () => {
    const mutate = vi.fn()
    vi.mocked(useMutation).mockReturnValue({ mutate, isPending: false, isError: false, error: null } as any)
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [mockChar], total: 1, page: 1, limit: 20 },
    } as any)
    renderPage()
    // open delete modal via row button, then confirm
    fireEvent.click(screen.getAllByRole('button', { name: /delete/i })[0])
    const dialog = screen.getByRole('dialog')
    fireEvent.click(screen.getAllByRole('button', { name: /^delete$/i }).find(
      (btn) => dialog.contains(btn) && btn.classList.contains('bg-red-600'),
    )!)
    expect(mutate).toHaveBeenCalledWith(1)
  })

  it('shows mutation error alert', () => {
    vi.mocked(useMutation).mockReturnValue({ mutate: vi.fn(), isPending: false, isError: true, error: new Error('fail') } as any)
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [], total: 0, page: 1, limit: 20 },
    } as any)
    renderPage()
    expect(screen.getByText(/operation failed/i)).toBeInTheDocument()
  })

  it('covers onSuccess callbacks for all mutations', async () => {
    const callbacks: Array<((...args: any[]) => any) | undefined> = []
    vi.mocked(useMutation).mockImplementation(({ onSuccess }: any = {}) => {
      callbacks.push(onSuccess)
      return { mutate: vi.fn(), isPending: false, isError: false, error: null } as any
    })
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false,
      isError: false,
      data: { data: [], total: 0, page: 1, limit: 20 },
    } as any)
    renderPage()
    await act(async () => {
      for (const cb of callbacks) await cb?.()
    })
  })
})
