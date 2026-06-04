import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect } from 'vitest'
import { Input } from './Input'

describe('Input', () => {
  it('renders label and input', () => {
    render(<Input label="Email" name="email" />)
    expect(screen.getByLabelText('Email')).toBeInTheDocument()
  })

  it('associates label with input via name as id fallback', () => {
    render(<Input label="Username" name="username" />)
    const input = screen.getByLabelText('Username')
    expect(input).toHaveAttribute('id', 'username')
  })

  it('prefers explicit id over name', () => {
    render(<Input label="Search" id="search-box" name="q" />)
    expect(screen.getByLabelText('Search')).toHaveAttribute('id', 'search-box')
  })

  it('shows error message and sets aria-invalid', () => {
    render(<Input label="Email" name="email" error="Required" />)
    expect(screen.getByRole('alert')).toHaveTextContent('Required')
    expect(screen.getByLabelText('Email')).toHaveAttribute('aria-invalid', 'true')
  })

  it('links input to error via aria-describedby', () => {
    render(<Input label="Email" name="email" error="Invalid email" />)
    const input = screen.getByLabelText('Email')
    const errorEl = screen.getByRole('alert')
    expect(input.getAttribute('aria-describedby')).toBe(errorEl.id)
  })

  it('does not render error element when error is absent', () => {
    render(<Input label="Email" name="email" />)
    expect(screen.queryByRole('alert')).not.toBeInTheDocument()
    expect(screen.getByLabelText('Email')).not.toHaveAttribute('aria-invalid')
  })

  it('forwards user typing to the input', async () => {
    render(<Input label="Name" name="name" />)
    const input = screen.getByLabelText('Name')
    await userEvent.type(input, 'Italo')
    expect(input).toHaveValue('Italo')
  })

  it('forwards placeholder', () => {
    render(<Input label="Search" name="q" placeholder="Type here…" />)
    expect(screen.getByPlaceholderText('Type here…')).toBeInTheDocument()
  })
})
