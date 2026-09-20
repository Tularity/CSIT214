import Placeholder from '../components/Placeholder'

export default function Responders() {
  return (
    <Placeholder
      title="Responders"
      summary="Name, skill and whether the responder is available, assigned or off duty."
      waitingOn="GET /api/responders"
    />
  )
}
