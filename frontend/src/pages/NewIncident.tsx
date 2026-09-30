import { useState, type FormEvent } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { postJSON } from '../api'
import { typeLabels } from '../labels'
import type { Incident, IncidentType } from '../types'

type FormValues = {
  type: IncidentType
  location: string
  description: string
  peopleAffected: string
  severity: string
  vulnerable: boolean
}

const initialValues: FormValues = {
  type: 'flood',
  location: '',
  description: '',
  peopleAffected: '',
  severity: '3',
  vulnerable: false,
}

const severityHints: Record<string, string> = {
  '1': '1 — Minor',
  '2': '2 — Moderate',
  '3': '3 — Serious',
  '4': '4 — Severe',
  '5': '5 — Life-threatening',
}

export default function NewIncident() {
  const navigate = useNavigate()
  const [values, setValues] = useState<FormValues>(initialValues)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState<string | null>(null)

  function update<K extends keyof FormValues>(field: K, value: FormValues[K]) {
    setValues((current) => ({ ...current, [field]: value }))
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    setSubmitting(true)
    setError(null)

    try {
      const created = await postJSON<Incident>('/api/incidents', {
        type: values.type,
        location: values.location,
        description: values.description,
        people_affected: Number(values.peopleAffected),
        severity: Number(values.severity),
        vulnerable: values.vulnerable,
      })
      navigate(`/incidents?created=${created.id}`)
    } catch (failure) {
      setError((failure as Error).message)
      setSubmitting(false)
    }
  }

  return (
    <section className="panel">
      <h2>Register an incident</h2>
      <p className="panel-note">
        The priority is worked out by the system from severity, people affected and whether
        vulnerable people are involved.
      </p>

      {/* noValidate: the server's validation wording is the one the operator should see. */}
      <form className="form" onSubmit={submit} noValidate>
        <label className="field">
          <span className="field-label">Type</span>
          <select value={values.type} onChange={(event) => update('type', event.target.value as IncidentType)}>
            {Object.entries(typeLabels).map(([value, label]) => (
              <option key={value} value={value}>
                {label}
              </option>
            ))}
          </select>
        </label>

        <label className="field">
          <span className="field-label">Location</span>
          <input
            type="text"
            value={values.location}
            placeholder="Street address or landmark"
            onChange={(event) => update('location', event.target.value)}
          />
        </label>

        <label className="field">
          <span className="field-label">Description</span>
          <textarea
            rows={3}
            value={values.description}
            placeholder="What has happened and what is needed"
            onChange={(event) => update('description', event.target.value)}
          />
        </label>

        <div className="field-row">
          <label className="field">
            <span className="field-label">People affected</span>
            <input
              type="number"
              inputMode="numeric"
              value={values.peopleAffected}
              onChange={(event) => update('peopleAffected', event.target.value)}
            />
          </label>

          <label className="field">
            <span className="field-label">Severity</span>
            <select value={values.severity} onChange={(event) => update('severity', event.target.value)}>
              {Object.entries(severityHints).map(([value, label]) => (
                <option key={value} value={value}>
                  {label}
                </option>
              ))}
            </select>
          </label>
        </div>

        <label className="field-check">
          <input
            type="checkbox"
            checked={values.vulnerable}
            onChange={(event) => update('vulnerable', event.target.checked)}
          />
          Elderly people, children or people with reduced mobility are involved
        </label>

        {error && <p className="panel-error">{error}</p>}

        <div className="form-actions">
          <button type="submit" className="button button-primary" disabled={submitting}>
            {submitting ? 'Registering…' : 'Register incident'}
          </button>
          <Link to="/incidents" className="button">
            Cancel
          </Link>
        </div>
      </form>
    </section>
  )
}
