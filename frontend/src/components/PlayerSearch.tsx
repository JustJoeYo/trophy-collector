import {useMemo, useState } from 'react'
import { Link } from 'react-router-dom'
import { Player, Match, Hero } from '@/types'

const SAMPLE_HEROES: Hero[] = [
    { hero_id: 1, name: 'Infernus', image_url: '' },
    { hero_id: 2, name: 'McGinnis', image_url: '' },
    { hero_id: 3, name: 'Seven', image_url: '' },
    { hero_id: 4, name: 'Bebop', image_url: 'example.com/bebop.jpg' },
]

const SAMPLE_PLAYERS: Player[] = [
    {
        account_id: '76561198000000001',
        player_name: 'Test Player',
        avatar_url: '',
        wins: 42,
        losses: 20,
    },
    {

        account_id: '76561198000000002',
        player_name: 'Nova',
        avatar_url: '',
        wins: 31,
        losses: 27,
    },
    {
        account_id: '76561198000000003',
        player_name: 'Tyler',
        avatar_url: '',
        wins: 28,
        losses: 30,
    },
]

const SAMPLE_MATCHES: Record<string, Match[]> = {
    '76561198000000001': [
        { match_id: 1, hero_id: 1, outcome: 'Win', kills: 12, deaths: 3, assists: 8, duration_secs: 2400 },
        { match_id: 2, hero_id: 2, outcome: 'Loss', kills: 5, deaths: 9, assists: 11, duration_secs: 1800 },
    ],
    '76561198000000002': [
        { match_id: 3, hero_id: 3, outcome: 'Win', kills: 10, deaths: 4, assists: 7, duration_secs: 2350 },
    ],
    '76561198000000003': [
        { match_id: 4, hero_id: 4, outcome: 'Loss', kills: 6, deaths: 8, assists: 5, duration_secs: 1950 },
    ],
}

function PlayerMatchesQuery(player: Player, query: string) {
    const normalizedQuery = query.toLowerCase()

    if (!normalizedQuery) {
        return true
    }

    const searchableText = [
        player.account_id,
        player.player_name,
        String(player.wins),
        String(player.losses),
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
                placeholder="Search by SteamID or Player name"
            />
            
            <ul>
                {filteredPlayers.map((player) => (
                    <li key={player.account_id}>
                        <strong>{player.player_name}</strong> - {player.account_id} - {player.wins}W / {player.losses}L
                        {' '}
                        <Link to={`/player/${player.account_id}`}>View profile</Link>
                    </li>
                ))}
            </ul>
        </section>
    )
}

export { SAMPLE_HEROES, SAMPLE_MATCHES, SAMPLE_PLAYERS }