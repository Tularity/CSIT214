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
	mux.HandleFunc("GET /api/incidents/{id}", a.GetIncident)
	mux.HandleFunc("GET /api/shelters", a.ListShelters)
	mux.HandleFunc("GET /api/responders", a.ListResponders)
	return mux
}
