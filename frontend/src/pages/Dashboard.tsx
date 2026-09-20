import Placeholder from '../components/Placeholder'

export default function Dashboard() {
  return (
    <Placeholder
      title="Operations overview"
      summary="Incident counts per status and per priority band, shelter occupancy and the number of available responders."
      waitingOn="GET /api/dashboard"
    />
  )
}
