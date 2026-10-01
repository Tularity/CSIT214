package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/Tularity/CSIT214/backend/db"
	"github.com/Tularity/CSIT214/backend/models"
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
	id, ok := incidentID(w, r)
	if !ok {
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

func (a *API) CreateIncident(w http.ResponseWriter, r *http.Request) {
	var input models.NewIncident
	if !decodeJSON(w, r, &input) {
		return
	}

	input = input.Normalized()
	if err := input.Validate(); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	incident, err := a.store.CreateIncident(r.Context(), actorFrom(r), input)
	if err != nil {
		log.Printf("create incident: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not register the incident")
		return
	}

	w.Header().Set("Location", fmt.Sprintf("/api/incidents/%d", incident.ID))
	writeJSON(w, http.StatusCreated, incident)
}
