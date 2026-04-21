import { useMemo } from 'react'
import { useParams, Link } from 'react-router-dom'
import { SAMPLE_HEROES, SAMPLE_PLAYERS, SAMPLE_MATCHES } from '../components/PlayerSearch'


    function findPlayerByAccountId(accountId: string | undefined) {
        if (!accountId) {
            return null
        }

        return SAMPLE_PLAYERS.find((player) => player.account_id === accountId) || null
    }

    function getHeroName(heroId: number): string {
        const hero = SAMPLE_HEROES.find((currentHero) => currentHero.hero_id === heroId)
        return hero?.name || `Hero ${heroId}`
    }

    function formatDuration(durationSecs: number): string {
        const minutes = Math.floor(durationSecs / 60)
        const seconds = durationSecs % 60
        return `${minutes}:${seconds.toString().padStart(2, '0')}`
    }

    export default function PlayerDetailPage() {
        const { id } = useParams()
        const player = findPlayerByAccountId(id)

        const playerMatches = useMemo(() => {
            if (!player) {
                return []
            }
            return SAMPLE_MATCHES[player.account_id] || []
        }, [player])

        if (!player) {
            return (
                <main>
                    <h1>Player not found</h1>
                    <p>No player exists for ID {id}.</p>
                    <Link to="/">Back to Home</Link>
                </main>
            )
        }

        return (
            <main>
                <h1>{player.player_name}</h1>
                <p>Account ID: {player.account_id}</p>
                <p>Record: {player.wins}W / {player.losses}L</p>

                <section>
                    <h2>Top Heroes</h2>
                    <table>
                        <thead>
                            <tr>
                                <th>Hero</th>
                                <th>Wins</th>
                                <th>Losses</th>
                            </tr>
                        </thead>
                        <tbody>
                            {SAMPLE_HEROES.map((hero) => (
                                <tr key={hero.hero_id}>
                                    <td>{hero.name}</td>
                                    <td>--</td>
                                    <td>--</td>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                </section>

                <section>
                    <h2>Recent Matches</h2>
                    {playerMatches.length > 0 ? (
                        <table>
                            <thead>
                                <tr>
                                    <th>Hero</th>
                                    <th>Outcome</th>
                                    <th>Kills</th>
                                    <th>Deaths</th>
                                    <th>Assists</th>
                                    <th>Duration</th>
                                </tr>
                            </thead>
                            <tbody>
                                {playerMatches.map((match) => (
                                    <tr key={match.match_id}>
                                        <td>{getHeroName(match.hero_id)}</td>
                                        <td>{match.outcome}</td>
                                        <td>{match.kills}</td>
                                        <td>{match.deaths}</td>
                                        <td>{match.assists}</td>
                                        <td>{formatDuration(match.duration_secs)}</td>
                                    </tr>
                                ))}
                            </tbody>
                        </table>
                    ) : (
                        <p>No recent matches found.</p>
                    )}
                </section>

                <p>
                    <Link to="/">Back to Home</Link>
                </p>
            </main>
        )
    }