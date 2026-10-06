package handlers

import (
	"log"
	"net/http"
)

func (a *API) Dashboard(w http.ResponseWriter, r *http.Request) {
	dashboard, err := a.store.Dashboard(r.Context())
	if err != nil {
		log.Printf("dashboard: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not load the dashboard")
		return
	}
	writeJSON(w, http.StatusOK, dashboard)
}
