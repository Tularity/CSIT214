package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Tularity/CSIT214/backend/models"
)

var ErrNotFound = errors.New("record not found")

const incidentColumns = `id, type, location, description, people_affected,
	vulnerable, severity, priority_score, priority_band, status,
	shelter_id, responder_id, reported_at`

type scanner interface {
	Scan(dest ...any) error
}

func (s *Store) ListIncidents(ctx context.Context) ([]models.Incident, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT `+incidentColumns+` FROM incidents ORDER BY priority_score DESC, id`)
	if err != nil {
		return nil, fmt.Errorf("list incidents: %w", err)
	}
	defer rows.Close()

	incidents := []models.Incident{}
	for rows.Next() {
		incident, err := scanIncident(rows)
		if err != nil {
			return nil, fmt.Errorf("list incidents: %w", err)
		}
		incidents = append(incidents, incident)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list incidents: %w", err)
	}
	return incidents, nil
}

func (s *Store) Incident(ctx context.Context, id int64) (models.Incident, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT `+incidentColumns+` FROM incidents WHERE id = ?`, id)

	incident, err := scanIncident(row)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Incident{}, ErrNotFound
	}
	if err != nil {
		return models.Incident{}, fmt.Errorf("incident %d: %w", id, err)
	}
	return incident, nil
}

func (s *Store) ListShelters(ctx context.Context) ([]models.Shelter, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, address, capacity_total, capacity_used FROM shelters ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list shelters: %w", err)
	}
	defer rows.Close()

	shelters := []models.Shelter{}
	for rows.Next() {
		var shelter models.Shelter
		if err := rows.Scan(&shelter.ID, &shelter.Name, &shelter.Address,
			&shelter.CapacityTotal, &shelter.CapacityUsed); err != nil {
			return nil, fmt.Errorf("list shelters: %w", err)
		}
		shelter.CapacityRemaining = shelter.CapacityTotal - shelter.CapacityUsed
		shelters = append(shelters, shelter)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list shelters: %w", err)
	}
	return shelters, nil
}

func (s *Store) ListResponders(ctx context.Context) ([]models.Responder, error) {
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, name, skill, status FROM responders ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("list responders: %w", err)
	}
	defer rows.Close()

	responders := []models.Responder{}
	for rows.Next() {
		var responder models.Responder
		if err := rows.Scan(&responder.ID, &responder.Name,
			&responder.Skill, &responder.Status); err != nil {
			return nil, fmt.Errorf("list responders: %w", err)
		}
		responders = append(responders, responder)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list responders: %w", err)
	}
	return responders, nil
}

func scanIncident(row scanner) (models.Incident, error) {
	var (
		incident    models.Incident
		vulnerable  int
		shelterID   sql.NullInt64
		responderID sql.NullInt64
	)

	err := row.Scan(
		&incident.ID, &incident.Type, &incident.Location, &incident.Description,
		&incident.PeopleAffected, &vulnerable, &incident.Severity,
		&incident.PriorityScore, &incident.PriorityBand, &incident.Status,
		&shelterID, &responderID, &incident.ReportedAt,
	)
	if err != nil {
		return models.Incident{}, err
	}

	incident.Vulnerable = vulnerable != 0
	if shelterID.Valid {
		incident.ShelterID = &shelterID.Int64
	}
	if responderID.Valid {
		incident.ResponderID = &responderID.Int64
	}
	return incident, nil
}
