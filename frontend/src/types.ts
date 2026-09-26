export type Health = { status: string }

export type IncidentType = 'flood' | 'fire' | 'storm' | 'medical' | 'other'

export type IncidentStatus =
  | 'registered'
  | 'assessed'
  | 'assigned'
  | 'in_progress'
  | 'resolved'

export type PriorityBand = 'Critical' | 'High' | 'Medium' | 'Low'

export type Incident = {
  id: number
  type: IncidentType
  location: string
  description: string
  people_affected: number
  vulnerable: boolean
  severity: number
  priority_score: number
  priority_band: PriorityBand
  status: IncidentStatus
  shelter_id: number | null
  responder_id: number | null
  reported_at: string
}

export type Shelter = {
  id: number
  name: string
  address: string
  capacity_total: number
  capacity_used: number
  capacity_remaining: number
}

export type Responder = {
  id: number
  name: string
  skill: string
  status: 'available' | 'assigned' | 'off_duty'
}
