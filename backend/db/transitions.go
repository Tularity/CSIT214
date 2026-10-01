package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Tularity/CSIT214/backend/models"
)

const (
	statusRegistered = "registered"
	statusAssessed   = "assessed"
	statusAssigned   = "assigned"
	statusInProgress = "in_progress"
	statusResolved   = "resolved"
)

// checkStep enforces registered → assessed → assigned → in_progress → resolved.
// A shelter is optional (a medical call may need none), so dispatching and
// resolving are allowed straight after assessment, but nothing may skip it.
func checkStep(incident models.Incident, otherwise string, allowed ...string) error {
	for _, status := range allowed {
		if incident.Status == status {
			return nil
		}
	}
	switch incident.Status {
	case statusRegistered:
		return &ConflictError{"Incident must be assessed first"}
	case statusResolved:
		return &ConflictError{"Incident is already resolved"}
	}
	return &ConflictError{otherwise}
}

func incidentInTx(ctx context.Context, tx *sql.Tx, id int64) (models.Incident, error) {
	row := tx.QueryRowContext(ctx, `SELECT `+incidentColumns+` FROM incidents WHERE id = ?`, id)
	incident, err := scanIncident(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Incident{}, &NotFoundError{"Incident"}
	}
	if err != nil {
		return models.Incident{}, fmt.Errorf("load incident %d: %w", id, err)
	}
	return incident, nil
}

func (s *Store) transition(ctx context.Context, id int64, change func(tx *sql.Tx, incident models.Incident, at string) error) (models.Incident, error) {
	at := timestamp(time.Now())
	err := s.inTransaction(ctx, func(tx *sql.Tx) error {
		incident, err := incidentInTx(ctx, tx, id)
		if err != nil {
			return err
		}
		return change(tx, incident, at)
	})
	if err != nil {
		return models.Incident{}, err
	}
	return s.Incident(ctx, id)
}

func setStatus(ctx context.Context, tx *sql.Tx, id int64, status string) error {
	if _, err := tx.ExecContext(ctx, `UPDATE incidents SET status = ? WHERE id = ?`, status, id); err != nil {
		return fmt.Errorf("set incident %d to %s: %w", id, status, err)
	}
	return nil
}

func (s *Store) AssessIncident(ctx context.Context, actor string, id int64) (models.Incident, error) {
	return s.transition(ctx, id, func(tx *sql.Tx, incident models.Incident, at string) error {
		if err := checkStep(incident, "Incident has already been assessed", statusRegistered); err != nil {
			return err
		}
		if err := setStatus(ctx, tx, id, statusAssessed); err != nil {
			return err
		}
		return recordAudit(ctx, tx, at, actor, "assess", "incident", id,
			fmt.Sprintf("Assessed at priority score %d (%s)", incident.PriorityScore, incident.PriorityBand))
	})
}

func (s *Store) AssignShelter(ctx context.Context, actor string, id, shelterID int64) (models.Incident, error) {
	return s.transition(ctx, id, func(tx *sql.Tx, incident models.Incident, at string) error {
		otherwise := "A shelter can only be assigned before a responder is dispatched"
		if incident.ShelterID != nil {
			otherwise = "Incident already has a shelter"
		}
		if err := checkStep(incident, otherwise, statusAssessed); err != nil {
			return err
		}

		var name string
		var total, used int
		err := tx.QueryRowContext(ctx,
			`SELECT name, capacity_total, capacity_used FROM shelters WHERE id = ?`, shelterID).
			Scan(&name, &total, &used)
		if errors.Is(err, sql.ErrNoRows) {
			return &NotFoundError{"Shelter"}
		}
		if err != nil {
			return fmt.Errorf("load shelter %d: %w", shelterID, err)
		}

		if remaining := total - used; remaining < incident.PeopleAffected {
			return &ConflictError{fmt.Sprintf("Shelter has only %d places remaining", remaining)}
		}

		if _, err := tx.ExecContext(ctx,
			`UPDATE shelters SET capacity_used = capacity_used + ? WHERE id = ?`,
			incident.PeopleAffected, shelterID); err != nil {
			return fmt.Errorf("reserve places at shelter %d: %w", shelterID, err)
		}
		if _, err := tx.ExecContext(ctx,
			`UPDATE incidents SET shelter_id = ?, status = ? WHERE id = ?`,
			shelterID, statusAssigned, id); err != nil {
			return fmt.Errorf("assign shelter to incident %d: %w", id, err)
		}
		return recordAudit(ctx, tx, at, actor, "assign-shelter", "incident", id,
			fmt.Sprintf("Assigned shelter %d (%s) for %d people", shelterID, name, incident.PeopleAffected))
	})
}

func (s *Store) AssignResponder(ctx context.Context, actor string, id, responderID int64) (models.Incident, error) {
	return s.transition(ctx, id, func(tx *sql.Tx, incident models.Incident, at string) error {
		if err := checkStep(incident, "Incident already has a responder", statusAssessed, statusAssigned); err != nil {
			return err
		}

		var name, status string
		err := tx.QueryRowContext(ctx,
			`SELECT name, status FROM responders WHERE id = ?`, responderID).Scan(&name, &status)
		if errors.Is(err, sql.ErrNoRows) {
			return &NotFoundError{"Responder"}
		}
		if err != nil {
			return fmt.Errorf("load responder %d: %w", responderID, err)
		}
		if status != "available" {
			return &ConflictError{"Responder is not available"}
		}

		if _, err := tx.ExecContext(ctx,
			`UPDATE responders SET status = 'assigned' WHERE id = ?`, responderID); err != nil {
			return fmt.Errorf("mark responder %d assigned: %w", responderID, err)
		}
		if _, err := tx.ExecContext(ctx,
			`UPDATE incidents SET responder_id = ?, status = ? WHERE id = ?`,
			responderID, statusInProgress, id); err != nil {
			return fmt.Errorf("assign responder to incident %d: %w", id, err)
		}
		return recordAudit(ctx, tx, at, actor, "assign-responder", "incident", id,
			fmt.Sprintf("Assigned responder %d (%s)", responderID, name))
	})
}

func (s *Store) ResolveIncident(ctx context.Context, actor string, id int64) (models.Incident, error) {
	return s.transition(ctx, id, func(tx *sql.Tx, incident models.Incident, at string) error {
		if err := checkStep(incident, "Incident cannot be resolved yet",
			statusAssessed, statusAssigned, statusInProgress); err != nil {
			return err
		}

		released := []string{"Resolved"}
		if incident.ShelterID != nil {
			if _, err := tx.ExecContext(ctx,
				`UPDATE shelters SET capacity_used = capacity_used - ? WHERE id = ?`,
				incident.PeopleAffected, *incident.ShelterID); err != nil {
				return fmt.Errorf("release places at shelter %d: %w", *incident.ShelterID, err)
			}
			released = append(released,
				fmt.Sprintf("%d places released at shelter %d", incident.PeopleAffected, *incident.ShelterID))
		}
		if incident.ResponderID != nil {
			// Only an assigned responder goes back to available; one who went
			// off duty while on the job stays off duty.
			if _, err := tx.ExecContext(ctx,
				`UPDATE responders SET status = 'available' WHERE id = ? AND status = 'assigned'`,
				*incident.ResponderID); err != nil {
				return fmt.Errorf("release responder %d: %w", *incident.ResponderID, err)
			}
			released = append(released,
				fmt.Sprintf("responder %d returned to available", *incident.ResponderID))
		}

		if err := setStatus(ctx, tx, id, statusResolved); err != nil {
			return err
		}
		return recordAudit(ctx, tx, at, actor, "resolve", "incident", id, strings.Join(released, "; "))
	})
}
