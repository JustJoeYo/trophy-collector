import { useQuery } from '@tanstack/react-query'
import { getPlayer, getPlayerMatches, getPlayerHeroes } from '../api/client'

export function usePlayer(steamId: string) {
  return useQuery({
    queryKey: ['player', steamId],
    queryFn: () => getPlayer(steamId),
    enabled: !!steamId,
  })
}

export function usePlayerMatches(steamId: string) {
  return useQuery({
    queryKey: ['player', steamId, 'matches'],
    queryFn: () => getPlayerMatches(steamId),
    enabled: !!steamId,
  })
}

export function usePlayerHeroes(steamId: string) {
  return useQuery({
    queryKey: ['player', steamId, 'heroes'],
    queryFn: () => getPlayerHeroes(steamId),
    enabled: !!steamId,
  })
}
