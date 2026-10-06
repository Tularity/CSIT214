package db

import (
	"context"
	"fmt"
	"strings"

	"github.com/Tularity/CSIT214/backend/models"
)

type AuditFilter struct {
	EntityType string
	EntityID   int64
	Limit      int
}

func (s *Store) ListAudit(ctx context.Context, filter AuditFilter) ([]models.AuditEntry, error) {
	var conditions []string
	var arguments []any
	if filter.EntityType != "" {
		conditions = append(conditions, "entity_type = ?")
		arguments = append(arguments, filter.EntityType)
	}
	if filter.EntityID > 0 {
		conditions = append(conditions, "entity_id = ?")
		arguments = append(arguments, filter.EntityID)
	}

	query := `SELECT id, at, actor, action, entity_type, entity_id, detail FROM audit`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY at DESC, id DESC LIMIT ?"
	arguments = append(arguments, filter.Limit)

	rows, err := s.db.QueryContext(ctx, query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("list audit: %w", err)
	}
	defer rows.Close()

	entries := []models.AuditEntry{}
	for rows.Next() {
		var entry models.AuditEntry
		if err := rows.Scan(&entry.ID, &entry.At, &entry.Actor, &entry.Action,
			&entry.EntityType, &entry.EntityID, &entry.Detail); err != nil {
			return nil, fmt.Errorf("list audit: %w", err)
		}
		entries = append(entries, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("list audit: %w", err)
	}
	return entries, nil
}
