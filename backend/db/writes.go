package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Tularity/CSIT214/backend/models"
)

func timestamp(at time.Time) string {
	return at.UTC().Truncate(time.Second).Format(time.RFC3339)
}

// inTransaction commits only when work succeeds, so a change and its audit
// row are either both stored or neither is.
func (s *Store) inTransaction(ctx context.Context, work func(tx *sql.Tx) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	if err := work(tx); err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func recordAudit(ctx context.Context, tx *sql.Tx, at, actor, action, entityType string, entityID int64, detail string) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO audit (at, actor, action, entity_type, entity_id, detail)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		at, actor, action, entityType, entityID, detail)
	if err != nil {
		return fmt.Errorf("record audit for %s %d: %w", entityType, entityID, err)
	}
	return nil
}

func (s *Store) CreateIncident(ctx context.Context, actor string, input models.NewIncident) (models.Incident, error) {
	score := models.PriorityScore(input.Severity, input.PeopleAffected, input.Vulnerable)
	band := models.PriorityBand(score)
	reportedAt := timestamp(time.Now())

	var id int64
	err := s.inTransaction(ctx, func(tx *sql.Tx) error {
		result, err := tx.ExecContext(ctx,
			`INSERT INTO incidents
			   (type, location, description, people_affected, vulnerable, severity,
			    priority_score, priority_band, status, reported_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, 'registered', ?)`,
			input.Type, input.Location, input.Description, input.PeopleAffected,
			input.Vulnerable, input.Severity, score, band, reportedAt)
		if err != nil {
			return fmt.Errorf("insert incident: %w", err)
		}
		if id, err = result.LastInsertId(); err != nil {
			return fmt.Errorf("read new incident id: %w", err)
		}
		return recordAudit(ctx, tx, reportedAt, actor, "register", "incident", id,
			fmt.Sprintf("Registered with priority score %d (%s)", score, band))
	})
	if err != nil {
		return models.Incident{}, err
	}
	return s.Incident(ctx, id)
}
