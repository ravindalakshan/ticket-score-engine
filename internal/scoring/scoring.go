package scoring

import (
	"database/sql"
	"fmt"
	"time"
)

type CategoryScore struct {
	CategoryName string
	Date         string // YYYY-MM-DD
	Score        float64
	RatingCount  int
}

// GetCategoryScores calculates weighted scores per category over a time range.
func GetCategoryScores(db *sql.DB, start, end time.Time) ([]CategoryScore, error) {
	// Determine if we should use weekly aggregation
	duration := end.Sub(start)
	useWeekly := duration > 30*24*time.Hour // More than 30 days

	var query string
	if useWeekly {
		query = `
            SELECT 
                rc.name AS category,
                STRFTIME('%Y-%V', r.created_at) as week,
                COUNT(r.id) as count,
                SUM((r.rating * 1.0 / 5.0) * rc.weight) as weighted_score,
                SUM(rc.weight) as total_weight
            FROM ratings r
            JOIN rating_categories rc ON r.rating_category_id = rc.id
            WHERE r.created_at BETWEEN ? AND ?
            GROUP BY rc.name, STRFTIME('%Y-%W', r.created_at)
            ORDER BY rc.name, week;
        `
	} else {
		query = `
             SELECT 
                rc.name AS category,
                DATE(r.created_at) as date,
                COUNT(r.id) as count,
                SUM((r.rating * 1.0 / 5.0) * rc.weight) as weighted_score,
                SUM(rc.weight) as total_weight
            FROM ratings r
            JOIN rating_categories rc ON r.rating_category_id = rc.id
            WHERE r.created_at BETWEEN ? AND ?
            GROUP BY rc.name, DATE(r.created_at)
            ORDER BY rc.name, DATE(r.created_at);
        `
	}

	rows, err := db.Query(query, start, end)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var results []CategoryScore

	for rows.Next() {
		var category string
		var date string
		var count int
		var weightedSum float64
		var totalWeight float64

		if err := rows.Scan(&category, &date, &count, &weightedSum, &totalWeight); err != nil {
			return nil, err
		}

		normalizedScore := (weightedSum / totalWeight) * 100

		results = append(results, CategoryScore{
			CategoryName: category,
			Date:         date,
			Score:        normalizedScore,
			RatingCount:  count,
		})
	}

	return results, nil
}
