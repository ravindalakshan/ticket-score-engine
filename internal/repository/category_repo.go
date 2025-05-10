package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ticket-score-engine/internal/domain"
)

type CategoryRepository interface {
	GetCategoryScores(ctx context.Context, start, end time.Time) ([]domain.CategoryScore, error)
}

type categoryRepo struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) CategoryRepository {
	return &categoryRepo{db: db}
}

func (r *categoryRepo) GetCategoryScores(ctx context.Context, start, end time.Time) ([]domain.CategoryScore, error) {
	isWeekly := end.Sub(start) > 30*24*time.Hour

	var query string
	if isWeekly {
		query = `
			SELECT 
				rc.name AS category,
				STRFTIME('%Y-%V', r.created_at) as period,
				COUNT(r.id) as count,
				SUM((r.rating * 1.0 / 5.0) * rc.weight) as weighted_score,
				SUM(rc.weight) as total_weight
			FROM ratings r
			JOIN rating_categories rc ON r.rating_category_id = rc.id
			WHERE r.created_at BETWEEN ? AND ?
			GROUP BY rc.name, STRFTIME('%Y-%V', r.created_at)
			ORDER BY rc.name, period`
	} else {
		query = `
			SELECT 
				rc.name AS category,
				DATE(r.created_at) as period,
				COUNT(r.id) as count,
				SUM((r.rating * 1.0 / 5.0) * rc.weight) as weighted_score,
				SUM(rc.weight) as total_weight
			FROM ratings r
			JOIN rating_categories rc ON r.rating_category_id = rc.id
			WHERE r.created_at BETWEEN ? AND ?
			GROUP BY rc.name, DATE(r.created_at)
			ORDER BY rc.name, period`
	}

	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to query category scores: %w", err)
	}
	defer rows.Close()

	var scores []domain.CategoryScore
	for rows.Next() {
		var cs domain.CategoryScore
		var weightedSum, totalWeight float64

		if err := rows.Scan(
			&cs.CategoryName,
			&cs.Date,
			&cs.RatingCount,
			&weightedSum,
			&totalWeight,
		); err != nil {
			return nil, fmt.Errorf("failed to scan category score: %w", err)
		}

		if totalWeight > 0 {
			cs.Score = (weightedSum / totalWeight) * 100
		}
		//cs.IsWeekly = isWeekly

		scores = append(scores, cs)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return scores, nil
}
