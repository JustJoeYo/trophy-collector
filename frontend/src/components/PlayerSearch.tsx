import {useMemo, useState } from 'react'
import { Link } from 'react-router-dom'

type HeroStat = {
    name: string
    wins: number
    losses: number
}

type MatchStat = {
    hero: string
    outcome: 'Win' | 'Loss'
    kills: number
    deaths: number
    assists: number
}


export type Player = {
    id: string
    steamID: string
    name: string
    wins: number
    losses: number
    heroes: HeroStat[]
    matches: MatchStat[]
}

const SAMPLE_PLAYERS: Player[] = [
    {
        id: '1',
        steamID: '76561198000000001',
        name: 'Test Player',
        wins: 42,
        losses: 20,
        heroes: [
            { name: 'Infernus', wins: 10, losses: 4 },
            { name: 'McGinnis', wins: 8, losses: 6 },
        ],
        matches: [
            { hero: 'Infernus', outcome: 'Win', kills: 15, deaths: 5, assists: 8 },
            { hero: 'McGinnis', outcome: 'Loss', kills: 5, deaths: 9, assists: 11 },
        ],
    },
    {

        id: '2',
        steamID: '76561198000000002',
        name: 'Nova',
        wins: 31,
        losses: 27,
        heroes: [{ name: 'Seven', wins: 14, losses: 9 }],
        matches: [{ hero: 'Seven', outcome: 'Win', kills: 10, deaths: 4, assists: 7 }],
    },
    {
        id: '3',
        steamID: '76561198000000003',
        name: 'Tyler',
        wins: 28,
        losses: 30,
        heroes: [{ name: 'Bebop', wins: 11, losses: 12 }],
        matches: [{ hero: 'Bebop', outcome: 'Loss', kills: 6, deaths: 8, assists: 5 }],
    },
]

function PlayerMatchesQuery(player: Player, query: string) {
    const normalizedQuery = query.toLowerCase()

    if (!normalizedQuery) {
        return true
    }

    const searchableText = [
        player.id,
        player.steamID,
        player.name,
        String(player.wins),
        String(player.losses),
        ...player.heroes.flatMap((hero) => [hero.name, String(hero.wins), String(hero.losses)]),
        ...player.matches.flatMap((match) => [
            match.hero,
            match.outcome,
            String(match.kills),
            String(match.deaths),
            String(match.assists),
        ]),
    ]
        .join(' ')
        .toLowerCase()

return searchableText.includes(normalizedQuery)
}

export default function PlayerSearch() {
    const [query, setQuery] = useState('')

    const filteredPlayers = useMemo(() => {
        return SAMPLE_PLAYERS.filter((player) => PlayerMatchesQuery(player, query))
    }, [query])

    return (
        <section>
            <label htmlFor="player-search">Search players</label>
            <input
                id="player-search"
                type="text"
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder="Search by name, hero, stats..."
            />
            
            <ul>
                {filteredPlayers.map((player) => (
                    <li key={player.id}>
                        <strong>{player.name}</strong> - {player.steamID} - {player.wins}W / {player.losses}L
                        {' '}
                        <Link to={`/player/${player.id}`}>View profile</Link>
                    </li>
                ))}
            </ul>
        </section>
    )
}
