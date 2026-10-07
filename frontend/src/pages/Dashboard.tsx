import { Link } from 'react-router-dom'
import OccupancyMeter from '../components/OccupancyMeter'
import { statusLabels } from '../labels'
import type { Dashboard as DashboardData, IncidentStatus, PriorityBand } from '../types'
import { useApiData } from '../useApiData'

const statusOrder: IncidentStatus[] = ['registered', 'assessed', 'assigned', 'in_progress', 'resolved']
const bandOrder: PriorityBand[] = ['Critical', 'High', 'Medium', 'Low']

export default function Dashboard() {
  const state = useApiData<DashboardData>('/api/dashboard')

  if (state.kind === 'loading') {
    return (
      <section className="panel">
        <h2>Operations overview</h2>
        <p className="panel-note">Loading the latest figures…</p>
      </section>
    )
  }

  if (state.kind === 'failed') {
    return (
      <section className="panel">
        <h2>Operations overview</h2>
        <p className="panel-error">{state.message}</p>
      </section>
    )
  }

  const data = state.data
  const respondersTotal = Object.values(data.responders_by_status).reduce((sum, count) => sum + count, 0)

  return (
    <>
      <section className="kpis" aria-label="Headline figures">
        <div className="kpi kpi-hero">
          <span className="kpi-label">Open incidents</span>
          <span className="kpi-value">{data.incidents_open}</span>
          <span className="kpi-note">of {data.incidents_total} on the board</span>
        </div>
        <div className="kpi">
          <span className="kpi-label">Responders available</span>
          <span className="kpi-value">{data.responders_available}</span>
          <span className="kpi-note">of {respondersTotal} on the roster</span>
        </div>
        <div className="kpi">
          <span className="kpi-label">Shelter places in use</span>
          <span className="kpi-value">{data.occupancy_percent}%</span>
          <span className="kpi-note">
            {data.capacity_used} of {data.capacity_total} places
          </span>
        </div>
      </section>

      <section className="panel">
        <h2>Incidents by status</h2>
        <div className="tally">
          {statusOrder.map((status) => (
            <div key={status} className="tally-item">
              <span className="tally-value">{data.incidents_by_status[status]}</span>
              <span className={`status status-${status}`}>{statusLabels[status]}</span>
            </div>
          ))}
        </div>
      </section>

      <section className="panel">
        <h2>Incidents by priority</h2>
        <div className="tally">
          {bandOrder.map((band) => (
            <div key={band} className="tally-item">
              <span className="tally-value">{data.incidents_by_band[band]}</span>
              <span className={`band band-${band.toLowerCase()}`}>{band}</span>
            </div>
          ))}
        </div>
        <p className="panel-note">
          <Link to="/incidents">Open the incident list</Link> to work through them in priority order.
        </p>
      </section>

      <section className="panel">
        <h2>Shelter occupancy</h2>
        <div className="meters">
          {data.shelters.map((shelter) => (
            <OccupancyMeter
              key={shelter.id}
              name={shelter.name}
              used={shelter.capacity_used}
              total={shelter.capacity_total}
              percent={shelter.occupancy_percent}
            />
          ))}
        </div>
      </section>
    </>
  )
}
