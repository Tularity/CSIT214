import { useState } from 'react'
import { ApiError, postJSON } from '../api'
import { useNotify } from '../notices'
import type { Incident, Responder, Shelter, ShelterAlternative } from '../types'
import ShelterAlternatives from './ShelterAlternatives'

type IncidentActionsProps = {
  incident: Incident
  shelters: Shelter[]
  responders: Responder[]
  onChanged: () => void
}

function alternativesIn(body: unknown): ShelterAlternative[] | null {
  if (body && typeof body === 'object' && 'alternatives' in body) {
    const alternatives = (body as { alternatives: unknown }).alternatives
    if (Array.isArray(alternatives)) return alternatives as ShelterAlternative[]
  }
  return null
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
  const [refused, setRefused] = useState<{ name: string; message: string; alternatives: ShelterAlternative[] } | null>(null)

  // Every shelter stays selectable: if one is short of room the backend says
  // exactly how many places are left, which is more useful than a greyed out row.
  const roomyShelter = shelters.find((shelter) => shelter.capacity_remaining >= incident.people_affected)
  const freeResponder = responders.find((responder) => responder.status === 'available')
  const chosenShelter = shelterId || String(roomyShelter?.id ?? shelters[0]?.id ?? '')
  const chosenResponder = responderId || String(freeResponder?.id ?? responders[0]?.id ?? '')

  async function run(action: string, success: string, payload?: unknown): Promise<unknown> {
    setBusy(true)
    try {
      await postJSON<Incident>(`/api/incidents/${incident.id}/${action}`, payload)
      setShelterId('')
      setResponderId('')
      setRefused(null)
      notify('success', success)
      onChanged()
      return null
    } catch (failure) {
      return failure
    } finally {
      setBusy(false)
    }
  }

  async function runOrReport(action: string, success: string, payload?: unknown) {
    const failure = await run(action, success, payload)
    if (failure) notify('error', (failure as Error).message)
  }

  // A refusal that comes with alternatives is shown in the panel beside the
  // shelter choice rather than as a notice: it is a question to answer, and a
  // notice would linger after the people had been placed elsewhere.
  async function placeIn(id: string, name: string) {
    const failure = await run('assign-shelter', `${incident.people_affected} people placed at ${name}.`, {
      shelter_id: Number(id),
    })
    if (!failure) return

    const alternatives = failure instanceof ApiError ? alternativesIn(failure.body) : null
    if (alternatives) {
      setRefused({ name, message: (failure as Error).message, alternatives })
    } else {
      setRefused(null)
      notify('error', (failure as Error).message)
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
            <button className="button button-primary" disabled={busy} onClick={() => runOrReport('assess', `Incident #${incident.id} assessed.`)}>
              Assess
            </button>
          </div>
        )}

        {canPlace && (
          <div className="action">
            <label className="field">
              <span className="field-label">Shelter for {incident.people_affected} people</span>
              <select
                value={chosenShelter}
                onChange={(event) => {
                  setShelterId(event.target.value)
                  setRefused(null)
                }}
              >
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
              onClick={() => placeIn(chosenShelter, nameOf(shelters, chosenShelter))}
            >
              Assign shelter
            </button>
          </div>
        )}

        {canPlace && refused && (
          <ShelterAlternatives
            people={incident.people_affected}
            refusedName={refused.name}
            message={refused.message}
            alternatives={refused.alternatives}
            busy={busy}
            onChoose={(alternative) => placeIn(String(alternative.id), alternative.name)}
          />
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
                runOrReport('assign-responder', `${nameOf(responders, chosenResponder)} dispatched.`, {
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
              onClick={() => runOrReport('resolve', `Incident #${incident.id} resolved; its shelter places and responder are free again.`)}
            >
              Resolve
            </button>
          </div>
        )}
      </div>
    </section>
  )
}
