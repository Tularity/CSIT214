import { useCallback, useEffect, useRef, useState, type ReactNode } from 'react'
import { NoticeContext, type Notice } from '../notices'

const successLifetimeMs = 5000

export default function NoticeProvider({ children }: { children: ReactNode }) {
  const [notices, setNotices] = useState<Notice[]>([])
  const nextId = useRef(1)
  const timers = useRef(new Map<number, number>())

  const dismiss = useCallback((id: number) => {
    setNotices((current) => current.filter((notice) => notice.id !== id))
    window.clearTimeout(timers.current.get(id))
    timers.current.delete(id)
  }, [])

  // Errors stay until the operator closes them: a rejected step is something
  // they need to read, not something that should disappear while they look away.
  const notify = useCallback(
    (kind: Notice['kind'], message: string) => {
      const id = nextId.current++
      setNotices((current) => [...current, { id, kind, message }])
      if (kind === 'success') {
        timers.current.set(id, window.setTimeout(() => dismiss(id), successLifetimeMs))
      }
    },
    [dismiss],
  )

  useEffect(() => {
    const pending = timers.current
    return () => pending.forEach((timer) => window.clearTimeout(timer))
  }, [])

  return (
    <NoticeContext.Provider value={notify}>
      <div className="notices" aria-live="polite">
        {notices.map((notice) => (
          <div key={notice.id} className={`notice notice-${notice.kind}`} role={notice.kind === 'error' ? 'alert' : 'status'}>
            <span>{notice.message}</span>
            <button type="button" className="notice-close" aria-label="Dismiss" onClick={() => dismiss(notice.id)}>
              ×
            </button>
          </div>
        ))}
      </div>
      {children}
    </NoticeContext.Provider>
  )
}
