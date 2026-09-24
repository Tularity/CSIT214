package handlers

import (
	"log"
	"net/http"
)

func (a *API) ListShelters(w http.ResponseWriter, r *http.Request) {
	shelters, err := a.store.ListShelters(r.Context())
	if err != nil {
		log.Printf("list shelters: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not load shelters")
		return
	}
	writeJSON(w, http.StatusOK, shelters)
}

func (a *API) ListResponders(w http.ResponseWriter, r *http.Request) {
	responders, err := a.store.ListResponders(r.Context())
	if err != nil {
		log.Printf("list responders: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not load responders")
		return
	}
	writeJSON(w, http.StatusOK, responders)
}
