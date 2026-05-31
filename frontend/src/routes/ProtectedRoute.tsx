import { Navigate, Outlet, useLocation } from 'react-router-dom'
import { useAuth } from '../lib/auth/AuthProvider'

interface ProtectedRouteProps {
  requireAdmin?: boolean
}

export function ProtectedRoute({ requireAdmin = false }: ProtectedRouteProps) {
  const { isAuthenticated, isAdmin } = useAuth()
  const location = useLocation()

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />
  }

  if (requireAdmin && !isAdmin) {
    return <Navigate to="/" replace />
  }

  return <Outlet />
}
