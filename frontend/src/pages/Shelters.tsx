import OccupancyMeter from '../components/OccupancyMeter'
import type { Shelter } from '../types'
import { useApiData } from '../useApiData'

function percentUsed(shelter: Shelter): number {
  return shelter.capacity_total === 0 ? 0 : Math.round((shelter.capacity_used / shelter.capacity_total) * 100)
}

export default function Shelters() {
  const state = useApiData<Shelter[]>('/api/shelters')

  return (
    <section className="panel">
      <h2>Shelters</h2>
      <p className="panel-note">Places are reserved when an incident is placed and given back when it is resolved.</p>

      {state.kind === 'loading' && <p className="panel-note">Loading shelters…</p>}
      {state.kind === 'failed' && <p className="panel-error">{state.message}</p>}
      {state.kind === 'ready' && state.data.length === 0 && (
        <p className="panel-note">No shelters have been set up.</p>
      )}

      {state.kind === 'ready' && state.data.length > 0 && (
        <table className="table">
          <thead>
            <tr>
              <th scope="col">Shelter</th>
              <th scope="col" className="numeric">Capacity</th>
              <th scope="col" className="numeric">In use</th>
              <th scope="col" className="numeric">Remaining</th>
              <th scope="col" className="meter-column">Occupancy</th>
            </tr>
          </thead>
          <tbody>
            {state.data.map((shelter) => (
              <tr key={shelter.id}>
                <td>
                  <strong>{shelter.name}</strong>
                  <div className="panel-note">{shelter.address}</div>
                </td>
                <td className="numeric">{shelter.capacity_total}</td>
                <td className="numeric">{shelter.capacity_used}</td>
                <td className="numeric">{shelter.capacity_remaining}</td>
                <td className="meter-column">
                  <OccupancyMeter
                    name={shelter.name}
                    used={shelter.capacity_used}
                    total={shelter.capacity_total}
                    percent={percentUsed(shelter)}
                    compact
                  />
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </section>
  )
}
