import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { CharacterCard } from './CharacterCard'
import type { Character } from '../lib/api/types'

const mockCharacter: Character = {
  id: 1,
  name: 'Naruto Uzumaki',
  description: 'The seventh Hokage.',
  voice_actor: 'Junko Takeuchi',
  image_url: 'https://example.com/naruto.jpg',
}

describe('CharacterCard', () => {
  it('renders character name', () => {
    render(<CharacterCard character={mockCharacter} />)
    expect(screen.getByText('Naruto Uzumaki')).toBeInTheDocument()
  })

  it('renders character image when image_url is set', () => {
    render(<CharacterCard character={mockCharacter} />)
    const img = screen.getByRole('img', { name: 'Naruto Uzumaki' })
    expect(img).toHaveAttribute('src', 'https://example.com/naruto.jpg')
  })

  it('renders initials placeholder when image_url is absent', () => {
    render(<CharacterCard character={{ ...mockCharacter, image_url: '' }} />)
    expect(screen.getByText('NU')).toBeInTheDocument()
  })

  it('renders voice actor when present', () => {
    render(<CharacterCard character={mockCharacter} />)
    expect(screen.getByText(/Junko Takeuchi/)).toBeInTheDocument()
  })

  it('does not render voice actor row when absent', () => {
    render(<CharacterCard character={{ ...mockCharacter, voice_actor: '' }} />)
    expect(screen.queryByText(/Junko Takeuchi/)).not.toBeInTheDocument()
  })

  it('renders description when present', () => {
    render(<CharacterCard character={mockCharacter} />)
    expect(screen.getByText('The seventh Hokage.')).toBeInTheDocument()
  })

  it('has correct data-testid', () => {
    render(<CharacterCard character={mockCharacter} />)
    expect(screen.getByTestId('character-card-1')).toBeInTheDocument()
  })
})
