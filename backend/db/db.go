package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

//go:embed schema.sql
var schema string

const DefaultPath = "data.db"

type Store struct {
	db *sql.DB
}

// Open creates the schema if needed and, when the database is empty, loads
// seedPath; an empty seedPath means the usual locations are searched.
func Open(path, seedPath string) (*Store, error) {
	if path == "" {
		path = DefaultPath
	}

	// _txlock=immediate makes a write transaction take the lock when it starts,
	// so two concurrent writers wait on busy_timeout instead of one failing
	// half-way through when it tries to upgrade a read lock.
	database, err := sql.Open("sqlite", path+"?_pragma=foreign_keys(1)&_pragma=busy_timeout(5000)&_txlock=immediate")
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", path, err)
	}

	if err := database.Ping(); err != nil {
		database.Close()
		return nil, fmt.Errorf("ping %s: %w", path, err)
	}

	if _, err := database.Exec(schema); err != nil {
		database.Close()
		return nil, fmt.Errorf("apply schema: %w", err)
	}

	if err := seedIfEmpty(database, seedPath); err != nil {
		database.Close()
		return nil, err
	}

	return &Store{db: database}, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

func seedIfEmpty(database *sql.DB, seedPath string) error {
	var incidents int
	if err := database.QueryRow(`SELECT COUNT(*) FROM incidents`).Scan(&incidents); err != nil {
		return fmt.Errorf("count incidents: %w", err)
	}
	if incidents > 0 {
		return nil
	}

	path, err := findSeedFile(seedPath)
	if err != nil {
		log.Printf("skipping seed data: %v", err)
		return nil
	}

	statements, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read seed file %s: %w", path, err)
	}
	if _, err := database.Exec(string(statements)); err != nil {
		return fmt.Errorf("load seed file %s: %w", path, err)
	}

	log.Printf("loaded seed data from %s", path)
	return nil
}

// findSeedFile looks beside the working directory first so that both `go run .`
// from backend and a binary started from the repository root find the same file.
func findSeedFile(configured string) (string, error) {
	if configured != "" {
		if _, err := os.Stat(configured); err != nil {
			return "", fmt.Errorf("seed file %s is not readable: %w", configured, err)
		}
		return configured, nil
	}

	candidates := []string{
		filepath.Join("seed", "seed.sql"),
		filepath.Join("..", "seed", "seed.sql"),
	}
	for _, candidate := range candidates {
		if _, err := os.Stat(candidate); err == nil {
			return candidate, nil
		}
	}
	return "", fmt.Errorf("no seed file found in %v", candidates)
}
