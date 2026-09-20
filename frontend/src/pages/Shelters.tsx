import Placeholder from '../components/Placeholder'

export default function Shelters() {
  return (
    <Placeholder
      title="Shelters"
      summary="Name, address, total capacity, places used and places remaining."
      waitingOn="GET /api/shelters"
    />
  )
}
