import { Link } from 'react-router-dom'
import { Empty, LoadFailure, Loading } from '../components/LoadStates'
import { formatDateTime } from '../format'
import { actionLabels } from '../labels'
import type { AuditEntry } from '../types'
import { useApiResource } from '../useApiData'

function recordLink(entry: AuditEntry) {
  if (entry.entity_type === 'incident') {
    return <Link to={`/incidents/${entry.entity_id}`}>Incident #{entry.entity_id}</Link>
  }
  return `${entry.entity_type} #${entry.entity_id}`
}

export default function ActivityLog() {
  const [state, reload] = useApiResource<AuditEntry[]>('/api/audit')

  return (
    <section className="panel">
      <h2>Activity log</h2>
      <p className="panel-note">
        Every change made through the system, newest first. Entries cannot be edited or removed.
      </p>

      {state.kind === 'loading' && <Loading what="the log" />}
      {state.kind === 'failed' && <LoadFailure message={state.message} onRetry={reload} />}
      {state.kind === 'ready' && state.data.length === 0 && (
        <Empty>
          <p>Nothing has been recorded yet.</p>
        </Empty>
      )}

      {state.kind === 'ready' && state.data.length > 0 && (
        <table className="table">
          <thead>
            <tr>
              <th scope="col">When</th>
              <th scope="col">Who</th>
              <th scope="col">Action</th>
              <th scope="col">Record</th>
              <th scope="col">Detail</th>
            </tr>
          </thead>
          <tbody>
            {state.data.map((entry) => (
              <tr key={entry.id}>
                <td className="nowrap">{formatDateTime(entry.at)}</td>
                <td>{entry.actor}</td>
                <td className="nowrap">{actionLabels[entry.action] ?? entry.action}</td>
                <td className="nowrap">{recordLink(entry)}</td>
                <td>{entry.detail}</td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
