import { useState } from 'react'
import { useForm } from 'react-hook-form'
import { zodResolver } from '@hookform/resolvers/zod'
import { z } from 'zod'
import { Link, useNavigate } from 'react-router-dom'
import { useMutation } from '@tanstack/react-query'
import { login, register as registerUser } from '../lib/api/auth'
import { extractErrorMessage } from '../lib/api/client'
import { useAuth } from '../lib/auth/AuthProvider'
import { Button } from '../components/Button'
import { Input } from '../components/Input'

const registerSchema = z.object({
  username: z
    .string()
    .min(3, 'Username must be at least 3 characters')
    .max(50, 'Username is too long'),
  email: z.string().email('Invalid email address'),
  password: z.string().min(5, 'Password must be at least 5 characters'),
})

type RegisterValues = z.infer<typeof registerSchema>

export function RegisterPage() {
  const navigate = useNavigate()
  const { setSession } = useAuth()
  const [submitError, setSubmitError] = useState<string | null>(null)

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
  } = useForm<RegisterValues>({
    resolver: zodResolver(registerSchema),
    defaultValues: { username: '', email: '', password: '' },
  })

  const mutation = useMutation({
    mutationFn: async (values: RegisterValues) => {
      await registerUser(values)
      // Auto-login right after a successful registration so the user lands
      // already signed in. The login endpoint is rate-limited (5/min/IP)
      // which is fine for the human-driven flow.
      const auth = await login({ username: values.username, password: values.password })
      return { token: auth.token, username: values.username }
    },
    onSuccess: ({ token, username }) => {
      setSession(token, username)
      navigate('/', { replace: true })
    },
    onError: (error) => {
      setSubmitError(extractErrorMessage(error, 'Could not create the account.'))
    },
  })

  const onSubmit = (values: RegisterValues) => {
    setSubmitError(null)
    mutation.mutate(values)
  }

  return (
    <div className="max-w-md mx-auto bg-white p-6 rounded-lg border border-slate-200 shadow-sm">
      <h1 className="text-2xl font-bold mb-1">Create an account</h1>
      <p className="text-sm text-slate-600 mb-6">Track your favourite anime.</p>

      <form onSubmit={handleSubmit(onSubmit)} className="flex flex-col gap-4" noValidate>
        <Input
          label="Username"
          autoComplete="username"
          error={errors.username?.message}
          data-testid="register-username"
          {...register('username')}
        />
        <Input
          label="Email"
          type="email"
          autoComplete="email"
          error={errors.email?.message}
          data-testid="register-email"
          {...register('email')}
        />
        <Input
          label="Password"
          type="password"
          autoComplete="new-password"
          error={errors.password?.message}
          data-testid="register-password"
          {...register('password')}
        />

        {submitError && (
          <p className="text-sm text-red-600" role="alert" data-testid="register-error">
            {submitError}
          </p>
        )}

        <Button
          type="submit"
          fullWidth
          disabled={isSubmitting || mutation.isPending}
          data-testid="register-submit"
        >
          {mutation.isPending ? 'Creating account...' : 'Create account'}
        </Button>
      </form>

      <p className="text-sm text-slate-600 mt-6 text-center">
        Already have one?{' '}
        <Link to="/login" className="text-brand-700 font-medium underline">
          Sign in
        </Link>
      </p>
    </div>
  )
}
