import { fetchJson } from './client'
import type { Item, ItemStats } from '@/types'

export const getItems = () =>
  fetchJson<Item[]>('/api/v1/items')

export const getItemStats = () =>
  fetchJson<ItemStats[]>('/api/v1/items/stats')
