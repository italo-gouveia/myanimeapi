import { useEffect, useState } from 'react'

export interface ToastMessage {
  id: number
  text: string
  type?: 'success' | 'info' | 'error'
}

interface ToastItemProps {
  toast: ToastMessage
  onDone: (id: number) => void
}

function ToastItem({ toast, onDone }: ToastItemProps) {
  const [visible, setVisible] = useState(false)

  useEffect(() => {
    // Slide in
    const t1 = window.setTimeout(() => setVisible(true), 10)
    // Start fade-out after 2.3s
    const t2 = window.setTimeout(() => setVisible(false), 2300)
    // Remove after fade completes
    const t3 = window.setTimeout(() => onDone(toast.id), 2700)
    return () => { clearTimeout(t1); clearTimeout(t2); clearTimeout(t3) }
  }, [toast.id, onDone])

  const colors = {
    success: 'bg-green-600',
    info:    'bg-brand-600',
    error:   'bg-red-600',
  }[toast.type ?? 'success']

  return (
    <div
      className={`
        flex items-center gap-2 px-4 py-2.5 rounded-xl text-white text-sm font-medium shadow-lg
        transition-all duration-300 ease-out
        ${colors}
        ${visible ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-2'}
      `}
      role="status"
      aria-live="polite"
    >
      {toast.type === 'error' ? '✕' : '✓'} {toast.text}
    </div>
  )
}

interface ToastContainerProps {
  toasts: ToastMessage[]
  onDone: (id: number) => void
}

export function ToastContainer({ toasts, onDone }: ToastContainerProps) {
  if (toasts.length === 0) return null
  return (
    <div className="fixed bottom-6 left-1/2 -translate-x-1/2 z-50 flex flex-col items-center gap-2 pointer-events-none">
      {toasts.map((t) => (
        <ToastItem key={t.id} toast={t} onDone={onDone} />
      ))}
    </div>
  )
}

// ── Hook ──────────────────────────────────────────────────────────────────────
let _id = 0

export function useToast() {
  const [toasts, setToasts] = useState<ToastMessage[]>([])

  const push = (text: string, type: ToastMessage['type'] = 'success') => {
    const id = ++_id
    setToasts((prev) => [...prev, { id, text, type }])
  }

  const remove = (id: number) => setToasts((prev) => prev.filter((t) => t.id !== id))

  return { toasts, push, remove }
}
