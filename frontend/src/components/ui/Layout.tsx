import { Outlet, NavLink, useNavigate } from 'react-router-dom'
import { Trophy, Sword, Crown, Search } from 'lucide-react'
import { useState } from 'react'
import clsx from 'clsx'

export default function Layout() {
  const [search, setSearch] = useState('')
  const navigate = useNavigate()

  function handleSearch(e: React.FormEvent) {
    e.preventDefault()
    const id = search.trim()
    if (id) {
      navigate(`/player/${id}`)
      setSearch('')
    }
  }

  return (
    <div className="min-h-screen flex flex-col">
      {/* Navbar */}
      <header className="border-b border-surface-600 bg-surface-800/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 h-16 flex items-center justify-between gap-6">
          {/* Logo */}
          <NavLink to="/" className="flex items-center gap-2 shrink-0">
            <Trophy className="text-brand-400" size={22} />
            <span className="font-bold text-white tracking-tight">
              Trophy<span className="text-brand-400">Collector</span>
            </span>
          </NavLink>

          {/* Search */}
          <form onSubmit={handleSearch} className="flex-1 max-w-md">
            <div className="relative">
              <Search className="absolute left-3 top-1/2 -translate-y-1/2 text-slate-500" size={16} />
              <input
                value={search}
                onChange={e => setSearch(e.target.value)}
                placeholder="Search by Steam ID or username..."
                className="w-full bg-surface-900 border border-surface-600 rounded-lg pl-9 pr-4 py-2 text-sm
                           text-slate-200 placeholder-slate-500 outline-none
                           focus:border-brand-500 focus:ring-1 focus:ring-brand-500 transition"
              />
            </div>
          </form>

          {/* Nav links */}
          <nav className="flex items-center gap-1 shrink-0">
            {[
              { to: '/heroes',      label: 'Heroes',      icon: Sword  },
              { to: '/leaderboard', label: 'Leaderboard', icon: Crown  },
            ].map(({ to, label, icon: Icon }) => (
              <NavLink
                key={to}
                to={to}
                className={({ isActive }) =>
                  clsx(
                    'flex items-center gap-1.5 px-3 py-2 rounded-lg text-sm font-medium transition',
                    isActive
                      ? 'bg-brand-600/20 text-brand-400'
                      : 'text-slate-400 hover:text-slate-200 hover:bg-surface-700',
                  )
                }
              >
                <Icon size={15} />
                {label}
              </NavLink>
            ))}
          </nav>
        </div>
      </header>

      {/* Page content */}
      <main className="flex-1 max-w-7xl mx-auto w-full px-4 py-8">
        <Outlet />
      </main>

      {/* Footer */}
      <footer className="border-t border-surface-600 py-6 text-center text-sm text-slate-500">
        Trophy Collector — Deadlock stats tracker. Not affiliated with Valve.
      </footer>
    </div>
  )
}
