package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/Tularity/CSIT214/backend/db"
)

// newTestServer gives every test its own database file loaded from the same
// seed the demonstration uses, so a test failure means the code changed rather
// than an earlier test having left something behind.
func newTestServer(t *testing.T) *httptest.Server {
	t.Helper()

	t.Setenv("CSIT214_SEED_PATH", filepath.Join("..", "..", "seed", "seed.sql"))

	store, err := db.Open(filepath.Join(t.TempDir(), "test.db"))
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() {
		if err := store.Close(); err != nil {
			t.Errorf("close test database: %v", err)
		}
	})

	server := httptest.NewServer(New(store).Routes())
	t.Cleanup(server.Close)
	return server
}

func get(t *testing.T, server *httptest.Server, path string, target any) *http.Response {
	t.Helper()

	response, err := http.Get(server.URL + path)
	if err != nil {
		t.Fatalf("GET %s: %v", path, err)
	}
	t.Cleanup(func() { response.Body.Close() })

	if contentType := response.Header.Get("Content-Type"); contentType != "application/json; charset=utf-8" {
		t.Errorf("GET %s returned Content-Type %q, want JSON", path, contentType)
	}
	if target != nil {
		if err := json.NewDecoder(response.Body).Decode(target); err != nil {
			t.Fatalf("decode body of GET %s: %v", path, err)
		}
	}
	return response
}

func TestHealthReportsOk(t *testing.T) {
	server := newTestServer(t)

	var body map[string]string
	response := get(t, server, "/api/health", &body)

	if response.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", response.StatusCode, http.StatusOK)
	}
	if body["status"] != "ok" {
		t.Errorf("status field = %q, want %q", body["status"], "ok")
	}
}
