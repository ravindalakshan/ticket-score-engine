package scoring_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ticket-score-engine/internal/domain"
	"ticket-score-engine/internal/scoring"
)

type mockCategoryRepo struct {
	mock.Mock
}

func (m *mockCategoryRepo) GetCategoryScores(ctx context.Context, start, end time.Time) ([]domain.CategoryScore, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]domain.CategoryScore), args.Error(1)
}

func TestGetCategoryScores(t *testing.T) {
	mockRepo := new(mockCategoryRepo)
	scorer := scoring.NewCategoryScorer(mockRepo)

	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()

	expected := []domain.CategoryScore{
		{CategoryName: "Efficiency", Date: "2025-05-01", RatingCount: 10, Score: 87.5},
		{CategoryName: "Communication", Date: "2025-05-01", RatingCount: 12, Score: 90.0},
	}

	mockRepo.On("GetCategoryScores", mock.Anything, start, end).Return(expected, nil)

	result, err := scorer.GetCategoryScores(context.Background(), start, end)

	assert.NoError(t, err)
	assert.Equal(t, expected, result)

	mockRepo.AssertExpectations(t)
}
