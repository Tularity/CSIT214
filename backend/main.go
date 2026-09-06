package main

import (
	"database/sql"
	"log"
	"os"

	_ "modernc.org/sqlite"
)

const defaultDBPath = "data.db"

func main() {
	path := os.Getenv("CSIT214_DB_PATH")
	if path == "" {
		path = defaultDBPath
	}

	database, err := sql.Open("sqlite", path)
	if err != nil {
		log.Fatalf("open %s: %v", path, err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("ping %s: %v", path, err)
	}

	var version string
	if err := database.QueryRow(`SELECT sqlite_version()`).Scan(&version); err != nil {
		log.Fatalf("query version: %v", err)
	}
	log.Printf("sqlite %s ready at %s", version, path)
}
