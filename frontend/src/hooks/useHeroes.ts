import { useQuery } from '@tanstack/react-query'
import { getHeroes, getLeaderboard } from '../api/client'

export function useHeroes() {
  return useQuery({
    queryKey: ['heroes'],
    queryFn: getHeroes,
    staleTime: 60 * 60 * 1000, // Hero list is stable — cache for 1 hour
  })
}

export function useLeaderboard() {
  return useQuery({
    queryKey: ['leaderboard'],
    queryFn: getLeaderboard,
    staleTime: 10 * 60 * 1000, // Cache for 10 minutes
  })
}
