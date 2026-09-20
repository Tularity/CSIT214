type PlaceholderProps = {
  title: string
  summary: string
  waitingOn: string
}

export default function Placeholder({ title, summary, waitingOn }: PlaceholderProps) {
  return (
    <section className="panel">
      <h2>{title}</h2>
      <p>{summary}</p>
      <p className="panel-note">Waiting on {waitingOn}.</p>
    </section>
  )
}
