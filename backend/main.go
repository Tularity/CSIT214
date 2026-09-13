package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "modernc.org/sqlite"

	"github.com/Tularity/CSIT214/backend/handlers"
)

const defaultDBPath = "data.db"

func main() {
	database, err := openDatabase()
	if err != nil {
		log.Fatalf("database unavailable: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handlers.Health)

	log.Print("listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func openDatabase() (*sql.DB, error) {
	path := os.Getenv("CSIT214_DB_PATH")
	if path == "" {
		path = defaultDBPath
	}

	database, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}
	if err := database.Ping(); err != nil {
		database.Close()
		return nil, err
	}
	return database, nil
}
