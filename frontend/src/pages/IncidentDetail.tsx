import { Link, useParams } from 'react-router-dom'
import IncidentActions from '../components/IncidentActions'
import { LoadFailure, Loading } from '../components/LoadStates'
import IncidentHistory from '../components/IncidentHistory'
import PriorityBadge from '../components/PriorityBadge'
import { formatDateTime } from '../format'
import { statusLabels, typeLabels } from '../labels'
import type { AuditEntry, Incident, Responder, Shelter } from '../types'
import { useApiResource } from '../useApiData'

export default function IncidentDetail() {
  const { id } = useParams()
  const [incidentState, reloadIncident] = useApiResource<Incident>(`/api/incidents/${id}`)
  const [shelterState, reloadShelters] = useApiResource<Shelter[]>('/api/shelters')
  const [responderState, reloadResponders] = useApiResource<Responder[]>('/api/responders')
  const [historyState, reloadHistory] = useApiResource<AuditEntry[]>(
    `/api/audit?entity_type=incident&entity_id=${id}`,
  )

  function refreshAll() {
    reloadIncident()
    reloadShelters()
    reloadResponders()
    reloadHistory()
  }

  if (incidentState.kind === 'loading') {
    return (
      <section className="panel">
        <Loading what="the incident" />
      </section>
    )
  }

  if (incidentState.kind === 'failed') {
    return (
      <section className="panel">
        <h2>Incident {id}</h2>
        <LoadFailure message={incidentState.message} onRetry={refreshAll} />
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
            <p className="panel-note">Reported {formatDateTime(incident.reported_at)}</p>
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

      <IncidentActions
        incident={incident}
        shelters={shelters}
        responders={responders}
        onChanged={refreshAll}
      />

      <IncidentHistory state={historyState} />
    </>
  )
}
