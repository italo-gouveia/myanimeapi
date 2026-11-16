import { describe, it, expect, vi } from 'vitest'
import { render, screen, act } from '@testing-library/react'
import { ToastContainer, useToast } from './Toast'
import { renderHook } from '@testing-library/react'

describe('ToastContainer', () => {
  it('renders nothing when toasts array is empty', () => {
    const { container } = render(<ToastContainer toasts={[]} onDone={vi.fn()} />)
    expect(container.firstChild).toBeNull()
  })

  it('renders a success toast', () => {
    render(
      <ToastContainer
        toasts={[{ id: 1, text: 'Saved!', type: 'success' }]}
        onDone={vi.fn()}
      />,
    )
    expect(screen.getByText(/saved!/i)).toBeInTheDocument()
  })

  it('renders an error toast', () => {
    render(
      <ToastContainer
        toasts={[{ id: 2, text: 'Something went wrong', type: 'error' }]}
        onDone={vi.fn()}
      />,
    )
    expect(screen.getByText(/something went wrong/i)).toBeInTheDocument()
  })

  it('renders an info toast', () => {
    render(
      <ToastContainer
        toasts={[{ id: 3, text: 'Info message', type: 'info' }]}
        onDone={vi.fn()}
      />,
    )
    expect(screen.getByText(/info message/i)).toBeInTheDocument()
  })

  it('renders multiple toasts', () => {
    render(
      <ToastContainer
        toasts={[
          { id: 1, text: 'First', type: 'success' },
          { id: 2, text: 'Second', type: 'error' },
        ]}
        onDone={vi.fn()}
      />,
    )
    expect(screen.getByText(/first/i)).toBeInTheDocument()
    expect(screen.getByText(/second/i)).toBeInTheDocument()
  })

  it('renders toast with accessible role', () => {
    render(
      <ToastContainer
        toasts={[{ id: 1, text: 'Done', type: 'success' }]}
        onDone={vi.fn()}
      />,
    )
    expect(screen.getByRole('status')).toBeInTheDocument()
  })
})

describe('useToast hook', () => {
  it('starts with no toasts', () => {
    const { result } = renderHook(() => useToast())
    expect(result.current.toasts).toHaveLength(0)
  })

  it('adds a toast with push', () => {
    const { result } = renderHook(() => useToast())
    act(() => { result.current.push('Hello', 'success') })
    expect(result.current.toasts).toHaveLength(1)
    expect(result.current.toasts[0].text).toBe('Hello')
    expect(result.current.toasts[0].type).toBe('success')
  })

  it('removes a toast with remove', () => {
    const { result } = renderHook(() => useToast())
    act(() => { result.current.push('Bye', 'info') })
    const id = result.current.toasts[0].id
    act(() => { result.current.remove(id) })
    expect(result.current.toasts).toHaveLength(0)
  })

  it('defaults to success type', () => {
    const { result } = renderHook(() => useToast())
    act(() => { result.current.push('Default') })
    expect(result.current.toasts[0].type).toBe('success')
  })
})
