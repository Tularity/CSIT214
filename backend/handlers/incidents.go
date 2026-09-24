package handlers

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Tularity/CSIT214/backend/db"
)

func (a *API) ListIncidents(w http.ResponseWriter, r *http.Request) {
	incidents, err := a.store.ListIncidents(r.Context())
	if err != nil {
		log.Printf("list incidents: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not load incidents")
		return
	}
	writeJSON(w, http.StatusOK, incidents)
}

func (a *API) GetIncident(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "Incident id must be a positive number")
		return
	}

	incident, err := a.store.Incident(r.Context(), id)
	if errors.Is(err, db.ErrNotFound) {
		writeError(w, http.StatusNotFound, "Incident not found")
		return
	}
	if err != nil {
		log.Printf("incident %d: %v", id, err)
		writeError(w, http.StatusInternalServerError, "Could not load the incident")
		return
	}
	writeJSON(w, http.StatusOK, incident)
}
