import { describe, it, expect, vi } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router-dom'
import { RegisterPage } from './RegisterPage'

vi.mock('@tanstack/react-query', () => ({
  useMutation: vi.fn(),
}))

vi.mock('react-router-dom', async () => {
  const actual = await vi.importActual<typeof import('react-router-dom')>('react-router-dom')
  return { ...actual, useNavigate: () => vi.fn() }
})

vi.mock('../lib/api/auth', () => ({
  login: vi.fn(),
  register: vi.fn(),
}))

vi.mock('../lib/api/client', () => ({
  extractErrorMessage: vi.fn((_err, fallback) => fallback),
}))

vi.mock('../lib/auth/AuthProvider', () => ({
  useAuth: vi.fn(),
}))

// @hookform/resolvers v3 is not compatible with Zod v4 — zodResolver throws instead
// of returning errors, so we mock react-hook-form for the validation-error tests.
vi.mock('react-hook-form', async () => {
  const actual = await vi.importActual<typeof import('react-hook-form')>('react-hook-form')
  return { ...actual, useForm: vi.fn() }
})

import { useMutation } from '@tanstack/react-query'
import { useAuth } from '../lib/auth/AuthProvider'
import { useForm } from 'react-hook-form'

function setupMutationMocks(isPending = false) {
  vi.mocked(useAuth).mockReturnValue({
    isAuthenticated: false, isAdmin: false, token: null, username: null,
    role: 'user', userId: null, setSession: vi.fn(), clearSession: vi.fn(),
  })
  vi.mocked(useMutation).mockReturnValue({
    mutate: vi.fn(), isPending, isError: false, error: null,
  } as unknown as ReturnType<typeof useMutation>)
}

function setupFormMock(errors: Record<string, { message: string }> = {}) {
  vi.mocked(useForm).mockReturnValue({
    register: (name: string) => ({ name, ref: vi.fn(), onChange: vi.fn(), onBlur: vi.fn() }),
    handleSubmit: (fn: (v: unknown) => void) => (e?: Event) => { e?.preventDefault?.(); fn({}) },
    formState: { errors, isSubmitting: false },
  } as unknown as ReturnType<typeof useForm>)
}

function renderPage() {
  setupMutationMocks()
  setupFormMock()
  return render(<MemoryRouter><RegisterPage /></MemoryRouter>)
}

describe('RegisterPage', () => {
  it('renders all form fields', () => {
    renderPage()
    expect(screen.getByTestId('register-username')).toBeInTheDocument()
    expect(screen.getByTestId('register-email')).toBeInTheDocument()
    expect(screen.getByTestId('register-password')).toBeInTheDocument()
    expect(screen.getByTestId('register-submit')).toBeInTheDocument()
  })

  it('renders heading', () => {
    renderPage()
    expect(screen.getByRole('heading', { name: /create an account/i })).toBeInTheDocument()
  })

  it('shows validation error for short username', () => {
    setupMutationMocks()
    setupFormMock({ username: { message: 'Username must be at least 3 characters' } })
    render(<MemoryRouter><RegisterPage /></MemoryRouter>)
    expect(screen.getByText(/at least 3 characters/i)).toBeInTheDocument()
  })

  it('shows validation error for invalid email', () => {
    setupMutationMocks()
    setupFormMock({ email: { message: 'Invalid email address' } })
    render(<MemoryRouter><RegisterPage /></MemoryRouter>)
    expect(screen.getByText(/invalid email/i)).toBeInTheDocument()
  })

  it('shows validation error for short password', () => {
    setupMutationMocks()
    setupFormMock({ password: { message: 'Password must be at least 5 characters' } })
    render(<MemoryRouter><RegisterPage /></MemoryRouter>)
    expect(screen.getByText(/at least 5 characters/i)).toBeInTheDocument()
  })

  it('renders sign-in link', () => {
    renderPage()
    expect(screen.getByRole('link', { name: /sign in/i })).toBeInTheDocument()
  })

  it('shows pending state on button when mutation is pending', () => {
    setupMutationMocks(true)
    setupFormMock()
    render(<MemoryRouter><RegisterPage /></MemoryRouter>)
    expect(screen.getByTestId('register-submit')).toBeDisabled()
    expect(screen.getByText(/creating account/i)).toBeInTheDocument()
  })
})
