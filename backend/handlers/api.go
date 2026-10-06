package handlers

import (
	"net/http"

	"github.com/Tularity/CSIT214/backend/db"
)

type API struct {
	store *db.Store
}

func New(store *db.Store) *API {
	return &API{store: store}
}

func (a *API) Routes() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/health", a.Health)
	mux.HandleFunc("GET /api/incidents", a.ListIncidents)
	mux.HandleFunc("POST /api/incidents", a.CreateIncident)
	mux.HandleFunc("GET /api/incidents/{id}", a.GetIncident)
	mux.HandleFunc("POST /api/incidents/{id}/assess", a.AssessIncident)
	mux.HandleFunc("POST /api/incidents/{id}/assign-shelter", a.AssignShelter)
	mux.HandleFunc("POST /api/incidents/{id}/assign-responder", a.AssignResponder)
	mux.HandleFunc("POST /api/incidents/{id}/resolve", a.ResolveIncident)
	mux.HandleFunc("GET /api/shelters", a.ListShelters)
	mux.HandleFunc("GET /api/responders", a.ListResponders)
	mux.HandleFunc("GET /api/audit", a.ListAudit)
	mux.HandleFunc("GET /api/dashboard", a.Dashboard)
	return mux
}
