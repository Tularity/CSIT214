import type { Health } from '../types'
import { useApiData } from '../useApiData'

export default function HealthBadge() {
  const state = useApiData<Health>('/api/health')

  if (state.kind === 'loading') {
    return <span className="badge badge-waiting">Checking API</span>
  }
  if (state.kind === 'failed') {
    return (
      <span className="badge badge-down" title={state.message}>
        API unreachable
      </span>
    )
  }
  return <span className="badge badge-up">API {state.data.status}</span>
}
