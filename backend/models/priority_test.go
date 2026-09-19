package models

import "testing"

func TestPriorityScore(t *testing.T) {
	cases := []struct {
		name            string
		severity        int
		peopleAffected  int
		vulnerable      bool
		wantScore       int
		wantBand        string
	}{
		{"worked example from the specification", 5, 30, true, 100, BandCritical},
		{"people affected are capped at thirty", 5, 120, true, 100, BandCritical},
		{"vulnerable bonus lifts the band", 2, 5, true, 45, BandHigh},
		{"same incident without vulnerable people", 2, 5, false, 25, BandMedium},
		{"smallest possible incident", 1, 1, false, 11, BandLow},
		{"band boundary at sixty", 4, 20, false, 60, BandCritical},
		{"band boundary at forty", 3, 10, false, 40, BandHigh},
		{"band boundary at twenty", 1, 10, false, 20, BandMedium},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			score := PriorityScore(testCase.severity, testCase.peopleAffected, testCase.vulnerable)
			if score != testCase.wantScore {
				t.Errorf("PriorityScore(%d, %d, %t) = %d, want %d",
					testCase.severity, testCase.peopleAffected, testCase.vulnerable, score, testCase.wantScore)
			}
			if band := PriorityBand(score); band != testCase.wantBand {
				t.Errorf("PriorityBand(%d) = %q, want %q", score, band, testCase.wantBand)
			}
		})
	}
}
