package scoring

import (
	"database/sql"
	"fmt"
	"time"
)

// OverallScoreResult represents the result of an overall score calculation
type OverallScoreResult struct {
	Score       float64
	RatingCount int
}

// PeriodComparisonResult represents the comparison between two time periods
type PeriodComparisonResult struct {
	PercentageChange float64
	CurrentScore     float64
	PreviousScore    float64
	CurrentCount     int
	PreviousCount    int
}

// GetOverallScore calculates the weighted average score across all categories for a given time period
func GetOverallScore(db *sql.DB, start, end time.Time) (*OverallScoreResult, error) {
	query := `
        SELECT 
            SUM((r.rating * 1.0 / 5.0) * rc.weight) as total_weighted_score,
            SUM(rc.weight) as total_weight,
            COUNT(r.id) as rating_count
        FROM ratings r
        JOIN rating_categories rc ON r.rating_category_id = rc.id
        WHERE r.created_at BETWEEN ? AND ?;
    `

	var (
		totalWeightedScore float64
		totalWeight        float64
		ratingCount        int
	)

	err := db.QueryRow(query, start, end).Scan(&totalWeightedScore, &totalWeight, &ratingCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return &OverallScoreResult{Score: 0, RatingCount: 0}, nil
		}
		return nil, fmt.Errorf("failed to calculate overall score: %w", err)
	}

	if totalWeight == 0 {
		return &OverallScoreResult{Score: 0, RatingCount: ratingCount}, nil
	}

	overallScore := (totalWeightedScore / totalWeight) * 100
	return &OverallScoreResult{
		Score:       overallScore,
		RatingCount: ratingCount,
	}, nil
}

// GetPeriodComparison compares the overall scores between two time periods
func GetPeriodComparison(db *sql.DB, currentStart, currentEnd, previousStart, previousEnd time.Time) (*PeriodComparisonResult, error) {
	// Get current period score
	currentResult, err := GetOverallScore(db, currentStart, currentEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get current period score: %w", err)
	}

	// Get previous period score
	previousResult, err := GetOverallScore(db, previousStart, previousEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous period score: %w", err)
	}

	// Calculate percentage change
	var change float64
	if previousResult.Score != 0 {
		change = ((currentResult.Score - previousResult.Score) / previousResult.Score) * 100
	}

	return &PeriodComparisonResult{
		PercentageChange: change,
		CurrentScore:     currentResult.Score,
		PreviousScore:    previousResult.Score,
		CurrentCount:     currentResult.RatingCount,
		PreviousCount:    previousResult.RatingCount,
	}, nil
}

// Helper function to calculate time ranges for common period comparisons
func GetComparisonPeriods(period string) (time.Time, time.Time, time.Time, time.Time, error) {
	now := time.Now()
	var currentStart, currentEnd, previousStart, previousEnd time.Time

	switch period {
	case "week":
		currentStart = now.AddDate(0, 0, -7)
		currentEnd = now
		previousStart = now.AddDate(0, 0, -14)
		previousEnd = now.AddDate(0, 0, -7)
	case "month":
		currentStart = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
		currentEnd = now
		prevMonth := now.AddDate(0, -1, 0)
		previousStart = time.Date(prevMonth.Year(), prevMonth.Month(), 1, 0, 0, 0, 0, time.UTC)
		previousEnd = time.Date(prevMonth.Year(), prevMonth.Month()+1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Nanosecond)
	default:
		return time.Time{}, time.Time{}, time.Time{}, time.Time{}, fmt.Errorf("invalid period: %s", period)
	}

	return currentStart, currentEnd, previousStart, previousEnd, nil
}
