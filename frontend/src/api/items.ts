import { fetchJson } from './client'
import type { Image, Item, ItemStats } from '@/types'

export const getItems = () =>
  fetchJson<Item[]>('/api/v1/items')

export const getItemStats = () =>
  fetchJson<ItemStats[]>('/api/v1/items/stats')

export const getImage = () =>
  fetchJson<Image[]>('/api/v1/images')