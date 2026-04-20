import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Trophy, Search, Zap } from 'lucide-react'

export default function HomePage() {
  const [input, setInput] = useState('')
  const navigate = useNavigate()

  function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    const id = input.trim()
    if (id) navigate(`/player/${id}`)
  }

  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] gap-12">
      {/* Hero section */}
      <div className="text-center space-y-4">
        <div className="flex justify-center mb-6">
          <div className="p-4 rounded-2xl bg-brand-600/10 border border-brand-600/20">
            <Trophy size={48} className="text-brand-400" />
          </div>
        </div>
        <h1 className="text-5xl font-bold text-white tracking-tight">
          Trophy Collector
        </h1>
        <p className="text-lg text-slate-400 max-w-md">
          Track your Deadlock stats, hero performance, and match history.
        </p>
      </div>

      {/* Search */}
      <form onSubmit={handleSubmit} className="w-full max-w-lg space-y-3">
        <div className="relative">
          <Search className="absolute left-4 top-1/2 -translate-y-1/2 text-slate-500" size={18} />
          <input
            value={input}
            onChange={e => setInput(e.target.value)}
            placeholder="Enter Steam ID or username..."
            className="w-full bg-surface-800 border border-surface-600 rounded-xl
                       pl-12 pr-4 py-4 text-slate-200 placeholder-slate-500
                       outline-none focus:border-brand-500 focus:ring-1 focus:ring-brand-500
                       text-lg transition"
            autoFocus
          />
        </div>
        <button
          type="submit"
          className="w-full py-3 rounded-xl bg-brand-600 hover:bg-brand-500
                     text-white font-semibold transition active:scale-95"
        >
          Search Player
        </button>
      </form>

      {/* Feature highlights */}
      <div className="grid grid-cols-3 gap-4 w-full max-w-lg text-center">
        {[
          { icon: Zap,    label: 'Real-time stats'  },
          { icon: Trophy, label: 'Hero analytics'   },
          { icon: Search, label: 'Match history'    },
        ].map(({ icon: Icon, label }) => (
          <div key={label} className="card flex flex-col items-center gap-2 py-4">
            <Icon size={20} className="text-brand-400" />
            <span className="text-xs text-slate-400">{label}</span>
          </div>
        ))}
      </div>
    </div>
  )
}
