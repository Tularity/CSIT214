package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tularity/CSIT214/backend/db"
	"github.com/Tularity/CSIT214/backend/handlers"
)

const defaultAddr = ":8080"

func main() {
	store, err := db.Open(os.Getenv("CSIT214_DB_PATH"))
	if err != nil {
		log.Fatalf("database unavailable: %v", err)
	}
	defer store.Close()

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
