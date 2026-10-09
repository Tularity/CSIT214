import type { ReactNode } from 'react'

export function Loading({ what }: { what: string }) {
  return (
    <p className="load-state" role="status">
      <span className="spinner" aria-hidden="true" />
      Loading {what}…
    </p>
  )
}

export function Empty({ children }: { children: ReactNode }) {
  return <div className="load-state load-state-empty">{children}</div>
}

export function LoadFailure({ message, onRetry }: { message: string; onRetry: () => void }) {
  return (
    <div className="load-state load-state-failed" role="alert">
      <span>{message}</span>
      <button type="button" className="button" onClick={onRetry}>
        Try again
      </button>
    </div>
  )
}
