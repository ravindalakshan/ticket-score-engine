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

type mockTicketRepo struct {
	mock.Mock
}

func (m *mockTicketRepo) GetScoresByTicket(ctx context.Context, start, end time.Time) ([]domain.TicketCategoryScore, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).([]domain.TicketCategoryScore), args.Error(1)
}

func TestGetTicketScores(t *testing.T) {
	mockRepo := new(mockTicketRepo)
	scorer := scoring.NewTicketScorer(mockRepo)

	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()

	expectedScores := []domain.TicketCategoryScore{
		{TicketID: 1, CategoryName: "Responsiveness", Score: 85.0},
		{TicketID: 1, CategoryName: "Knowledge", Score: 90.0},
	}

	mockRepo.On("GetScoresByTicket", mock.Anything, start, end).Return(expectedScores, nil)

	result, err := scorer.GetTicketScores(context.Background(), start, end)

	assert.NoError(t, err)
	assert.Equal(t, expectedScores, result)

	mockRepo.AssertExpectations(t)
}
