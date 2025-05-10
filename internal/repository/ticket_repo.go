package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"ticket-score-engine/internal/domain"
)

type TicketRepository interface {
	GetScoresByTicket(ctx context.Context, start, end time.Time) ([]domain.TicketCategoryScore, error)
}

type ticketRepo struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepo{db: db}
}

func (r *ticketRepo) GetScoresByTicket(ctx context.Context, start, end time.Time) ([]domain.TicketCategoryScore, error) {
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

	rows, err := r.db.QueryContext(ctx, query, start, end)
	if err != nil {
		return nil, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	var scores []domain.TicketCategoryScore
	for rows.Next() {
		var score domain.TicketCategoryScore
		var weightedSum, totalWeight float64

		if err := rows.Scan(&score.TicketID, &score.CategoryName, &weightedSum, &totalWeight); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if totalWeight > 0 {
			score.Score = (weightedSum / totalWeight) * 100
		}
		scores = append(scores, score)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return scores, nil
}
