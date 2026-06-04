import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, it, expect, vi } from 'vitest'
import { Button } from './Button'

describe('Button', () => {
  it('renders children', () => {
    render(<Button>Click me</Button>)
    expect(screen.getByRole('button', { name: 'Click me' })).toBeInTheDocument()
  })

  it('defaults to type=button to avoid accidental form submission', () => {
    render(<Button>Save</Button>)
    expect(screen.getByRole('button')).toHaveAttribute('type', 'button')
  })

  it('forwards type=submit when specified', () => {
    render(<Button type="submit">Submit</Button>)
    expect(screen.getByRole('button')).toHaveAttribute('type', 'submit')
  })

  it('calls onClick when clicked', async () => {
    const onClick = vi.fn()
    render(<Button onClick={onClick}>Go</Button>)
    await userEvent.click(screen.getByRole('button'))
    expect(onClick).toHaveBeenCalledOnce()
  })

  it('does not call onClick when disabled', async () => {
    const onClick = vi.fn()
    render(<Button disabled onClick={onClick}>Disabled</Button>)
    await userEvent.click(screen.getByRole('button'))
    expect(onClick).not.toHaveBeenCalled()
  })

  it('applies fullWidth class when fullWidth=true', () => {
    render(<Button fullWidth>Wide</Button>)
    expect(screen.getByRole('button')).toHaveClass('w-full')
  })

  it('renders primary variant by default', () => {
    render(<Button>Primary</Button>)
    expect(screen.getByRole('button')).toHaveClass('bg-brand-600')
  })

  it('renders danger variant', () => {
    render(<Button variant="danger">Delete</Button>)
    expect(screen.getByRole('button')).toHaveClass('bg-red-600')
  })

  it('renders secondary variant', () => {
    render(<Button variant="secondary">Cancel</Button>)
    expect(screen.getByRole('button')).toHaveClass('bg-white')
  })

  it('merges extra className', () => {
    render(<Button className="my-custom">Custom</Button>)
    expect(screen.getByRole('button')).toHaveClass('my-custom')
  })
})
