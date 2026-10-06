package handlers

import (
	"fmt"
	"net/http"
	"testing"
)

type dashboardBody struct {
	IncidentsTotal      int            `json:"incidents_total"`
	IncidentsOpen       int            `json:"incidents_open"`
	IncidentsByStatus   map[string]int `json:"incidents_by_status"`
	IncidentsByBand     map[string]int `json:"incidents_by_band"`
	CapacityTotal       int            `json:"capacity_total"`
	CapacityUsed        int            `json:"capacity_used"`
	OccupancyPercent    int            `json:"occupancy_percent"`
	RespondersAvailable int            `json:"responders_available"`
}

func TestDashboardSummarisesTheSampleData(t *testing.T) {
	server := newTestServer(t)

	var dashboard dashboardBody
	get(t, server, "/api/dashboard", &dashboard)

	wantStatus := map[string]int{"registered": 3, "assessed": 1, "assigned": 1, "in_progress": 1, "resolved": 2}
	for status, want := range wantStatus {
		if got := dashboard.IncidentsByStatus[status]; got != want {
			t.Errorf("%s = %d, want %d", status, got, want)
		}
	}
	wantBand := map[string]int{"Critical": 4, "High": 1, "Medium": 2, "Low": 1}
	for band, want := range wantBand {
		if got := dashboard.IncidentsByBand[band]; got != want {
			t.Errorf("%s = %d, want %d", band, got, want)
		}
	}
	if dashboard.IncidentsTotal != 8 || dashboard.IncidentsOpen != 6 {
		t.Errorf("total/open = %d/%d, want 8/6", dashboard.IncidentsTotal, dashboard.IncidentsOpen)
	}
	if dashboard.CapacityTotal != 420 || dashboard.CapacityUsed != 94 || dashboard.OccupancyPercent != 22 {
		t.Errorf("occupancy = %d/%d (%d%%), want 94/420 (22%%)",
			dashboard.CapacityUsed, dashboard.CapacityTotal, dashboard.OccupancyPercent)
	}
	if dashboard.RespondersAvailable != 4 {
		t.Errorf("available responders = %d, want 4", dashboard.RespondersAvailable)
	}
}

// The acceptance check for the activity log: one full pass through the chain
// must leave at least four entries for that incident, one per write.
func TestEveryStepOfTheChainIsAudited(t *testing.T) {
	server := newTestServer(t)

	var created map[string]any
	send(t, server, http.MethodPost, "/api/incidents", map[string]any{
		"type": "fire", "location": "Coastvale Library", "people_affected": 6, "severity": 3,
	}, &created)
	id := int(created["id"].(float64))

	for _, step := range []struct {
		action string
		body   any
	}{
		{"assess", nil},
		{"assign-shelter", map[string]any{"shelter_id": 1}},
		{"assign-responder", map[string]any{"responder_id": 3}},
		{"resolve", nil},
	} {
		response := send(t, server, http.MethodPost, fmt.Sprintf("/api/incidents/%d/%s", id, step.action), step.body, nil)
		if response.StatusCode != http.StatusOK {
			t.Fatalf("%s returned %d", step.action, response.StatusCode)
		}
	}

	var entries []map[string]any
	get(t, server, fmt.Sprintf("/api/audit?entity_type=incident&entity_id=%d", id), &entries)

	wantActions := []string{"resolve", "assign-responder", "assign-shelter", "assess", "register"}
	if len(entries) != len(wantActions) {
		t.Fatalf("got %d audit entries, want %d", len(entries), len(wantActions))
	}
	for i, want := range wantActions {
		if entries[i]["action"] != want {
			t.Errorf("entry %d action = %v, want %s", i, entries[i]["action"], want)
		}
		if entries[i]["actor"] == "" {
			t.Errorf("entry %d has no actor", i)
		}
	}

	var dashboard dashboardBody
	get(t, server, "/api/dashboard", &dashboard)
	if dashboard.IncidentsByStatus["resolved"] != 3 {
		t.Errorf("resolved count = %d after the chain, want 3", dashboard.IncidentsByStatus["resolved"])
	}
}
