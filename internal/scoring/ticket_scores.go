package scoring

import (
	"context"
	"time"

	"ticket-score-engine/internal/domain"
	"ticket-score-engine/internal/repository"
)

type TicketScorer struct {
	repo repository.TicketRepository
}

func NewTicketScorer(repo repository.TicketRepository) *TicketScorer {
	return &TicketScorer{repo: repo}
}

func (s *TicketScorer) GetTicketScores(ctx context.Context, start, end time.Time) ([]domain.TicketCategoryScore, error) {
	return s.repo.GetScoresByTicket(ctx, start, end)
}
