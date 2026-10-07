type OccupancyMeterProps = {
  name: string
  used: number
  total: number
  percent: number
}

function occupancyLevel(percent: number): { level: 'room' | 'filling' | 'full'; label: string } {
  if (percent >= 100) return { level: 'full', label: 'Full' }
  if (percent >= 90) return { level: 'full', label: 'Nearly full' }
  if (percent >= 75) return { level: 'filling', label: 'Filling up' }
  return { level: 'room', label: 'Room available' }
}

export default function OccupancyMeter({ name, used, total, percent }: OccupancyMeterProps) {
  const { level, label } = occupancyLevel(percent)
  const summary = `${name}: ${used} of ${total} places used (${percent}%), ${label.toLowerCase()}`

  return (
    <div className="meter-row" title={summary}>
      <div className="meter-heading">
        <span className="meter-name">{name}</span>
        <span className="meter-figures">
          {used} / {total} · <span className={`meter-state meter-state-${level}`}>{label}</span>
        </span>
      </div>
      <div
        className={`meter meter-${level}`}
        role="meter"
        aria-label={summary}
        aria-valuemin={0}
        aria-valuemax={total}
        aria-valuenow={used}
      >
        <div className="meter-fill" style={{ width: `${Math.min(percent, 100)}%` }} />
      </div>
    </div>
  )
}
