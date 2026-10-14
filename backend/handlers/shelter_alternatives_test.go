package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
)

type fullShelterBody struct {
	Error        string `json:"error"`
	Alternatives []struct {
		ID                int64  `json:"id"`
		Name              string `json:"name"`
		CapacityRemaining int    `json:"capacity_remaining"`
	} `json:"alternatives"`
}

// CR-01: incident 3 (five people) refused by shelter 2 (two places left)
// should be offered the three shelters that can take five, roomiest first.
func TestFullShelterOffersSheltersWithRoom(t *testing.T) {
	server := newTestServer(t)

	var body fullShelterBody
	response := send(t, server, http.MethodPost, "/api/incidents/3/assign-shelter",
		map[string]any{"shelter_id": 2}, &body)

	if response.StatusCode != http.StatusConflict || body.Error != "Shelter has only 2 places remaining" {
		t.Fatalf("got %d %q, want 409 with the unchanged message", response.StatusCode, body.Error)
	}

	want := []struct {
		id        int64
		remaining int
	}{{4, 200}, {1, 96}, {3, 28}}
	if len(body.Alternatives) != len(want) {
		t.Fatalf("got %d alternatives %+v, want %d", len(body.Alternatives), body.Alternatives, len(want))
	}
	for i, expected := range want {
		got := body.Alternatives[i]
		if got.ID != expected.id || got.CapacityRemaining != expected.remaining || got.Name == "" {
			t.Errorf("alternative %d = %+v, want shelter %d with %d places", i, got, expected.id, expected.remaining)
		}
	}

	var placed map[string]any
	send(t, server, http.MethodPost, "/api/incidents/3/assign-shelter",
		map[string]any{"shelter_id": body.Alternatives[0].ID}, &placed)
	if placed["status"] != "assigned" {
		t.Errorf("choosing the first alternative left status %v, want assigned", placed["status"])
	}
}

func TestNoAlternativesIsAnEmptyList(t *testing.T) {
	server := newTestServer(t)

	var created map[string]any
	send(t, server, http.MethodPost, "/api/incidents", map[string]any{
		"type": "flood", "location": "Coastvale township", "people_affected": 250, "severity": 5,
	}, &created)
	id := int(created["id"].(float64))
	send(t, server, http.MethodPost, fmt.Sprintf("/api/incidents/%d/assess", id), nil, nil)

	response := send(t, server, http.MethodPost, fmt.Sprintf("/api/incidents/%d/assign-shelter", id),
		map[string]any{"shelter_id": 4}, nil)
	var raw map[string]json.RawMessage
	if err := json.NewDecoder(response.Body).Decode(&raw); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if response.StatusCode != http.StatusConflict {
		t.Fatalf("status = %d, want 409", response.StatusCode)
	}
	if string(raw["alternatives"]) != "[]" {
		t.Errorf("alternatives = %s, want [] when no shelter can take 250 people", raw["alternatives"])
	}
}
