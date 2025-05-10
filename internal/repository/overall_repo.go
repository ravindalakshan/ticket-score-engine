package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type OverallRepository interface {
	GetOverallScore(ctx context.Context, start, end time.Time) (float64, int, error)
}

type overallRepo struct {
	db *sql.DB
}

func NewOverallRepository(db *sql.DB) OverallRepository {
	return &overallRepo{db: db}
}

func (r *overallRepo) GetOverallScore(ctx context.Context, start, end time.Time) (float64, int, error) {
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

	err := r.db.QueryRowContext(ctx, query, start, end).Scan(&totalWeightedScore, &totalWeight, &ratingCount)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, 0, nil
		}
		return 0, 0, fmt.Errorf("query error: %w", err)
	}

	if totalWeight == 0 {
		return 0, ratingCount, nil
	}

	score := (totalWeightedScore / totalWeight) * 100
	return score, ratingCount, nil
}
