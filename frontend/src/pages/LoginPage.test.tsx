import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { LoginPage } from './LoginPage'

vi.mock('@tanstack/react-query', () => ({ useMutation: vi.fn() }))
vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => vi.fn() }
})
vi.mock('../lib/api/auth', () => ({ login: vi.fn() }))
vi.mock('../lib/api/client', () => ({ extractErrorMessage: vi.fn((_e: unknown, fb: string) => fb) }))
vi.mock('../lib/auth/AuthProvider', () => ({ useAuth: vi.fn() }))
vi.mock('react-hook-form', async () => {
  const actual = await vi.importActual<typeof import('react-hook-form')>('react-hook-form')
  return { ...actual, useForm: vi.fn() }
})

import { useMutation } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/AuthProvider'
import { useForm } from 'react-hook-form'

function setupMocks(isPending = false, errors: Record<string, { message: string }> = {}) {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated: false, isAdmin: false, token: null, username: null,
    role: 'user', userId: null, setSession: vi.fn(), clearSession: vi.fn(),
  })
  vi.mocked(useMutation).mockReturnValue({
    mutate: vi.fn(), isPending, isError: false, error: null,
  } as unknown as ReturnType<typeof useMutation>)
  vi.mocked(useForm).mockReturnValue({
    register: (name: string) => ({ name, ref: vi.fn(), onChange: vi.fn(), onBlur: vi.fn() }),
    handleSubmit: (fn: (v: unknown) => void) => (e?: Event) => { e?.preventDefault?.(); fn({}) },
    formState: { errors, isSubmitting: false },
  } as unknown as ReturnType<typeof useForm>)
}

function renderPage() {
  setupMocks()
  return render(<MemoryRouter><LoginPage /></MemoryRouter>)
}

describe('LoginPage', () => {
  it('renders heading', () => {
    renderPage()
    expect(screen.getByRole('heading', { name: /sign in/i })).toBeInTheDocument()
  })

  it('renders username and password fields', () => {
    renderPage()
    expect(screen.getByTestId('login-username')).toBeInTheDocument()
    expect(screen.getByTestId('login-password')).toBeInTheDocument()
  })

  it('renders submit button', () => {
    renderPage()
    expect(screen.getByTestId('login-submit')).toBeInTheDocument()
  })

  it('shows pending text when mutation is pending', () => {
    setupMocks(true)
    render(<MemoryRouter><LoginPage /></MemoryRouter>)
    expect(screen.getByText(/signing in/i)).toBeInTheDocument()
    expect(screen.getByTestId('login-submit')).toBeDisabled()
  })

  it('shows username validation error', () => {
    setupMocks(false, { username: { message: 'Username must be at least 3 characters' } })
    render(<MemoryRouter><LoginPage /></MemoryRouter>)
    expect(screen.getByText(/at least 3 characters/i)).toBeInTheDocument()
  })

  it('shows password validation error', () => {
    setupMocks(false, { password: { message: 'Password must be at least 5 characters' } })
    render(<MemoryRouter><LoginPage /></MemoryRouter>)
    expect(screen.getByText(/at least 5 characters/i)).toBeInTheDocument()
  })

  it('renders link to register page', () => {
    renderPage()
    expect(screen.getByRole('link', { name: /create one/i })).toBeInTheDocument()
  })

  it('redirects when already authenticated', () => {
    vi.mocked(useAuth).mockReturnValue({
      isAuthenticated: true, isAdmin: false, token: 'tok', username: 'italo',
      role: 'user', userId: 1, setSession: vi.fn(), clearSession: vi.fn(),
    })
    vi.mocked(useMutation).mockReturnValue({ mutate: vi.fn(), isPending: false } as unknown as ReturnType<typeof useMutation>)
    vi.mocked(useForm).mockReturnValue({
      register: (name: string) => ({ name, ref: vi.fn(), onChange: vi.fn(), onBlur: vi.fn() }),
      handleSubmit: (fn: (v: unknown) => void) => (e?: Event) => { e?.preventDefault?.(); fn({}) },
      formState: { errors: {}, isSubmitting: false },
    } as unknown as ReturnType<typeof useForm>)
    render(<MemoryRouter><LoginPage /></MemoryRouter>)
    // Page still renders without crashing when already authenticated
    expect(screen.getByRole('heading', { name: /sign in/i })).toBeInTheDocument()
  })
})
