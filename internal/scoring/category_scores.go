package scoring

import (
	"context"
	"time"

	"ticket-score-engine/internal/domain"
	"ticket-score-engine/internal/repository"
)

type CategoryScorer struct {
	repo repository.CategoryRepository
}

func NewCategoryScorer(repo repository.CategoryRepository) *CategoryScorer {
	return &CategoryScorer{repo: repo}
}

func (s *CategoryScorer) GetCategoryScores(ctx context.Context, start, end time.Time) ([]domain.CategoryScore, error) {
	return s.repo.GetCategoryScores(ctx, start, end)
}
