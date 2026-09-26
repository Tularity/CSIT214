import { useEffect, useState } from 'react'
import { getJSON } from './api'

export type ApiState<T> =
  | { kind: 'loading' }
  | { kind: 'ready'; data: T }
  | { kind: 'failed'; message: string }

export function useApiData<T>(path: string): ApiState<T> {
  const [state, setState] = useState<ApiState<T>>({ kind: 'loading' })

  useEffect(() => {
    let active = true

    getJSON<T>(path)
      .then((data) => {
        if (active) setState({ kind: 'ready', data })
      })
      .catch((error: Error) => {
        if (active) setState({ kind: 'failed', message: error.message })
      })

    return () => {
      active = false
    }
  }, [path])

  return state
}
