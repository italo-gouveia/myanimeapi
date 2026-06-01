import { NavLink, Outlet } from 'react-router-dom'
import { clsx } from 'clsx'

function sideNavClass({ isActive }: { isActive: boolean }) {
  return clsx(
    'flex items-center gap-2 px-3 py-2 rounded-md text-sm font-medium transition-colors',
    isActive
      ? 'bg-brand-100 text-brand-700'
      : 'text-slate-600 hover:text-slate-900 hover:bg-slate-100',
  )
}

export function AdminLayout() {
  return (
    <div className="flex gap-8 min-h-[calc(100vh-10rem)]">
      {/* Sidebar */}
      <aside className="w-48 shrink-0">
        <p className="text-xs font-semibold uppercase tracking-wider text-slate-400 mb-3 px-3">
          Admin
        </p>
        <nav className="flex flex-col gap-1" aria-label="Admin navigation">
          <NavLink to="/admin" end className={sideNavClass}>
            <span>📊</span> Dashboard
          </NavLink>
          <NavLink to="/admin/analytics" className={sideNavClass}>
            <span>📈</span> Analytics
          </NavLink>
          <NavLink to="/admin/etl" className={sideNavClass}>
            <span>🔄</span> ETL Sync
          </NavLink>
          <NavLink to="/admin/users" className={sideNavClass}>
            <span>👥</span> Users
          </NavLink>
        </nav>
      </aside>

      {/* Content */}
      <div className="flex-1 min-w-0">
        <Outlet />
      </div>
    </div>
  )
}
