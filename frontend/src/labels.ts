import type { Incident } from './types'

export const typeLabels: Record<Incident['type'], string> = {
  flood: 'Flood',
  fire: 'Fire',
  storm: 'Storm',
  medical: 'Medical',
  other: 'Other',
}

export const statusLabels: Record<Incident['status'], string> = {
  registered: 'Registered',
  assessed: 'Assessed',
  assigned: 'Shelter assigned',
  in_progress: 'In progress',
  resolved: 'Resolved',
}
