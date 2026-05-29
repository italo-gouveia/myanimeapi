import { Link } from 'react-router-dom'

export function NotFoundPage() {
  return (
    <div className="text-center py-16">
      <h1 className="text-3xl font-bold">404</h1>
      <p className="text-slate-600 mt-2">Page not found.</p>
      <Link to="/" className="inline-block mt-4 text-brand-700 underline">
        Back to catalog
      </Link>
    </div>
  )
}
