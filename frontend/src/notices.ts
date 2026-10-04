import { createContext, useContext } from 'react'

export type Notice = {
  id: number
  kind: 'success' | 'error'
  message: string
}

export type Notify = (kind: Notice['kind'], message: string) => void

export const NoticeContext = createContext<Notify>(() => {})

export function useNotify(): Notify {
  return useContext(NoticeContext)
}
