package repository_test

import (
	"context"
	"testing"
	"time"

	"ticket-score-engine/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetCategoryScores(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	start := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC) // Less than 30 days => daily

	rows := sqlmock.NewRows([]string{
		"category", "period", "count", "weighted_score", "total_weight",
	}).AddRow(
		"Support", "2024-01-02", 10, 45.0, 5.0,
	)

	mock.ExpectQuery("SELECT .* FROM ratings").
		WithArgs(start, end).
		WillReturnRows(rows)

	repo := repository.NewCategoryRepository(db)
	scores, err := repo.GetCategoryScores(context.Background(), start, end)

	assert.NoError(t, err)
	assert.Len(t, scores, 1)
	assert.Equal(t, "Support", scores[0].CategoryName)
	assert.Equal(t, "2024-01-02", scores[0].Date)
	assert.Equal(t, 10, scores[0].RatingCount)
	assert.InDelta(t, 900.0, scores[0].Score, 0.01) // (45/5)*100 = 900
}
