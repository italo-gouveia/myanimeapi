import { Link, NavLink, Outlet, useNavigate } from 'react-router-dom'
import { useAuth } from '../lib/auth/AuthProvider'
import { clsx } from 'clsx'

function navLinkClass({ isActive }: { isActive: boolean }) {
  return clsx(
    'px-3 py-2 rounded-md text-sm font-medium transition-colors',
    isActive
      ? 'bg-brand-100 text-brand-700'
      : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100',
  )
}

export function Layout() {
  const { isAuthenticated, isAdmin, username, clearSession } = useAuth()
  const navigate = useNavigate()

  const handleLogout = () => {
    clearSession()
    navigate('/')
  }

  return (
    <div className="min-h-screen flex flex-col">
      <header className="bg-white border-b border-slate-200">
        <div className="max-w-6xl mx-auto px-4 h-16 flex items-center justify-between">
          <Link to="/" className="text-xl font-bold text-brand-700">
            MyAnimeAPI
          </Link>
          <nav className="flex items-center gap-1" aria-label="Primary">
            <NavLink to="/" end className={navLinkClass}>
              Catalog
            </NavLink>
            {isAuthenticated && (
              <>
                <NavLink to="/favorites" className={navLinkClass}>
                  Favorites
                </NavLink>
                <NavLink to="/profile" className={navLinkClass}>
                  Profile
                </NavLink>
                {isAdmin && (
                  <NavLink to="/admin" className={navLinkClass}>
                    Admin
                  </NavLink>
                )}
              </>
            )}
            {isAuthenticated ? (
              <div className="flex items-center gap-3 ml-4">
                <span className="text-sm text-slate-600" data-testid="header-username">
                  {username}
                </span>
                <button
                  type="button"
                  onClick={handleLogout}
                  className="px-3 py-2 rounded-md text-sm font-medium text-slate-600 hover:text-slate-900 hover:bg-slate-100"
                  data-testid="logout-button"
                >
                  Sign out
                </button>
              </div>
            ) : (
              <div className="flex items-center gap-1 ml-4">
                <NavLink to="/login" className={navLinkClass}>
                  Sign in
                </NavLink>
                <NavLink
                  to="/register"
                  className="px-3 py-2 rounded-md text-sm font-medium bg-brand-600 text-white hover:bg-brand-700"
                >
                  Register
                </NavLink>
              </div>
            )}
          </nav>
        </div>
      </header>
      <main className="flex-1 max-w-6xl w-full mx-auto px-4 py-8">
        <Outlet />
      </main>
      <footer className="border-t border-slate-200 py-4 text-center text-xs text-slate-500">
        MyAnimeAPI
      </footer>
    </div>
  )
}
