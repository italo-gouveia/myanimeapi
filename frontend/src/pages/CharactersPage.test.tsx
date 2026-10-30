import { describe, it, expect, vi, beforeEach } from 'vitest'
import { render, screen, fireEvent } from '@testing-library/react'
import { CharactersPage } from './CharactersPage'

vi.mock('@tanstack/react-query', () => ({
  useQuery: vi.fn(),
  keepPreviousData: undefined,
}))

vi.mock('../lib/api/characters', () => ({
  listCharacters: vi.fn(),
  searchCharacters: vi.fn(),
}))

vi.mock('../components/CharacterCard', () => ({
  CharacterCard: ({ character }: { character: { id: number; name: string } }) => (
    <div data-testid={`char-card-${character.id}`}>{character.name}</div>
  ),
}))

import { useQuery } from '@tanstack/react-query'

const mockChar = { id: 1, name: 'Naruto Uzumaki', description: 'A ninja.', image_url: '', voice_actor: '' }

function renderPage() {
  return render(<CharactersPage />)
}

describe('CharactersPage', () => {
  beforeEach(() => vi.clearAllMocks())

  it('renders page heading', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: undefined } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('heading', { name: /characters/i })).toBeInTheDocument()
  })

  it('shows loading state', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: true, isError: false, data: undefined } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText(/loading/i)).toBeInTheDocument()
  })

  it('shows error state', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: true, data: undefined } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('alert')).toBeInTheDocument()
    expect(screen.getByText(/could not load characters/i)).toBeInTheDocument()
  })

  it('renders character grid', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false, isError: false,
      data: { data: [mockChar], total: 1 },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByTestId('characters-grid')).toBeInTheDocument()
    expect(screen.getByTestId('char-card-1')).toBeInTheDocument()
    expect(screen.getByText('Naruto Uzumaki')).toBeInTheDocument()
  })

  it('shows "No characters found." when total is 0', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false, isError: false,
      data: { data: [], total: 0 },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('No characters found.')).toBeInTheDocument()
  })

  it('shows pluralized count', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false, isError: false,
      data: { data: [mockChar, { ...mockChar, id: 2, name: 'Sasuke' }], total: 2 },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('2 characters')).toBeInTheDocument()
  })

  it('shows singular count for 1 character', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false, isError: false,
      data: { data: [mockChar], total: 1 },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByText('1 character')).toBeInTheDocument()
  })

  it('renders search form', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: undefined } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('searchbox', { name: /search characters/i })).toBeInTheDocument()
    expect(screen.getByRole('button', { name: /^search$/i })).toBeInTheDocument()
  })

  it('shows pagination when totalPages > 1', () => {
    vi.mocked(useQuery).mockReturnValue({
      isLoading: false, isError: false,
      data: { data: Array(20).fill(mockChar).map((c, i) => ({ ...c, id: i + 1 })), total: 50 },
    } as ReturnType<typeof useQuery>)
    renderPage()
    expect(screen.getByRole('navigation', { name: /pagination/i })).toBeInTheDocument()
  })

  it('submits search form', () => {
    vi.mocked(useQuery).mockReturnValue({ isLoading: false, isError: false, data: undefined } as ReturnType<typeof useQuery>)
    renderPage()
    const input = screen.getByRole('searchbox')
    fireEvent.change(input, { target: { value: 'Naruto' } })
    fireEvent.submit(screen.getByRole('searchbox').closest('form')!)
    expect(input).toBeInTheDocument()
  })
})
