import { useHealth } from '../useHealth'

export default function Dashboard() {
  const state = useHealth()

  return (
    <>
      <section className="panel">
        <h2>Backend connection</h2>
        <p>Response from <code>GET /api/health</code> through the development proxy:</p>
        {state.kind === 'loading' && <p className="panel-note">Requesting…</p>}
        {state.kind === 'failed' && <p className="panel-error">{state.message}</p>}
        {state.kind === 'ready' && <pre className="panel-code">{JSON.stringify(state.health, null, 2)}</pre>}
      </section>

      <section className="panel">
        <h2>Operations overview</h2>
        <p>
          Incident counts per status and per priority band, shelter occupancy and the number of
          available responders.
        </p>
        <p className="panel-note">Waiting on <code>GET /api/dashboard</code>.</p>
      </section>
    </>
  )
}
