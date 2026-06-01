import { Routes, Route } from 'react-router-dom'
import { Layout } from '../components/Layout'
import { HomePage } from '../pages/HomePage'
import { AnimeDetailPage } from '../pages/AnimeDetailPage'
import { LoginPage } from '../pages/LoginPage'
import { RegisterPage } from '../pages/RegisterPage'
import { FavoritesPage } from '../pages/FavoritesPage'
import { WatchlistPage } from '../pages/WatchlistPage'
import { ProfilePage } from '../pages/ProfilePage'
import { NotFoundPage } from '../pages/NotFoundPage'
import { AdminLayout } from '../pages/admin/AdminLayout'
import { AdminDashboard } from '../pages/admin/AdminDashboard'
import { AdminETLPage } from '../pages/admin/AdminETLPage'
import { AdminUsersPage } from '../pages/admin/AdminUsersPage'
import { ProtectedRoute } from './ProtectedRoute'

export function AppRoutes() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route index element={<HomePage />} />
        <Route path="animes/:id" element={<AnimeDetailPage />} />
        <Route path="login" element={<LoginPage />} />
        <Route path="register" element={<RegisterPage />} />

        {/* Authenticated routes */}
        <Route element={<ProtectedRoute />}>
          <Route path="favorites" element={<FavoritesPage />} />
          <Route path="watchlist" element={<WatchlistPage />} />
          <Route path="profile" element={<ProfilePage />} />
        </Route>

        {/* Admin-only routes */}
        <Route element={<ProtectedRoute requireAdmin />}>
          <Route path="admin" element={<AdminLayout />}>
            <Route index element={<AdminDashboard />} />
            <Route path="etl" element={<AdminETLPage />} />
            <Route path="users" element={<AdminUsersPage />} />
          </Route>
        </Route>

        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  )
}
