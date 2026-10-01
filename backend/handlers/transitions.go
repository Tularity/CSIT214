package handlers

import (
	"context"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/Tularity/CSIT214/backend/db"
	"github.com/Tularity/CSIT214/backend/models"
)

func incidentID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "Incident id must be a positive number")
		return 0, false
	}
	return id, true
}

func writeTransitionResult(w http.ResponseWriter, incident models.Incident, err error) {
	var conflict *db.ConflictError
	var missing *db.NotFoundError
	switch {
	case err == nil:
		writeJSON(w, http.StatusOK, incident)
	case errors.As(err, &conflict):
		writeError(w, http.StatusConflict, conflict.Message)
	case errors.As(err, &missing):
		writeError(w, http.StatusNotFound, missing.Error())
	default:
		log.Printf("update incident: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not update the incident")
	}
}

func (a *API) simpleTransition(change func(ctx context.Context, actor string, id int64) (models.Incident, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := incidentID(w, r)
		if !ok {
			return
		}
		incident, err := change(r.Context(), actorFrom(r), id)
		writeTransitionResult(w, incident, err)
	}
}

func (a *API) AssessIncident(w http.ResponseWriter, r *http.Request) {
	a.simpleTransition(a.store.AssessIncident)(w, r)
}

func (a *API) ResolveIncident(w http.ResponseWriter, r *http.Request) {
	a.simpleTransition(a.store.ResolveIncident)(w, r)
}

func (a *API) AssignShelter(w http.ResponseWriter, r *http.Request) {
	id, ok := incidentID(w, r)
	if !ok {
		return
	}
	var body struct {
		ShelterID int64 `json:"shelter_id"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	if body.ShelterID < 1 {
		writeError(w, http.StatusBadRequest, "shelter_id must be a positive number")
		return
	}

	incident, err := a.store.AssignShelter(r.Context(), actorFrom(r), id, body.ShelterID)
	writeTransitionResult(w, incident, err)
}

func (a *API) AssignResponder(w http.ResponseWriter, r *http.Request) {
	id, ok := incidentID(w, r)
	if !ok {
		return
	}
	var body struct {
		ResponderID int64 `json:"responder_id"`
	}
	if !decodeJSON(w, r, &body) {
		return
	}
	if body.ResponderID < 1 {
		writeError(w, http.StatusBadRequest, "responder_id must be a positive number")
		return
	}

	incident, err := a.store.AssignResponder(r.Context(), actorFrom(r), id, body.ResponderID)
	writeTransitionResult(w, incident, err)
}
