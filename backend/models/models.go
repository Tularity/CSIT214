package models

type Incident struct {
	ID             int64  `json:"id"`
	Type           string `json:"type"`
	Location       string `json:"location"`
	Description    string `json:"description"`
	PeopleAffected int    `json:"people_affected"`
	Vulnerable     bool   `json:"vulnerable"`
	Severity       int    `json:"severity"`
	PriorityScore  int    `json:"priority_score"`
	PriorityBand   string `json:"priority_band"`
	Status         string `json:"status"`
	ShelterID      *int64 `json:"shelter_id"`
	ResponderID    *int64 `json:"responder_id"`
	ReportedAt     string `json:"reported_at"`
}

type Shelter struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Address           string `json:"address"`
	CapacityTotal     int    `json:"capacity_total"`
	CapacityUsed      int    `json:"capacity_used"`
	CapacityRemaining int    `json:"capacity_remaining"`
}

type Responder struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Skill  string `json:"skill"`
	Status string `json:"status"`
}
