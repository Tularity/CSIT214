package db

import (
	"context"
	"fmt"

	"github.com/Tularity/CSIT214/backend/models"
)

func percent(part, whole int) int {
	if whole == 0 {
		return 0
	}
	return (part*100 + whole/2) / whole
}

func (s *Store) countBy(ctx context.Context, query string, counts map[string]int) error {
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var key string
		var count int
		if err := rows.Scan(&key, &count); err != nil {
			return err
		}
		counts[key] = count
	}
	return rows.Err()
}

// Dashboard zero-fills every status and band so the client can draw a fixed
// set of tiles without guessing which keys might be missing.
func (s *Store) Dashboard(ctx context.Context) (models.Dashboard, error) {
	dashboard := models.Dashboard{
		IncidentsByStatus: map[string]int{
			statusRegistered: 0, statusAssessed: 0, statusAssigned: 0, statusInProgress: 0, statusResolved: 0,
		},
		IncidentsByBand: map[string]int{
			models.BandCritical: 0, models.BandHigh: 0, models.BandMedium: 0, models.BandLow: 0,
		},
		RespondersByStatus: map[string]int{"available": 0, "assigned": 0, "off_duty": 0},
		Shelters:           []models.ShelterOccupancy{},
	}

	if err := s.countBy(ctx, `SELECT status, COUNT(*) FROM incidents GROUP BY status`,
		dashboard.IncidentsByStatus); err != nil {
		return models.Dashboard{}, fmt.Errorf("count incidents by status: %w", err)
	}
	if err := s.countBy(ctx, `SELECT priority_band, COUNT(*) FROM incidents GROUP BY priority_band`,
		dashboard.IncidentsByBand); err != nil {
		return models.Dashboard{}, fmt.Errorf("count incidents by band: %w", err)
	}
	if err := s.countBy(ctx, `SELECT status, COUNT(*) FROM responders GROUP BY status`,
		dashboard.RespondersByStatus); err != nil {
		return models.Dashboard{}, fmt.Errorf("count responders by status: %w", err)
	}

	for status, count := range dashboard.IncidentsByStatus {
		dashboard.IncidentsTotal += count
		if status != statusResolved {
			dashboard.IncidentsOpen += count
		}
	}
	dashboard.RespondersAvailable = dashboard.RespondersByStatus["available"]

	shelters, err := s.ListShelters(ctx)
	if err != nil {
		return models.Dashboard{}, err
	}
	for _, shelter := range shelters {
		dashboard.Shelters = append(dashboard.Shelters, models.ShelterOccupancy{
			ID:               shelter.ID,
			Name:             shelter.Name,
			CapacityTotal:    shelter.CapacityTotal,
			CapacityUsed:     shelter.CapacityUsed,
			OccupancyPercent: percent(shelter.CapacityUsed, shelter.CapacityTotal),
		})
		dashboard.CapacityTotal += shelter.CapacityTotal
		dashboard.CapacityUsed += shelter.CapacityUsed
	}
	dashboard.OccupancyPercent = percent(dashboard.CapacityUsed, dashboard.CapacityTotal)

	return dashboard, nil
}
