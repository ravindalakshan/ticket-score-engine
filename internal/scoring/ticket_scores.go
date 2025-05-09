package scoring

import (
	"database/sql"
	"fmt"
	"time"
)

// TicketCategoryScore represents aggregated category score per ticket
type TicketCategoryScore struct {
	TicketID     int
	CategoryName string
	Score        float64
}

// GetScoresByTicket returns per-ticket aggregated category scores within a date range
func GetScoresByTicket(db *sql.DB, start, end time.Time) ([]TicketCategoryScore, error) {
	query := `
        SELECT 
            r.ticket_id,
            rc.name AS category,
            SUM((r.rating * 1.0 / 5.0) * rc.weight) as weighted_score,
            SUM(rc.weight) as total_weight
        FROM ratings r
        JOIN rating_categories rc ON r.rating_category_id = rc.id
        WHERE r.created_at BETWEEN ? AND ?
        GROUP BY r.ticket_id, rc.name
        ORDER BY r.ticket_id, rc.name;
    `

	rows, err := db.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []TicketCategoryScore

	for rows.Next() {
		var ticketID int
		var category string
		var weightedSum, totalWeight float64

		if err := rows.Scan(&ticketID, &category, &weightedSum, &totalWeight); err != nil {
			return nil, err
		}

		normalizedScore := (weightedSum / totalWeight) * 100

		results = append(results, TicketCategoryScore{
			TicketID:     ticketID,
			CategoryName: category,
			Score:        normalizedScore,
		})
	}

	return results, nil
}
