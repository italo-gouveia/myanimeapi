import { clsx } from 'clsx'
import { forwardRef, type InputHTMLAttributes } from 'react'

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label: string
  error?: string
}

export const Input = forwardRef<HTMLInputElement, InputProps>(function Input(
  { label, error, id, className, ...rest },
  ref,
) {
  const inputId = id ?? rest.name
  const errorId = error ? `${inputId}-error` : undefined
  return (
    <div className="flex flex-col gap-1">
      <label htmlFor={inputId} className="text-sm font-medium text-slate-700">
        {label}
      </label>
      <input
        ref={ref}
        id={inputId}
        aria-invalid={error ? true : undefined}
        aria-describedby={errorId}
        className={clsx(
          'px-3 py-2 rounded-md border bg-white text-slate-900 text-sm',
          'focus:outline-none focus:ring-2 focus:ring-brand-500',
          error ? 'border-red-400' : 'border-slate-300',
          className,
        )}
        {...rest}
      />
      {error && (
        <span id={errorId} className="text-xs text-red-600" role="alert">
          {error}
        </span>
      )}
    </div>
  )
})
