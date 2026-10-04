import { useState } from 'react'
import { postJSON } from '../api'
import { useNotify } from '../notices'
import type { Incident, Responder, Shelter } from '../types'

type IncidentActionsProps = {
  incident: Incident
  shelters: Shelter[]
  responders: Responder[]
  onChanged: () => void
}

const responderStatusNotes: Record<Responder['status'], string> = {
  available: '',
  assigned: ' — on another job',
  off_duty: ' — off duty',
}

export default function IncidentActions({ incident, shelters, responders, onChanged }: IncidentActionsProps) {
  const notify = useNotify()
  const [busy, setBusy] = useState(false)
  const [shelterId, setShelterId] = useState('')
  const [responderId, setResponderId] = useState('')

  // Every shelter stays selectable: if one is short of room the backend says
  // exactly how many places are left, which is more useful than a greyed out row.
  const roomyShelter = shelters.find((shelter) => shelter.capacity_remaining >= incident.people_affected)
  const freeResponder = responders.find((responder) => responder.status === 'available')
  const chosenShelter = shelterId || String(roomyShelter?.id ?? shelters[0]?.id ?? '')
  const chosenResponder = responderId || String(freeResponder?.id ?? responders[0]?.id ?? '')

  async function run(action: string, success: string, payload?: unknown) {
    setBusy(true)
    try {
      await postJSON<Incident>(`/api/incidents/${incident.id}/${action}`, payload)
      setShelterId('')
      setResponderId('')
      notify('success', success)
      onChanged()
    } catch (failure) {
      notify('error', (failure as Error).message)
    } finally {
      setBusy(false)
    }
  }

  function nameOf<T extends { id: number; name: string }>(records: T[], id: string): string {
    return records.find((record) => String(record.id) === id)?.name ?? `#${id}`
  }

  const canAssess = incident.status === 'registered'
  const canPlace = incident.status === 'assessed'
  const canDispatch = incident.status === 'assessed' || incident.status === 'assigned'
  const canResolve = ['assessed', 'assigned', 'in_progress'].includes(incident.status)

  if (incident.status === 'resolved') {
    return (
      <section className="panel">
        <h2>Next step</h2>
        <p className="panel-note">This incident is closed. Its shelter places and responder have been released.</p>
      </section>
    )
  }

  return (
    <section className="panel">
      <h2>Next step</h2>

      <div className="actions">
        {canAssess && (
          <div className="action">
            <p>Confirm the details with the caller before anything is committed to it.</p>
            <button className="button button-primary" disabled={busy} onClick={() => run('assess', `Incident #${incident.id} assessed.`)}>
              Assess
            </button>
          </div>
        )}

        {canPlace && (
          <div className="action">
            <label className="field">
              <span className="field-label">Shelter for {incident.people_affected} people</span>
              <select value={chosenShelter} onChange={(event) => setShelterId(event.target.value)}>
                {shelters.map((shelter) => (
                  <option key={shelter.id} value={shelter.id}>
                    {shelter.name} — {shelter.capacity_remaining} places left
                  </option>
                ))}
              </select>
            </label>
            <button
              className="button button-primary"
              disabled={busy || chosenShelter === ''}
              onClick={() =>
                run(
                  'assign-shelter',
                  `${incident.people_affected} people placed at ${nameOf(shelters, chosenShelter)}.`,
                  { shelter_id: Number(chosenShelter) },
                )
              }
            >
              Assign shelter
            </button>
          </div>
        )}

        {canDispatch && (
          <div className="action">
            <label className="field">
              <span className="field-label">Responder</span>
              <select value={chosenResponder} onChange={(event) => setResponderId(event.target.value)}>
                {responders.map((responder) => (
                  <option key={responder.id} value={responder.id}>
                    {responder.name} ({responder.skill}){responderStatusNotes[responder.status]}
                  </option>
                ))}
              </select>
            </label>
            <button
              className={canPlace ? 'button' : 'button button-primary'}
              disabled={busy || chosenResponder === ''}
              onClick={() =>
                run('assign-responder', `${nameOf(responders, chosenResponder)} dispatched.`, {
                  responder_id: Number(chosenResponder),
                })
              }
            >
              Dispatch responder
            </button>
          </div>
        )}

        {canResolve && (
          <div className="action">
            <p>Closing the incident gives back its shelter places and frees its responder.</p>
            <button
              className={incident.status === 'in_progress' ? 'button button-primary' : 'button'}
              disabled={busy}
              onClick={() => run('resolve', `Incident #${incident.id} resolved; its shelter places and responder are free again.`)}
            >
              Resolve
            </button>
          </div>
        )}
      </div>
    </section>
  )
}
