import { useNavigate } from 'react-router-dom'

export default function NotFoundPage() {
  const navigate = useNavigate()
  return (
    <div className="flex flex-col items-center justify-center min-h-[60vh] gap-4 text-slate-400">
      <span className="text-6xl font-bold text-surface-600">404</span>
      <p>Page not found.</p>
      <button
        onClick={() => navigate('/')}
        className="px-4 py-2 rounded-lg bg-brand-600 hover:bg-brand-500 text-white text-sm font-medium transition"
      >
        Go home
      </button>
    </div>
  )
}
