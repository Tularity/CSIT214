package models

const (
	BandCritical = "Critical"
	BandHigh     = "High"
	BandMedium   = "Medium"
	BandLow      = "Low"
)

// Weightings come from the CoastLink Council triage rules agreed on 2026-09-19;
// they are deliberately simple so an operator can justify an ordering on the spot.
const (
	severityWeight    = 10
	peopleAffectedCap = 30
	vulnerableBonus   = 20
)

func PriorityScore(severity, peopleAffected int, vulnerable bool) int {
	if peopleAffected > peopleAffectedCap {
		peopleAffected = peopleAffectedCap
	}

	score := severity*severityWeight + peopleAffected
	if vulnerable {
		score += vulnerableBonus
	}
	return score
}

func PriorityBand(score int) string {
	switch {
	case score >= 60:
		return BandCritical
	case score >= 40:
		return BandHigh
	case score >= 20:
		return BandMedium
	default:
		return BandLow
	}
}
