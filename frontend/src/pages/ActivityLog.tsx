import Placeholder from '../components/Placeholder'

export default function ActivityLog() {
  return (
    <Placeholder
      title="Activity log"
      summary="Every write operation recorded as when, who, what action, which record and detail. Read only."
      waitingOn="GET /api/audit"
    />
  )
}
