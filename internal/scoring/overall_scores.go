package scoring

import (
	"database/sql"
	"fmt"
	"time"
)

// GetOverallScore computes the total weighted average score for the period.
func GetOverallScore(db *sql.DB, start, end time.Time) (float64, error) {
	query := `
        SELECT 
            SUM((r.rating * 1.0 / 5.0) * rc.weight) as weighted_score,
            SUM(rc.weight) as total_weight
        FROM ratings r
        JOIN rating_categories rc ON r.rating_category_id = rc.id
        WHERE r.created_at BETWEEN ? AND ?;
    `

	var weightedSum, totalWeight float64
	err := db.QueryRow(query, start, end).Scan(&weightedSum, &totalWeight)
	if err != nil {
		return 0, fmt.Errorf("failed to compute overall score: %w", err)
	}

	if totalWeight == 0 {
		return 0, nil // avoid division by zero
	}

	normalizedScore := (weightedSum / totalWeight) * 100
	return normalizedScore, nil
}
