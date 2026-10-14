import type { ShelterAlternative } from '../types'

type ShelterAlternativesProps = {
  people: number
  refusedName: string
  alternatives: ShelterAlternative[]
  busy: boolean
  onChoose: (alternative: ShelterAlternative) => void
}

export default function ShelterAlternatives({ people, refusedName, alternatives, busy, onChoose }: ShelterAlternativesProps) {
  if (alternatives.length === 0) {
    return (
      <div className="alternatives alternatives-none" role="region" aria-label="Shelters with room">
        <p className="alternatives-title">No shelter has room for all {people} people at once.</p>
        <p className="panel-note">Split the group into separate incidents, or free places by resolving others first.</p>
      </div>
    )
  }

  return (
    <div className="alternatives" role="region" aria-label="Shelters with room">
      <p className="alternatives-title">
        {refusedName} cannot take {people} people. These shelters can:
      </p>
      <ul className="alternatives-list">
        {alternatives.map((alternative) => (
          <li key={alternative.id}>
            <span>
              <strong>{alternative.name}</strong>
              <span className="panel-note"> — {alternative.capacity_remaining} places left</span>
            </span>
            <button type="button" className="button button-primary" disabled={busy} onClick={() => onChoose(alternative)}>
              Place here
            </button>
          </li>
        ))}
      </ul>
    </div>
  )
}
