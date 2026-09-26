import PriorityBadge from '../components/PriorityBadge'
import type { Incident } from '../types'
import { useApiData } from '../useApiData'

const typeLabels: Record<Incident['type'], string> = {
  flood: 'Flood',
  fire: 'Fire',
  storm: 'Storm',
  medical: 'Medical',
  other: 'Other',
}

const statusLabels: Record<Incident['status'], string> = {
  registered: 'Registered',
  assessed: 'Assessed',
  assigned: 'Shelter assigned',
  in_progress: 'In progress',
  resolved: 'Resolved',
}

export default function Incidents() {
  const state = useApiData<Incident[]>('/api/incidents')

  if (state.kind === 'loading') {
    return (
      <section className="panel">
        <h2>Incidents</h2>
        <p className="panel-note">Loading incidents…</p>
      </section>
    )
  }

  if (state.kind === 'failed') {
    return (
      <section className="panel">
        <h2>Incidents</h2>
        <p className="panel-error">{state.message}</p>
      </section>
    )
  }

  if (state.data.length === 0) {
    return (
      <section className="panel">
        <h2>Incidents</h2>
        <p className="panel-note">Nothing has been registered yet.</p>
      </section>
    )
  }

  return (
    <section className="panel">
      <h2>Incidents</h2>
      <p className="panel-note">
        {state.data.length} on the board, highest priority first.
      </p>

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
            <tr key={incident.id}>
              <td>
                <PriorityBadge band={incident.priority_band} score={incident.priority_score} />
              </td>
              <td>{typeLabels[incident.type]}</td>
              <td>
                {incident.location}
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
