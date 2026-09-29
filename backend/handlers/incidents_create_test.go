package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func send(t *testing.T, server *httptest.Server, method, path string, body, target any) *http.Response {
	t.Helper()

	var payload bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&payload).Encode(body); err != nil {
			t.Fatalf("encode body for %s %s: %v", method, path, err)
		}
	}

	request, err := http.NewRequest(method, server.URL+path, &payload)
	if err != nil {
		t.Fatalf("build %s %s: %v", method, path, err)
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		t.Fatalf("%s %s: %v", method, path, err)
	}
	t.Cleanup(func() { response.Body.Close() })

	if target != nil {
		if err := json.NewDecoder(response.Body).Decode(target); err != nil {
			t.Fatalf("decode body of %s %s: %v", method, path, err)
		}
	}
	return response
}

func TestCreateIncidentScoresAndRegistersIt(t *testing.T) {
	server := newTestServer(t)

	var created map[string]any
	response := send(t, server, http.MethodPost, "/api/incidents", map[string]any{
		"type":            "flood",
		"location":        "  Harbour Road underpass, Coastvale  ",
		"description":     "Two cars stranded in rising water",
		"people_affected": 30,
		"vulnerable":      true,
		"severity":        5,
	}, &created)

	if response.StatusCode != http.StatusCreated {
		t.Fatalf("status = %d, want %d (%v)", response.StatusCode, http.StatusCreated, created)
	}
	if created["priority_score"].(float64) != 100 || created["priority_band"] != "Critical" {
		t.Errorf("scored %v (%v), want 100 (Critical)", created["priority_score"], created["priority_band"])
	}
	if created["status"] != "registered" {
		t.Errorf("status = %v, want registered", created["status"])
	}
	if created["location"] != "Harbour Road underpass, Coastvale" {
		t.Errorf("location = %q, want it trimmed", created["location"])
	}

	id := created["id"].(float64)
	if location := response.Header.Get("Location"); location == "" {
		t.Error("201 response carries no Location header")
	}

	var incidents []map[string]any
	get(t, server, "/api/incidents", &incidents)
	if incidents[0]["id"].(float64) != id {
		t.Errorf("the new Critical incident is not first in the list, got id %v", incidents[0]["id"])
	}
}

func TestCreateIncidentRejectsInvalidInput(t *testing.T) {
	server := newTestServer(t)

	cases := []struct {
		name   string
		body   any
		status int
		error  string
	}{
		{"empty location", map[string]any{"type": "fire", "location": "", "people_affected": 4, "severity": 2},
			http.StatusBadRequest, "Location is required (max 120 characters)"},
		{"people affected out of range", map[string]any{"type": "fire", "location": "Depot", "people_affected": 0, "severity": 2},
			http.StatusBadRequest, "People affected must be between 1 and 500"},
		{"severity out of range", map[string]any{"type": "fire", "location": "Depot", "people_affected": 4, "severity": 9},
			http.StatusBadRequest, "Severity must be between 1 and 5"},
		{"not JSON at all", "people_affected=4",
			http.StatusBadRequest, "Request body must be a JSON object with the expected fields"},
	}

	for _, testCase := range cases {
		t.Run(testCase.name, func(t *testing.T) {
			var body map[string]string
			response := send(t, server, http.MethodPost, "/api/incidents", testCase.body, &body)
			if response.StatusCode != testCase.status || body["error"] != testCase.error {
				t.Errorf("got %d %q, want %d %q", response.StatusCode, body["error"], testCase.status, testCase.error)
			}
		})
	}

	var incidents []map[string]any
	get(t, server, "/api/incidents", &incidents)
	if len(incidents) != 8 {
		t.Errorf("rejected requests changed the incident count to %d", len(incidents))
	}
}
