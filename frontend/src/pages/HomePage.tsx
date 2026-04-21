import PlayerSearch from '../components/PlayerSearch'

export default function HomePage() {
  return (
    <main>
      <h1>Trophy Collector</h1>
      <p>Enter your Steam ID to view your Deadlock stats.</p>
      <PlayerSearch />
    </main>
  )
}
