import { Routes, Route } from 'react-router-dom'
import Layout from './components/ui/Layout'
import HomePage from './pages/HomePage'
import PlayerPage from './pages/PlayerPage'
import HeroesPage from './pages/HeroesPage'
import LeaderboardPage from './pages/LeaderboardPage'
import NotFoundPage from './pages/NotFoundPage'

export default function App() {
  return (
    <Routes>
      <Route element={<Layout />}>
        <Route path="/" element={<HomePage />} />
        <Route path="/player/:steamId" element={<PlayerPage />} />
        <Route path="/heroes" element={<HeroesPage />} />
        <Route path="/leaderboard" element={<LeaderboardPage />} />
        <Route path="*" element={<NotFoundPage />} />
      </Route>
    </Routes>
  )
}
