package handlers

import (
	"net/http"
	"strings"
	"unicode/utf8"
)

const (
	defaultActor   = "duty.officer"
	maxActorLength = 60
)

// actorFrom names the operator for the audit trail. The prototype has no
// sign-in, so a client may pass the name in X-Actor; otherwise the entry is
// attributed to the duty officer seat.
func actorFrom(r *http.Request) string {
	actor := strings.TrimSpace(r.Header.Get("X-Actor"))
	if actor == "" {
		return defaultActor
	}
	if utf8.RuneCountInString(actor) > maxActorLength {
		actor = string([]rune(actor)[:maxActorLength])
	}
	return actor
}
