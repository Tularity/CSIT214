package handlers

import (
	"fmt"
	"net/http"
	"testing"
)

func shelterByID(t *testing.T, shelters []map[string]any, id float64) map[string]any {
	t.Helper()
	for _, shelter := range shelters {
		if shelter["id"].(float64) == id {
			return shelter
		}
	}
	t.Fatalf("shelter %v missing from the list", id)
	return nil
}

func responderStatus(t *testing.T, responders []map[string]any, id float64) string {
	t.Helper()
	for _, responder := range responders {
		if responder["id"].(float64) == id {
			return responder["status"].(string)
		}
	}
	t.Fatalf("responder %v missing from the list", id)
	return ""
}

func expectConflict(t *testing.T, response *http.Response, body map[string]any, want string) {
	t.Helper()
	if response.StatusCode != http.StatusConflict || body["error"] != want {
		t.Errorf("got %d %q, want 409 %q", response.StatusCode, body["error"], want)
	}
}

// Incident 3 in the sample data is assessed with five people affected, and
// shelter 2 has exactly two places left: the case the client asked to see.
func TestShelterWithoutRoomIsRefused(t *testing.T) {
	server := newTestServer(t)

	var body map[string]any
	response := send(t, server, http.MethodPost, "/api/incidents/3/assign-shelter",
		map[string]any{"shelter_id": 2}, &body)
	expectConflict(t, response, body, "Shelter has only 2 places remaining")

	var shelters []map[string]any
	get(t, server, "/api/shelters", &shelters)
	if used := shelterByID(t, shelters, 2)["capacity_used"].(float64); used != 38 {
		t.Errorf("a refused placement still changed capacity_used to %v", used)
	}
}

func TestStepsCannotBeSkipped(t *testing.T) {
	server := newTestServer(t)

	cases := []struct {
		path string
		body any
		want string
	}{
		{"/api/incidents/1/resolve", nil, "Incident must be assessed first"},
		{"/api/incidents/1/assign-shelter", map[string]any{"shelter_id": 4}, "Incident must be assessed first"},
		{"/api/incidents/1/assign-responder", map[string]any{"responder_id": 2}, "Incident must be assessed first"},
		{"/api/incidents/3/assess", nil, "Incident has already been assessed"},
		{"/api/incidents/7/resolve", nil, "Incident is already resolved"},
		{"/api/incidents/2/assign-responder", map[string]any{"responder_id": 2}, "Incident already has a responder"},
	}
	for _, testCase := range cases {
		var body map[string]any
		response := send(t, server, http.MethodPost, testCase.path, testCase.body, &body)
		expectConflict(t, response, body, testCase.want)
	}
}

func TestUnavailableResponderIsRefused(t *testing.T) {
	server := newTestServer(t)

	for _, responderID := range []int{1, 5} {
		var body map[string]any
		response := send(t, server, http.MethodPost, "/api/incidents/4/assign-responder",
			map[string]any{"responder_id": responderID}, &body)
		expectConflict(t, response, body, "Responder is not available")
	}
}

func TestMissingRecordsAndBadBodiesAreReported(t *testing.T) {
	server := newTestServer(t)

	cases := []struct {
		path   string
		body   any
		status int
		want   string
	}{
		{"/api/incidents/999/assess", nil, http.StatusNotFound, "Incident not found"},
		{"/api/incidents/3/assign-shelter", map[string]any{"shelter_id": 99}, http.StatusNotFound, "Shelter not found"},
		{"/api/incidents/3/assign-responder", map[string]any{"responder_id": 99}, http.StatusNotFound, "Responder not found"},
		{"/api/incidents/3/assign-shelter", map[string]any{}, http.StatusBadRequest, "shelter_id must be a positive number"},
	}
	for _, testCase := range cases {
		var body map[string]any
		response := send(t, server, http.MethodPost, testCase.path, testCase.body, &body)
		if response.StatusCode != testCase.status || body["error"] != testCase.want {
			t.Errorf("POST %s: got %d %q, want %d %q",
				testCase.path, response.StatusCode, body["error"], testCase.status, testCase.want)
		}
	}
}

func TestResolvingReleasesShelterAndResponder(t *testing.T) {
	server := newTestServer(t)
	const incident, shelter, responder = 1, 4, 2 // 18 people, Riverside Civic Centre, Bianca Torres

	steps := []struct {
		action, wantStatus string
		body               any
	}{
		{"assess", "assessed", nil},
		{"assign-shelter", "assigned", map[string]any{"shelter_id": shelter}},
		{"assign-responder", "in_progress", map[string]any{"responder_id": responder}},
	}
	for _, step := range steps {
		var updated map[string]any
		response := send(t, server, http.MethodPost,
			fmt.Sprintf("/api/incidents/%d/%s", incident, step.action), step.body, &updated)
		if response.StatusCode != http.StatusOK || updated["status"] != step.wantStatus {
			t.Fatalf("%s: got %d status %v, want 200 %s", step.action, response.StatusCode, updated["status"], step.wantStatus)
		}
	}

	var shelters []map[string]any
	var responders []map[string]any
	get(t, server, "/api/shelters", &shelters)
	get(t, server, "/api/responders", &responders)
	if used := shelterByID(t, shelters, shelter)["capacity_used"].(float64); used != 18 {
		t.Errorf("after placement capacity_used = %v, want 18", used)
	}
	if status := responderStatus(t, responders, responder); status != "assigned" {
		t.Errorf("after dispatch responder status = %q, want assigned", status)
	}

	var resolved map[string]any
	send(t, server, http.MethodPost, fmt.Sprintf("/api/incidents/%d/resolve", incident), nil, &resolved)
	if resolved["status"] != "resolved" {
		t.Fatalf("resolve left status %v", resolved["status"])
	}

	get(t, server, "/api/shelters", &shelters)
	get(t, server, "/api/responders", &responders)
	if used := shelterByID(t, shelters, shelter)["capacity_used"].(float64); used != 0 {
		t.Errorf("after resolving capacity_used = %v, want the 18 places back (0)", used)
	}
	if status := responderStatus(t, responders, responder); status != "available" {
		t.Errorf("after resolving responder status = %q, want available", status)
	}
}
