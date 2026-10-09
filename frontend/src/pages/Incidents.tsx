import { Link, useSearchParams } from 'react-router-dom'
import { Empty, LoadFailure, Loading } from '../components/LoadStates'
import PriorityBadge from '../components/PriorityBadge'
import { statusLabels, typeLabels } from '../labels'
import type { Incident } from '../types'
import { useApiResource } from '../useApiData'

export default function Incidents() {
  const [state, reload] = useApiResource<Incident[]>('/api/incidents')
  const [searchParams] = useSearchParams()
  const createdId = Number(searchParams.get('created'))

  if (state.kind === 'loading') {
    return (
      <section className="panel">
        <h2>Incidents</h2>
        <Loading what="incidents" />
      </section>
    )
  }

  if (state.kind === 'failed') {
    return (
      <section className="panel">
        <h2>Incidents</h2>
        <LoadFailure message={state.message} onRetry={reload} />
      </section>
    )
  }

  if (state.data.length === 0) {
    return (
      <section className="panel">
        <h2>Incidents</h2>
        <Empty>
          <p>Nothing has been registered yet.</p>
          <Link to="/incidents/new" className="button button-primary">
            Register incident
          </Link>
        </Empty>
      </section>
    )
  }

  return (
    <section className="panel">
      <div className="panel-header">
        <div>
          <h2>Incidents</h2>
          <p className="panel-note">
            {state.data.length} on the board, highest priority first.
          </p>
        </div>
        <Link to="/incidents/new" className="button button-primary">
          Register incident
        </Link>
      </div>

      <table className="table">
        <thead>
          <tr>
            <th scope="col">Priority</th>
            <th scope="col">Type</th>
            <th scope="col">Location</th>
            <th scope="col" className="numeric">People</th>
            <th scope="col">Status</th>
          </tr>
        </thead>
        <tbody>
          {state.data.map((incident) => (
            <tr key={incident.id} className={incident.id === createdId ? 'row-new' : undefined}>
              <td>
                <PriorityBadge band={incident.priority_band} score={incident.priority_score} />
              </td>
              <td>{typeLabels[incident.type]}</td>
              <td>
                <Link to={`/incidents/${incident.id}`}>{incident.location}</Link>
                {incident.vulnerable && <span className="tag">Vulnerable people</span>}
              </td>
              <td className="numeric">{incident.people_affected}</td>
              <td>
                <span className={`status status-${incident.status}`}>
                  {statusLabels[incident.status]}
                </span>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  )
}
