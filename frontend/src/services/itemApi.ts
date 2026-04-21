import type { Item, ItemStats } from '@/types'

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? ''

async function fetchJson<T>(path: string): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${path}`)

    if (!response.ok) {
        throw new Error(`Request failed with status ${response.status}`)
    }

    return response.json() as Promise<T>
}

export async function getItems(): Promise<Item[]> {
    return fetchJson<Item[]>('/api/v1/items')
}

export async function getItemStats(): Promise<ItemStats[]> {
    return fetchJson<ItemStats[]>('/api/v1/items/stats')
}
