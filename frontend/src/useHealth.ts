import { useEffect, useState } from 'react'
import { getJSON } from './api'

export type Health = { status: string }

export type HealthState =
  | { kind: 'loading' }
  | { kind: 'ready'; health: Health }
  | { kind: 'failed'; message: string }

export function useHealth(): HealthState {
  const [state, setState] = useState<HealthState>({ kind: 'loading' })

  useEffect(() => {
    let active = true

    getJSON<Health>('/api/health')
      .then((health) => {
        if (active) setState({ kind: 'ready', health })
      })
      .catch((error: Error) => {
        if (active) setState({ kind: 'failed', message: error.message })
      })

    return () => {
      active = false
    }
  }, [])

  return state
}
