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

type AuditEntry struct {
	ID         int64  `json:"id"`
	At         string `json:"at"`
	Actor      string `json:"actor"`
	Action     string `json:"action"`
	EntityType string `json:"entity_type"`
	EntityID   int64  `json:"entity_id"`
	Detail     string `json:"detail"`
}

type ShelterOccupancy struct {
	ID               int64  `json:"id"`
	Name             string `json:"name"`
	CapacityTotal    int    `json:"capacity_total"`
	CapacityUsed     int    `json:"capacity_used"`
	OccupancyPercent int    `json:"occupancy_percent"`
}

type Dashboard struct {
	IncidentsTotal      int                `json:"incidents_total"`
	IncidentsOpen       int                `json:"incidents_open"`
	IncidentsByStatus   map[string]int     `json:"incidents_by_status"`
	IncidentsByBand     map[string]int     `json:"incidents_by_band"`
	Shelters            []ShelterOccupancy `json:"shelters"`
	CapacityTotal       int                `json:"capacity_total"`
	CapacityUsed        int                `json:"capacity_used"`
	OccupancyPercent    int                `json:"occupancy_percent"`
	RespondersAvailable int                `json:"responders_available"`
	RespondersByStatus  map[string]int     `json:"responders_by_status"`
}

type ShelterAlternative struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	CapacityRemaining int    `json:"capacity_remaining"`
}
