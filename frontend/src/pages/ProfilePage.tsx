import { useState } from 'react'
import { useQuery } from '@tanstack/react-query'
import { useNavigate } from 'react-router-dom'
import { getProfile } from '../lib/api/users'
import { exportUserData, type ExportFormat } from '../lib/api/export'
import { useAuth } from '../lib/auth/AuthProvider'
import { Button } from '../components/Button'

function Avatar({ username }: { username: string }) {
  const initials = username
    .split(/[\s_-]/)
    .map((w) => w[0]?.toUpperCase() ?? '')
    .slice(0, 2)
    .join('')

  return (
    <div
      className="w-20 h-20 rounded-full bg-brand-600 text-white flex items-center justify-center text-2xl font-bold select-none"
      aria-label={`${username} avatar`}
    >
      {initials}
    </div>
  )
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString(undefined, {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  })
}

export function ProfilePage() {
  const { clearSession } = useAuth()
  const navigate = useNavigate()

  const { data: user, isLoading, isError, error } = useQuery({
    queryKey: ['profile'],
    queryFn: getProfile,
  })

  const [exporting, setExporting] = useState(false)

  function handleLogout() {
    clearSession()
    navigate('/', { replace: true })
  }

  async function handleExport(format: ExportFormat) {
    setExporting(true)
    try { await exportUserData(format) } finally { setExporting(false) }
  }

  return (
    <div className="flex flex-col gap-8 max-w-lg">
      <header className="flex flex-col gap-2">
        <h1 className="text-3xl font-bold">Profile</h1>
        <p className="text-slate-600">Your account details.</p>
      </header>

      {isLoading && <p className="text-slate-500">Loading profile...</p>}

      {isError && (
        <p className="text-red-600" role="alert">
          Could not load profile. {(error as Error).message}
        </p>
      )}

      {user && (
        <div className="bg-white border border-slate-200 rounded-xl p-6 flex flex-col gap-6">
          {/* Avatar + name */}
          <div className="flex items-center gap-4">
            {user.profile_pic ? (
              <img
                src={user.profile_pic}
                alt={`${user.username} avatar`}
                className="w-20 h-20 rounded-full object-cover"
              />
            ) : (
              <Avatar username={user.username} />
            )}
            <div>
              <p className="text-xl font-bold text-slate-900" data-testid="profile-username">
                {user.username}
              </p>
              {user.is_admin && (
                <span className="inline-block bg-brand-100 text-brand-700 text-xs font-semibold px-2 py-0.5 rounded mt-1">
                  Admin
                </span>
              )}
            </div>
          </div>

          {/* Details */}
          <dl className="grid grid-cols-[auto_1fr] gap-x-4 gap-y-3 text-sm">
            <dt className="font-medium text-slate-500">Email</dt>
            <dd className="text-slate-900" data-testid="profile-email">
              {user.email}
            </dd>

            {user.bio && (
              <>
                <dt className="font-medium text-slate-500">Bio</dt>
                <dd className="text-slate-700 whitespace-pre-line" data-testid="profile-bio">
                  {user.bio}
                </dd>
              </>
            )}

            <dt className="font-medium text-slate-500">Member since</dt>
            <dd className="text-slate-900" data-testid="profile-since">
              {formatDate(user.created_at)}
            </dd>

            <dt className="font-medium text-slate-500">Status</dt>
            <dd>
              <span
                className={`inline-block text-xs font-semibold px-2 py-0.5 rounded ${
                  user.is_active
                    ? 'bg-green-100 text-green-700'
                    : 'bg-slate-200 text-slate-500'
                }`}
              >
                {user.is_active ? 'Active' : 'Inactive'}
              </span>
            </dd>
          </dl>

          {/* Export data */}
          <div className="pt-2 border-t border-slate-100 flex flex-col gap-2">
            <p className="text-xs font-semibold text-slate-500 uppercase tracking-wide">Export your data</p>
            <div className="flex gap-2">
              <Button
                variant="secondary"
                disabled={exporting}
                onClick={() => handleExport('json')}
                data-testid="export-json"
              >
                {exporting ? 'Exporting…' : '↓ JSON'}
              </Button>
              <Button
                variant="secondary"
                disabled={exporting}
                onClick={() => handleExport('csv')}
                data-testid="export-csv"
              >
                {exporting ? 'Exporting…' : '↓ CSV'}
              </Button>
            </div>
          </div>

          {/* Logout */}
          <div className="pt-2 border-t border-slate-100">
            <Button variant="secondary" onClick={handleLogout} data-testid="profile-logout">
              Sign out
            </Button>
          </div>
        </div>
      )}
    </div>
  )
}
