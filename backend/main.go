package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Tularity/CSIT214/backend/db"
	"github.com/Tularity/CSIT214/backend/handlers"
	"github.com/Tularity/CSIT214/backend/web"
)

const defaultPort = "8080"

func main() {
	locations := resolvePaths()
	store, err := db.Open(locations.database, locations.seed)
	if err != nil {
		log.Fatalf("database unavailable: %v", err)
	}
	defer store.Close()

	mux := http.NewServeMux()
	mux.Handle("/api/", handlers.New(store).Routes())
	mux.Handle("/", web.Handler())

	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}
	server := &http.Server{
		Addr:              ":" + port,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	// Listening before serving lets a busy port be reported in plain words,
	// which matters most to someone who started the program by double-clicking.
	listener, err := net.Listen("tcp", server.Addr)
	if err != nil {
		log.Fatalf("cannot listen on port %s (%v); close whatever is using it or set PORT to another number", port, err)
	}

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	go func() {
		if err := server.Serve(listener); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server stopped: %v", err)
		}
	}()

	address := "http://localhost:" + port
	log.Printf("listening on %s", server.Addr)
	log.Printf("open %s in a browser; press Ctrl+C to stop", address)
	openBrowser(address)

	<-shutdown
	log.Print("shutting down")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("forced shutdown: %v", err)
	}
}
