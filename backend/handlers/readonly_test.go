package handlers

import (
	"net/http"
	"sort"
	"testing"
)

func TestIncidentsComeBackHighestPriorityFirst(t *testing.T) {
	server := newTestServer(t)

	var incidents []map[string]any
	response := get(t, server, "/api/incidents", &incidents)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if len(incidents) != 8 {
		t.Fatalf("got %d incidents, want the 8 from the sample data", len(incidents))
	}

	for i := 1; i < len(incidents); i++ {
		previous := incidents[i-1]["priority_score"].(float64)
		current := incidents[i]["priority_score"].(float64)
		if current > previous {
			t.Errorf("incident at position %d scores %.0f, higher than %.0f above it",
				i, current, previous)
		}
	}

	if top := incidents[0]["priority_band"]; top != "Critical" {
		t.Errorf("first incident is banded %v, want Critical", top)
	}
}

func TestIncidentFieldNamesMatchTheDataModel(t *testing.T) {
	server := newTestServer(t)

	var incidents []map[string]any
	get(t, server, "/api/incidents", &incidents)

	want := []string{
		"description", "id", "location", "people_affected", "priority_band",
		"priority_score", "reported_at", "responder_id", "severity",
		"shelter_id", "status", "type", "vulnerable",
	}
	assertFieldNames(t, "incident", incidents[0], want)

	for _, incident := range incidents {
		if _, ok := incident["vulnerable"].(bool); !ok {
			t.Errorf("incident %v has vulnerable=%v, want a JSON boolean",
				incident["id"], incident["vulnerable"])
		}
	}
}

func TestSingleIncidentIsAddressableById(t *testing.T) {
	server := newTestServer(t)

	var incident map[string]any
	response := get(t, server, "/api/incidents/2", &incident)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if incident["id"].(float64) != 2 {
		t.Errorf("id = %v, want 2", incident["id"])
	}
	if incident["shelter_id"] == nil {
		t.Error("incident 2 is placed in a shelter but shelter_id came back null")
	}

	var unplaced map[string]any
	get(t, server, "/api/incidents/1", &unplaced)
	if unplaced["shelter_id"] != nil {
		t.Errorf("incident 1 has no shelter yet, want null, got %v", unplaced["shelter_id"])
	}
}

func TestMissingAndMalformedIncidentIdsAreRejected(t *testing.T) {
	server := newTestServer(t)

	cases := []struct {
		path       string
		wantStatus int
	}{
		{"/api/incidents/999", http.StatusNotFound},
		{"/api/incidents/0", http.StatusBadRequest},
		{"/api/incidents/abc", http.StatusBadRequest},
	}

	for _, testCase := range cases {
		var body map[string]string
		response := get(t, server, testCase.path, &body)

		if response.StatusCode != testCase.wantStatus {
			t.Errorf("GET %s returned %d, want %d",
				testCase.path, response.StatusCode, testCase.wantStatus)
		}
		if body["error"] == "" {
			t.Errorf("GET %s returned no error message for the operator", testCase.path)
		}
	}
}

func TestSheltersReportRemainingCapacity(t *testing.T) {
	server := newTestServer(t)

	var shelters []map[string]any
	response := get(t, server, "/api/shelters", &shelters)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if len(shelters) != 4 {
		t.Fatalf("got %d shelters, want the 4 from the sample data", len(shelters))
	}

	assertFieldNames(t, "shelter", shelters[0], []string{
		"address", "capacity_remaining", "capacity_total", "capacity_used", "id", "name",
	})

	for _, shelter := range shelters {
		total := shelter["capacity_total"].(float64)
		used := shelter["capacity_used"].(float64)
		remaining := shelter["capacity_remaining"].(float64)
		if remaining != total-used {
			t.Errorf("shelter %v reports %.0f places left, want %.0f",
				shelter["id"], remaining, total-used)
		}
	}
}

func TestRespondersExposeTheirAvailability(t *testing.T) {
	server := newTestServer(t)

	var responders []map[string]any
	response := get(t, server, "/api/responders", &responders)

	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if len(responders) != 6 {
		t.Fatalf("got %d responders, want the 6 from the sample data", len(responders))
	}

	assertFieldNames(t, "responder", responders[0], []string{"id", "name", "skill", "status"})

	allowed := map[string]bool{"available": true, "assigned": true, "off_duty": true}
	available := 0
	for _, responder := range responders {
		status := responder["status"].(string)
		if !allowed[status] {
			t.Errorf("responder %v has status %q, which is not one of the three allowed values",
				responder["id"], status)
		}
		if status == "available" {
			available++
		}
	}
	if available == 0 {
		t.Error("no responder is available, so nothing could be dispatched in a demonstration")
	}
}

func assertFieldNames(t *testing.T, kind string, record map[string]any, want []string) {
	t.Helper()

	got := make([]string, 0, len(record))
	for field := range record {
		got = append(got, field)
	}
	sort.Strings(got)

	if len(got) != len(want) {
		t.Fatalf("%s fields = %v, want %v", kind, got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("%s fields = %v, want %v", kind, got, want)
		}
	}
}
