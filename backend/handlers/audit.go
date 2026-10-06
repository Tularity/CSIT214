package handlers

import (
	"log"
	"net/http"
	"strconv"

	"github.com/Tularity/CSIT214/backend/db"
)

const (
	defaultAuditLimit = 200
	maxAuditLimit     = 1000
)

func (a *API) ListAudit(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	filter := db.AuditFilter{
		EntityType: query.Get("entity_type"),
		Limit:      defaultAuditLimit,
	}

	if raw := query.Get("entity_id"); raw != "" {
		id, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || id < 1 {
			writeError(w, http.StatusBadRequest, "entity_id must be a positive number")
			return
		}
		filter.EntityID = id
	}
	if raw := query.Get("limit"); raw != "" {
		limit, err := strconv.Atoi(raw)
		if err != nil || limit < 1 || limit > maxAuditLimit {
			writeError(w, http.StatusBadRequest, "limit must be between 1 and 1000")
			return
		}
		filter.Limit = limit
	}

	entries, err := a.store.ListAudit(r.Context(), filter)
	if err != nil {
		log.Printf("list audit: %v", err)
		writeError(w, http.StatusInternalServerError, "Could not load the activity log")
		return
	}
	writeJSON(w, http.StatusOK, entries)
}
