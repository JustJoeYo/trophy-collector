import { Routes, Route, Navigate } from 'react-router-dom'
import Navbar from './components/Navbar'
import HomePage from './pages/HomePage'
import HeroPage from './pages/HeroPage'
import PlayerDetailPage from './pages/PlayerDetailPage'


export default function App() {
  return (
    <>
      <Navbar />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/player/:id" element={<PlayerDetailPage />} />
        <Route path="/heroes" element={<HeroPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </>
  )
}


