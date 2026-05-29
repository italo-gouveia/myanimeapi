import { useEffect, useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link, useLocation, useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { login } from '../lib/api/auth'
import { extractErrorMessage } from '../lib/api/client'
import { useAuth } from '../lib/auth/AuthProvider'
import { Button } from '../components/Button'
import { Input } from '../components/Input'

const loginSchema = z.object({
  username: z.string().min(3, 'Username must be at least 3 characters'),
  password: z.string().min(5, 'Password must be at least 5 characters'),
})

type LoginValues = z.infer<typeof loginSchema>

export function LoginPage() {
  const { setSession, isAuthenticated } = useAuth()
  const navigate = useNavigate()
  const location = useLocation()
  const [submitError, setSubmitError] = useState<string | null>(null)

  const from = (location.state as { from?: string } | null)?.from ?? '/'

  useEffect(() => {
    if (isAuthenticated) {
      navigate(from, { replace: true })
    }
  }, [isAuthenticated, navigate, from])

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<LoginValues>({
    resolver: zodResolver(loginSchema),
    defaultValues: { username: '', password: '' },
  })

  const mutation = useMutation({
    mutationFn: login,
    onSuccess: (data, variables) => {
      setSession(data.token, variables.username)
      navigate(from, { replace: true })
    },
    onError: (error) => {
      setSubmitError(extractErrorMessage(error, 'Could not sign in. Check your credentials.'))
    },
  })

  const onSubmit = (values: LoginValues) => {
    setSubmitError(null)
    mutation.mutate(values)
  }

  return (
    <div className="max-w-md mx-auto bg-white p-6 rounded-lg border border-slate-200 shadow-sm">
      <h1 className="text-2xl font-bold mb-1">Sign in</h1>
      <p className="text-sm text-slate-600 mb-6">Welcome back to MyAnimeAPI.</p>

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4" noValidate>
        <Input
          label="Username"
          autoComplete="username"
          error={errors.username?.message}
          data-testid="login-username"
          {...register('username')}
        />
        <Input
          label="Password"
          type="password"
          autoComplete="current-password"
          error={errors.password?.message}
          data-testid="login-password"
          {...register('password')}
        />

        {submitError && (
          <p className="text-sm text-red-600" role="alert" data-testid="login-error">
            {submitError}
          </p>
        )}

        <Button
          type="submit"
          fullWidth
          disabled={isSubmitting || mutation.isPending}
          data-testid="login-submit"
        >
          {mutation.isPending ? 'Signing in...' : 'Sign in'}
        </Button>
      </form>

      <p className="text-sm text-slate-600 mt-6 text-center">
        No account yet?{' '}
        <Link to="/register" className="text-brand-700 font-medium underline">
          Create one
        </Link>
      </p>
    </div>
  )
}
