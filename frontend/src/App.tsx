import { Routes, Route, Navigate } from 'react-router-dom'
import Navbar from './components/Navbar'
import PlayerPage from './pages/PlayerPage'
import HomePage from './pages/HomePage'
import HeroPage from './pages/HeroPage'


export default function App() {
  return (
    <>
      <Navbar />
      <Routes>
        <Route path="/" element={<HomePage />} />
        <Route path="/player" element={<PlayerPage />} />
        <Route path="/hero" element={<HeroPage />} />
        <Route path="*" element={<Navigate to="/" replace />} />
      </Routes>
    </>
  )
}


