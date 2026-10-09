import { Empty, LoadFailure, Loading } from '../components/LoadStates'
import type { Responder } from '../types'
import { useApiResource } from '../useApiData'

const statusLabels: Record<Responder['status'], string> = {
  available: 'Available',
  assigned: 'On a job',
  off_duty: 'Off duty',
}

export default function Responders() {
  const [state, reload] = useApiResource<Responder[]>('/api/responders')

  return (
    <section className="panel">
      <h2>Responders</h2>
      <p className="panel-note">
        A responder is marked as on a job when dispatched and becomes available again when the incident is
        resolved.
      </p>

      {state.kind === 'loading' && <Loading what="responders" />}
      {state.kind === 'failed' && <LoadFailure message={state.message} onRetry={reload} />}
      {state.kind === 'ready' && state.data.length === 0 && (
        <Empty>
          <p>No responders are on the roster.</p>
        </Empty>
      )}

      {state.kind === 'ready' && state.data.length > 0 && (
        <table className="table">
          <thead>
            <tr>
              <th scope="col">Name</th>
              <th scope="col">Skill</th>
              <th scope="col">Status</th>
            </tr>
          </thead>
          <tbody>
            {state.data.map((responder) => (
              <tr key={responder.id}>
                <td>{responder.name}</td>
                <td>{responder.skill}</td>
                <td>
                  <span className={`status status-responder-${responder.status}`}>
                    {statusLabels[responder.status]}
                  </span>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
