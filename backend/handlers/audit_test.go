package handlers

import (
	"net/http"
	"testing"
)

func TestAuditIsNewestFirst(t *testing.T) {
	server := newTestServer(t)

	var entries []map[string]any
	response := get(t, server, "/api/audit", &entries)
	if response.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", response.StatusCode)
	}
	if len(entries) != 21 {
		t.Fatalf("got %d entries, want the 21 from the sample data", len(entries))
	}
	assertFieldNames(t, "audit entry", entries[0], []string{
		"action", "actor", "at", "detail", "entity_id", "entity_type", "id",
	})
	for i := 1; i < len(entries); i++ {
		if entries[i]["at"].(string) > entries[i-1]["at"].(string) {
			t.Fatalf("entry %d (%v) is newer than the one above it (%v)",
				i, entries[i]["at"], entries[i-1]["at"])
		}
	}
}

func TestAuditCanBeNarrowedToOneIncident(t *testing.T) {
	server := newTestServer(t)

	var entries []map[string]any
	get(t, server, "/api/audit?entity_type=incident&entity_id=7", &entries)
	if len(entries) != 5 {
		t.Fatalf("incident 7 has %d entries, want 5", len(entries))
	}
	for _, entry := range entries {
		if entry["entity_id"].(float64) != 7 {
			t.Errorf("filter let through an entry for %v", entry["entity_id"])
		}
	}

	var body map[string]string
	response := get(t, server, "/api/audit?entity_id=abc", &body)
	if response.StatusCode != http.StatusBadRequest {
		t.Errorf("a malformed entity_id returned %d, want 400", response.StatusCode)
	}
}
