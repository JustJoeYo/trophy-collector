import {useMemo, useState } from 'react'

export type Player = {
    id: string
    name: string
    hero: string
}

//Sample set for testing before calling backend API. Not sure if we want to keep id field as string or use steamID, or something similar.

const SAMPLE_PLAYERS: Player[] = [
    { id: '1', name: 'AceRunner', hero: 'Seven' },
    { id: '2', name: "Nova", hero: 'Lady Geist' },
    { id: '3', name: "Tyler", hero: 'Bebop' },
]

export default function PlayerSearch() {
    const [query, setQuery] = useState('')

    const filteredPlayers = useMemo(() => {
        const normalizedQuery = query.trim().toLowerCase()

        if (!normalizedQuery) {
            return SAMPLE_PLAYERS
        }

        return SAMPLE_PLAYERS.filter((player) => {
            const searchableText = `${player.name} ${player.hero}`.toLowerCase()
            return searchableText.includes(normalizedQuery)
        })

    }, [query])

    return (
        <section>
            <label htmlFor="player-search">Search Players</label>
            <input
                id="player-search"
                type="text"
                value={query}
                onChange={(event) => setQuery(event.target.value)}
                placeholder="Search by player name or hero"
            />

            <ul>
                {filteredPlayers.map((player) => (
                    <li key={player.id}>
                        {player.name} - {player.hero}
                    </li>
                ))}
            </ul>
        </section>
    )
}