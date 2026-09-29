package models

import (
	"strings"
	"testing"
)

func validIncident() NewIncident {
	return NewIncident{
		Type:           "flood",
		Location:       "Lowtide Crescent, Riverside",
		Description:    "Water over the footpath",
		PeopleAffected: 12,
		Severity:       3,
	}
}

func TestNewIncidentValidation(t *testing.T) {
	cases := []struct {
		name    string
		change  func(*NewIncident)
		wantErr string
	}{
		{"valid incident", func(*NewIncident) {}, ""},
		{"empty location", func(n *NewIncident) { n.Location = "" },
			"Location is required (max 120 characters)"},
		{"location of 121 characters", func(n *NewIncident) { n.Location = strings.Repeat("a", 121) },
			"Location is required (max 120 characters)"},
		{"location of 120 multi-byte characters", func(n *NewIncident) { n.Location = strings.Repeat("é", 120) }, ""},
		{"no people affected", func(n *NewIncident) { n.PeopleAffected = 0 },
			"People affected must be between 1 and 500"},
		{"too many people affected", func(n *NewIncident) { n.PeopleAffected = 501 },
			"People affected must be between 1 and 500"},
		{"severity below range", func(n *NewIncident) { n.Severity = 0 },
			"Severity must be between 1 and 5"},
		{"severity above range", func(n *NewIncident) { n.Severity = 6 },
			"Severity must be between 1 and 5"},
		{"unknown type", func(n *NewIncident) { n.Type = "earthquake" },
			"Type must be one of flood, fire, storm, medical or other"},
		{"overlong description", func(n *NewIncident) { n.Description = strings.Repeat("x", 1001) },
			"Description must be at most 1000 characters"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			incident := validIncident()
			testCase.change(&incident)

			err := incident.Normalized().Validate()
			got := ""
			if err != nil {
				got = err.Error()
			}
			if got != testCase.wantErr {
				t.Errorf("Validate() = %q, want %q", got, testCase.wantErr)
			}
		})
	}
}

func TestNormalizedTrimsBeforeValidating(t *testing.T) {
	incident := validIncident()
	incident.Location = "   "
	incident.Type = " Fire "

	normalized := incident.Normalized()
	if normalized.Type != "fire" {
		t.Errorf("Type = %q, want %q", normalized.Type, "fire")
	}
	if err := normalized.Validate(); err == nil || !strings.HasPrefix(err.Error(), "Location") {
		t.Errorf("a location of only spaces should be rejected, got %v", err)
	}
}
