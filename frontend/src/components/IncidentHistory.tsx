import { formatDateTime } from '../format'
import { actionLabels } from '../labels'
import type { AuditEntry } from '../types'
import type { ApiState } from '../useApiData'

export default function IncidentHistory({ state }: { state: ApiState<AuditEntry[]> }) {
  return (
    <section className="panel">
      <h2>History</h2>
      {state.kind === 'loading' && <p className="panel-note">Loading history…</p>}
      {state.kind === 'failed' && <p className="panel-error">{state.message}</p>}
      {state.kind === 'ready' && (
        <ol className="timeline">
          {[...state.data].reverse().map((entry) => (
            <li key={entry.id}>
              <span className="timeline-action">{actionLabels[entry.action] ?? entry.action}</span>
              <span className="timeline-meta">
                {formatDateTime(entry.at)} · {entry.actor}
              </span>
              <span className="timeline-detail">{entry.detail}</span>
            </li>
          ))}
        </ol>
      )}
    </section>
  )
}
