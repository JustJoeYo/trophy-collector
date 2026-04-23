export const API_BASE = window.location.hostname === 'localhost'
  ? 'http://localhost:8080'
  : 'https://trophy-collector-backend.onrender.com'

export async function fetchJson<T>(path: string): Promise<T> {
  const res = await fetch(`${API_BASE}${path}`)
  if (!res.ok) throw new Error(`${res.status} ${res.statusText}`)
  return res.json() as Promise<T>
}
