import { useState } from 'react'

export default function App() {
  const [data, setData] = useState('')

  async function ping() {
    const res = await fetch('/api/v1/health')
    setData(JSON.stringify(await res.json(), null, 2))
  }

  return (
    <div>
      <h1>trophy-collector</h1>
      <button onClick={ping}>ping backend</button>
      <pre>{data}</pre>
    </div>
  )
}
