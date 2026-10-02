import { useCallback, useEffect, useState } from 'react'
import { getJSON } from './api'

export type ApiState<T> =
  | { kind: 'loading' }
  | { kind: 'ready'; data: T }
  | { kind: 'failed'; message: string }

// While reloading, the previous data stays on screen instead of flashing back
// to a loading state after every action.
export function useApiResource<T>(path: string): [ApiState<T>, () => void] {
  const [state, setState] = useState<ApiState<T>>({ kind: 'loading' })
  const [version, setVersion] = useState(0)

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
  }, [path, version])

  const reload = useCallback(() => setVersion((current) => current + 1), [])
  return [state, reload]
}

export function useApiData<T>(path: string): ApiState<T> {
  return useApiResource<T>(path)[0]
}
