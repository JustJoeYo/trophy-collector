import PlayerSearch from '../components/PlayerSearch'
import hazePng from '../assets/haze.png'
import './HomePage.css'

export default function HomePage() {
  return (
    <div className="home">
      <img
        className="home-haze"
        src={hazePng}
        alt=""
        aria-hidden="true"
      />
      <div className="home-glow" aria-hidden="true" />
      <div className="home-bg-grid" aria-hidden="true" />
      <div className="home-hero">
        <p className="home-eyebrow">Deadlock Stats Tracker</p>
        <h1 className="home-title">
          <span className="home-title-trophy">Trophy</span>
          <span className="home-title-collector">Collector</span>
        </h1>
        <p className="home-sub">
          Full match history, hero performance, lane stats, and awards — pulled from every game you've ever played.
        </p>
        <div className="home-search-wrap">
          <PlayerSearch />
        </div>
        <p className="home-hint">
          Accepts Steam64 · Steam32 · STEAM_0:x:y · [U:1:x] · steamcommunity.com/profiles/…
        </p>
      </div>
    </div>
  )
}
