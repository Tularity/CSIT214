import Placeholder from '../components/Placeholder'

export default function Incidents() {
  return (
    <Placeholder
      title="Incidents"
      summary="Type, location, people affected, priority band and status, highest priority first."
      waitingOn="GET /api/incidents"
    />
  )
}
