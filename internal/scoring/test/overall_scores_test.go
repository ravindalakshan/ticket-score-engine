package scoring_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"ticket-score-engine/internal/scoring"
)

type mockOverallRepo struct {
	mock.Mock
}

func (m *mockOverallRepo) GetOverallScore(ctx context.Context, start, end time.Time) (float64, int, error) {
	args := m.Called(ctx, start, end)
	return args.Get(0).(float64), args.Int(1), args.Error(2)
}

func TestGetOverallScore_Success(t *testing.T) {
	mockRepo := new(mockOverallRepo)
	scorer := scoring.NewOverallScorer(mockRepo)

	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()

	expectedScore := 85.5
	expectedCount := 20

	mockRepo.On("GetOverallScore", mock.Anything, start, end).Return(expectedScore, expectedCount, nil)

	result, err := scorer.GetOverallScore(context.Background(), start, end)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedScore, result.Score)
	assert.Equal(t, expectedCount, result.RatingCount)

	mockRepo.AssertExpectations(t)
}

func TestGetOverallScore_Error(t *testing.T) {
	mockRepo := new(mockOverallRepo)
	scorer := scoring.NewOverallScorer(mockRepo)

	start := time.Now().AddDate(0, 0, -7)
	end := time.Now()

	mockRepo.On("GetOverallScore", mock.Anything, start, end).Return(0.0, 0, errors.New("db error"))

	result, err := scorer.GetOverallScore(context.Background(), start, end)

	assert.Error(t, err)
	assert.Nil(t, result)

	mockRepo.AssertExpectations(t)
}

func TestGetPeriodComparison_Success(t *testing.T) {
	mockRepo := new(mockOverallRepo)
	scorer := scoring.NewOverallScorer(mockRepo)

	now := time.Now()
	currentStart := now.AddDate(0, 0, -7)
	currentEnd := now
	previousStart := now.AddDate(0, 0, -14)
	previousEnd := now.AddDate(0, 0, -7)

	mockRepo.On("GetOverallScore", mock.Anything, currentStart, currentEnd).Return(90.0, 10, nil)
	mockRepo.On("GetOverallScore", mock.Anything, previousStart, previousEnd).Return(75.0, 8, nil)

	result, err := scorer.GetPeriodComparison(context.Background(), currentStart, currentEnd, previousStart, previousEnd)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.InDelta(t, 20.0, result.PercentageChange, 0.01)
	assert.Equal(t, 90.0, result.CurrentScore)
	assert.Equal(t, 75.0, result.PreviousScore)
	assert.Equal(t, 10, result.CurrentCount)
	assert.Equal(t, 8, result.PreviousCount)

	mockRepo.AssertExpectations(t)
}

func TestGetPeriodComparison_HandlesZeroPreviousScore(t *testing.T) {
	mockRepo := new(mockOverallRepo)
	scorer := scoring.NewOverallScorer(mockRepo)

	start1 := time.Now().AddDate(0, 0, -7)
	end1 := time.Now()
	start2 := time.Now().AddDate(0, 0, -14)
	end2 := time.Now().AddDate(0, 0, -7)

	mockRepo.On("GetOverallScore", mock.Anything, start1, end1).Return(80.0, 5, nil)
	mockRepo.On("GetOverallScore", mock.Anything, start2, end2).Return(0.0, 3, nil)

	result, err := scorer.GetPeriodComparison(context.Background(), start1, end1, start2, end2)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, float64(0), result.PercentageChange)
	assert.Equal(t, 80.0, result.CurrentScore)
	assert.Equal(t, 0.0, result.PreviousScore)
}

func TestGetComparisonPeriods_Invalid(t *testing.T) {
	_, _, _, _, err := scoring.GetComparisonPeriods("year")
	assert.Error(t, err)
	assert.Equal(t, "invalid period: year", err.Error())
}
