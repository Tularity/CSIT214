package models

import (
	"strings"
	"unicode/utf8"
)

const (
	maxLocationLength    = 120
	maxDescriptionLength = 1000
	minPeopleAffected    = 1
	maxPeopleAffected    = 500
	minSeverity          = 1
	maxSeverity          = 5
)

var incidentTypes = map[string]bool{
	"flood": true, "fire": true, "storm": true, "medical": true, "other": true,
}

type NewIncident struct {
	Type           string `json:"type"`
	Location       string `json:"location"`
	Description    string `json:"description"`
	PeopleAffected int    `json:"people_affected"`
	Vulnerable     bool   `json:"vulnerable"`
	Severity       int    `json:"severity"`
}

type ValidationError struct {
	Message string
}

func (e *ValidationError) Error() string {
	return e.Message
}

// Normalized trims the free-text fields so that validation and storage see
// exactly the same value.
func (n NewIncident) Normalized() NewIncident {
	n.Type = strings.ToLower(strings.TrimSpace(n.Type))
	n.Location = strings.TrimSpace(n.Location)
	n.Description = strings.TrimSpace(n.Description)
	return n
}

// Validate checks the rules in the order the client listed them, so the
// operator is told about the first field on the form that is wrong.
func (n NewIncident) Validate() error {
	if n.Location == "" || utf8.RuneCountInString(n.Location) > maxLocationLength {
		return &ValidationError{"Location is required (max 120 characters)"}
	}
	if n.PeopleAffected < minPeopleAffected || n.PeopleAffected > maxPeopleAffected {
		return &ValidationError{"People affected must be between 1 and 500"}
	}
	if n.Severity < minSeverity || n.Severity > maxSeverity {
		return &ValidationError{"Severity must be between 1 and 5"}
	}
	if !incidentTypes[n.Type] {
		return &ValidationError{"Type must be one of flood, fire, storm, medical or other"}
	}
	if utf8.RuneCountInString(n.Description) > maxDescriptionLength {
		return &ValidationError{"Description must be at most 1000 characters"}
	}
	return nil
}
