import type { PriorityBand } from '../types'

type PriorityBadgeProps = {
  band: PriorityBand
  score: number
}

export default function PriorityBadge({ band, score }: PriorityBadgeProps) {
  return (
    <span className={`band band-${band.toLowerCase()}`}>
      {band}
      <span className="band-score">{score}</span>
    </span>
  )
}
