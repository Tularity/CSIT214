import { Link, useParams } from 'react-router-dom'
import PriorityBadge from '../components/PriorityBadge'
import { statusLabels, typeLabels } from '../labels'
import type { Incident, Responder, Shelter } from '../types'
import { useApiResource } from '../useApiData'

function formatReportedAt(iso: string): string {
  return new Date(iso).toLocaleString(undefined, {
    dateStyle: 'medium',
    timeStyle: 'short',
  })
}

export default function IncidentDetail() {
  const { id } = useParams()
  const [incidentState] = useApiResource<Incident>(`/api/incidents/${id}`)
  const [shelterState] = useApiResource<Shelter[]>('/api/shelters')
  const [responderState] = useApiResource<Responder[]>('/api/responders')

  if (incidentState.kind === 'loading') {
    return (
      <section className="panel">
        <p className="panel-note">Loading incident…</p>
      </section>
    )
  }

  if (incidentState.kind === 'failed') {
    return (
      <section className="panel">
        <h2>Incident {id}</h2>
        <p className="panel-error">{incidentState.message}</p>
        <Link to="/incidents">Back to incidents</Link>
      </section>
    )
  }

  const incident = incidentState.data
  const shelters = shelterState.kind === 'ready' ? shelterState.data : []
  const responders = responderState.kind === 'ready' ? responderState.data : []
  const shelter = shelters.find((candidate) => candidate.id === incident.shelter_id)
  const responder = responders.find((candidate) => candidate.id === incident.responder_id)

  return (
    <>
      <p className="breadcrumb">
        <Link to="/incidents">Incidents</Link> / #{incident.id}
      </p>

      <section className="panel">
        <div className="panel-header">
          <div>
            <h2>
              {typeLabels[incident.type]} at {incident.location}
            </h2>
            <p className="panel-note">Reported {formatReportedAt(incident.reported_at)}</p>
          </div>
          <span className={`status status-${incident.status}`}>{statusLabels[incident.status]}</span>
        </div>

        <dl className="facts">
          <div>
            <dt>Priority</dt>
            <dd>
              <PriorityBadge band={incident.priority_band} score={incident.priority_score} />
            </dd>
          </div>
          <div>
            <dt>Severity</dt>
            <dd>{incident.severity} of 5</dd>
          </div>
          <div>
            <dt>People affected</dt>
            <dd>{incident.people_affected}</dd>
          </div>
          <div>
            <dt>Vulnerable people</dt>
            <dd>{incident.vulnerable ? 'Yes' : 'No'}</dd>
          </div>
          <div>
            <dt>Shelter</dt>
            <dd>{shelter ? shelter.name : incident.shelter_id === null ? 'Not assigned' : `#${incident.shelter_id}`}</dd>
          </div>
          <div>
            <dt>Responder</dt>
            <dd>
              {responder
                ? `${responder.name} (${responder.skill})`
                : incident.responder_id === null
                  ? 'Not dispatched'
                  : `#${incident.responder_id}`}
            </dd>
          </div>
        </dl>

        <h3 className="subheading">Description</h3>
        <p>{incident.description || <span className="panel-note">No description given.</span>}</p>
      </section>
    </>
  )
}
