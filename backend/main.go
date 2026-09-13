package main

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"github.com/Tularity/CSIT214/backend/handlers"
)

const (
	defaultAddr   = ":8080"
	defaultDBPath = "data.db"
)

func main() {
	database, err := openDatabase()
	if err != nil {
		log.Fatalf("database unavailable: %v", err)
	}
	defer database.Close()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", handlers.Health)

	server := &http.Server{
		Addr:              listenAddr(),
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("listening on %s", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	<-shutdown
	log.Print("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
}

func listenAddr() string {
	port := os.Getenv("PORT")
	if port == "" {
		return defaultAddr
	}
	return ":" + port
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
