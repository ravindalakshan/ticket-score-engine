package scoring

import (
	"context"
	"fmt"
	"time"

	"ticket-score-engine/internal/domain"
	"ticket-score-engine/internal/repository"
)

type OverallScorer struct {
	repo repository.OverallRepository
}

func NewOverallScorer(repo repository.OverallRepository) *OverallScorer {
	return &OverallScorer{repo: repo}
}

func (s *OverallScorer) GetOverallScore(ctx context.Context, start, end time.Time) (*domain.OverallScoreResult, error) {
	score, count, err := s.repo.GetOverallScore(ctx, start, end)
	if err != nil {
		return nil, err
	}

	return &domain.OverallScoreResult{
		Score:       score,
		RatingCount: count,
	}, nil
}

func (s *OverallScorer) GetPeriodComparison(ctx context.Context, currentStart, currentEnd, previousStart, previousEnd time.Time) (*domain.PeriodComparisonResult, error) {
	current, err := s.GetOverallScore(ctx, currentStart, currentEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get current period score: %w", err)
	}

	previous, err := s.GetOverallScore(ctx, previousStart, previousEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to get previous period score: %w", err)
	}

	var change float64
	if previous.Score != 0 {
		change = ((current.Score - previous.Score) / previous.Score) * 100
	}

	return &domain.PeriodComparisonResult{
		PercentageChange: change,
		CurrentScore:     current.Score,
		PreviousScore:    previous.Score,
		CurrentCount:     current.RatingCount,
		PreviousCount:    previous.RatingCount,
	}, nil
}

// function to calculate time ranges for common comparisons
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
