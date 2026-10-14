package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// These two tests are the prototype's acceptance script: the end-to-end chain
// the client asked to see, and each validation rule it listed. Run them with
// -v to get a readable record of what was done, expected and observed.

func dashboardNow(t *testing.T, server *httptest.Server) dashboardBody {
	t.Helper()
	var dashboard dashboardBody
	get(t, server, "/api/dashboard", &dashboard)
	return dashboard
}

func TestAcceptanceEndToEndChain(t *testing.T) {
	server := newTestServer(t)
	before := dashboardNow(t, server)

	step := func(name string, method, path string, body any, wantStatus int) map[string]any {
		t.Helper()
		var result map[string]any
		response := send(t, server, method, path, body, &result)
		if response.StatusCode != wantStatus {
			t.Fatalf("%s: %s %s returned %d (%v), want %d", name, method, path, response.StatusCode, result, wantStatus)
		}
		t.Logf("%-28s %s %s -> %d", name, method, path, response.StatusCode)
		return result
	}

	created := step("1. register incident", http.MethodPost, "/api/incidents", map[string]any{
		"type":            "flood",
		"location":        "Seaview Retirement Village, Coastvale",
		"description":     "Ground floor units flooding; several residents use walking frames",
		"people_affected": 14,
		"vulnerable":      true,
		"severity":        4,
	}, http.StatusCreated)
	id := int(created["id"].(float64))
	base := fmt.Sprintf("/api/incidents/%d", id)

	if created["priority_score"].(float64) != 74 || created["priority_band"] != "Critical" {
		t.Fatalf("scored %v (%v), want 4*10 + 14 + 20 = 74 (Critical)", created["priority_score"], created["priority_band"])
	}
	t.Logf("%-28s priority %v (%v), status %v", "   scored by the server", created["priority_score"], created["priority_band"], created["status"])

	var incidents []map[string]any
	get(t, server, "/api/incidents", &incidents)
	for position, incident := range incidents {
		if int(incident["id"].(float64)) == id {
			t.Logf("%-28s position %d of %d", "   placed in the list", position+1, len(incidents))
			if position > 0 && incidents[position-1]["priority_score"].(float64) < 74 {
				t.Errorf("an incident with a lower score is listed above it")
			}
		}
	}

	step("2. assess", http.MethodPost, base+"/assess", nil, http.StatusOK)

	var shelters []map[string]any
	get(t, server, "/api/shelters", &shelters)
	var shelterID float64
	for _, shelter := range shelters {
		if shelter["capacity_remaining"].(float64) >= 14 {
			shelterID = shelter["id"].(float64)
			break
		}
	}
	placed := step("3. place in a shelter", http.MethodPost, base+"/assign-shelter",
		map[string]any{"shelter_id": shelterID}, http.StatusOK)
	if placed["status"] != "assigned" {
		t.Fatalf("status after placement = %v, want assigned", placed["status"])
	}

	var responders []map[string]any
	get(t, server, "/api/responders", &responders)
	var responderID float64
	for _, responder := range responders {
		if responder["status"] == "available" {
			responderID = responder["id"].(float64)
			break
		}
	}
	dispatched := step("4. dispatch a responder", http.MethodPost, base+"/assign-responder",
		map[string]any{"responder_id": responderID}, http.StatusOK)
	if dispatched["status"] != "in_progress" {
		t.Fatalf("status after dispatch = %v, want in_progress", dispatched["status"])
	}

	during := dashboardNow(t, server)
	if during.CapacityUsed != before.CapacityUsed+14 || during.RespondersAvailable != before.RespondersAvailable-1 {
		t.Errorf("dashboard while in progress: %d places used, %d responders free; want %d and %d",
			during.CapacityUsed, during.RespondersAvailable, before.CapacityUsed+14, before.RespondersAvailable-1)
	}
	t.Logf("%-28s places used %d -> %d, responders free %d -> %d", "   dashboard follows",
		before.CapacityUsed, during.CapacityUsed, before.RespondersAvailable, during.RespondersAvailable)

	resolved := step("5. resolve", http.MethodPost, base+"/resolve", nil, http.StatusOK)
	if resolved["status"] != "resolved" {
		t.Fatalf("status after resolve = %v, want resolved", resolved["status"])
	}

	after := dashboardNow(t, server)
	if after.CapacityUsed != before.CapacityUsed || after.RespondersAvailable != before.RespondersAvailable {
		t.Errorf("after resolving: %d places used, %d responders free; want them back at %d and %d",
			after.CapacityUsed, after.RespondersAvailable, before.CapacityUsed, before.RespondersAvailable)
	}
	if after.IncidentsByStatus["resolved"] != before.IncidentsByStatus["resolved"]+1 {
		t.Errorf("resolved count went from %d to %d, want one more",
			before.IncidentsByStatus["resolved"], after.IncidentsByStatus["resolved"])
	}
	t.Logf("%-28s places and responder released; resolved %d -> %d", "   dashboard after",
		before.IncidentsByStatus["resolved"], after.IncidentsByStatus["resolved"])

	var history []map[string]any
	get(t, server, fmt.Sprintf("/api/audit?entity_type=incident&entity_id=%d", id), &history)
	if len(history) < 4 {
		t.Fatalf("activity log has %d entries for the incident, want at least 4", len(history))
	}
	for i := len(history) - 1; i >= 0; i-- {
		entry := history[i]
		if entry["actor"] == "" || entry["at"] == "" {
			t.Errorf("entry %v is missing who or when", entry["id"])
		}
		t.Logf("%-28s %v by %v at %v", "   log", entry["action"], entry["actor"], entry["at"])
	}
}

func TestAcceptanceValidationRules(t *testing.T) {
	server := newTestServer(t)

	rules := []struct {
		rule       string
		path       string
		body       any
		wantStatus int
		wantError  string
	}{
		{"empty location", "/api/incidents",
			map[string]any{"type": "storm", "location": "", "people_affected": 3, "severity": 2},
			http.StatusBadRequest, "Location is required (max 120 characters)"},
		{"location over 120 characters", "/api/incidents",
			map[string]any{"type": "storm", "location": fmt.Sprintf("%0121d", 0), "people_affected": 3, "severity": 2},
			http.StatusBadRequest, "Location is required (max 120 characters)"},
		{"people affected above 500", "/api/incidents",
			map[string]any{"type": "storm", "location": "Esplanade", "people_affected": 501, "severity": 2},
			http.StatusBadRequest, "People affected must be between 1 and 500"},
		{"severity of 0", "/api/incidents",
			map[string]any{"type": "storm", "location": "Esplanade", "people_affected": 3, "severity": 0},
			http.StatusBadRequest, "Severity must be between 1 and 5"},
		{"shelter without enough room", "/api/incidents/3/assign-shelter",
			map[string]any{"shelter_id": 2},
			http.StatusConflict, "Shelter has only 2 places remaining"},
		{"responder already on a job", "/api/incidents/4/assign-responder",
			map[string]any{"responder_id": 1},
			http.StatusConflict, "Responder is not available"},
		{"resolving a registered incident", "/api/incidents/5/resolve", nil,
			http.StatusConflict, "Incident must be assessed first"},
	}

	for _, rule := range rules {
		t.Run(rule.rule, func(t *testing.T) {
			var body map[string]any
			response := send(t, server, http.MethodPost, rule.path, rule.body, &body)
			t.Logf("expected %d %q, got %d %q", rule.wantStatus, rule.wantError, response.StatusCode, body["error"])
			if response.StatusCode != rule.wantStatus || body["error"] != rule.wantError {
				t.Fail()
			}
		})
	}

	var incidents []map[string]any
	get(t, server, "/api/incidents", &incidents)
	if len(incidents) != 8 {
		t.Errorf("rejected requests changed the data: %d incidents, want 8", len(incidents))
	}
}

// Change request CR-01, approved by the client on 13 October: a refused
// placement must offer the shelters that can take the whole group, and only
// those.
func TestAcceptanceCR01SheltersOfferedWhenFull(t *testing.T) {
	server := newTestServer(t)

	check := func(name string, people int) (map[string]int, string) {
		t.Helper()
		var created map[string]any
		send(t, server, http.MethodPost, "/api/incidents", map[string]any{
			"type": "flood", "location": "CR-01 check: " + name, "people_affected": people, "severity": 3,
		}, &created)
		id := int(created["id"].(float64))
		send(t, server, http.MethodPost, fmt.Sprintf("/api/incidents/%d/assess", id), nil, nil)

		var body fullShelterBody
		response := send(t, server, http.MethodPost, fmt.Sprintf("/api/incidents/%d/assign-shelter", id),
			map[string]any{"shelter_id": 2}, &body)
		if response.StatusCode != http.StatusConflict {
			t.Fatalf("%s: placing %d people at Northbay returned %d, want 409", name, people, response.StatusCode)
		}

		offered := map[string]int{}
		for _, alternative := range body.Alternatives {
			offered[alternative.Name] = alternative.CapacityRemaining
		}
		t.Logf("%-26s %d people refused with %q; offered %v", name, people, body.Error, offered)
		return offered, body.Error
	}

	offered, message := check("1. alternatives offered", 5)
	if message != "Shelter has only 2 places remaining" {
		t.Errorf("message changed to %q", message)
	}
	if len(offered) != 3 {
		t.Errorf("five people should be offered the three shelters with room, got %v", offered)
	}

	offered, _ = check("2. only shelters that fit", 50)
	for name, remaining := range offered {
		if remaining < 50 {
			t.Errorf("%s has only %d places but was offered for 50 people", name, remaining)
		}
	}
	for _, unsuitable := range []string{"Northbay Primary School Gym", "Pelican Point Surf Club"} {
		if _, ok := offered[unsuitable]; ok {
			t.Errorf("%s cannot take 50 people but was offered", unsuitable)
		}
	}
	if len(offered) != 2 {
		t.Errorf("fifty people fit only at Riverside and Harbour View, got %v", offered)
	}
}
