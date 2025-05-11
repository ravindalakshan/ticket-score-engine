package repository_test

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"ticket-score-engine/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetOverallScore_Success(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	start := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	repo := repository.NewOverallRepository(db)

	totalWeightedScore := 75.0
	totalWeight := 100.0
	ratingCount := 15

	mock.ExpectQuery("SELECT SUM\\(\\(r.rating \\* 1.0 / 5.0\\) \\* rc.weight\\)").
		WithArgs(start, end).
		WillReturnRows(sqlmock.NewRows([]string{"total_weighted_score", "total_weight", "rating_count"}).
			AddRow(totalWeightedScore, totalWeight, ratingCount))

	score, count, err := repo.GetOverallScore(context.Background(), start, end)

	assert.NoError(t, err)
	assert.Equal(t, ratingCount, count)
	assert.InDelta(t, 75.0, score, 0.01)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOverallScore_ZeroWeight(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	start := time.Now().AddDate(0, -1, 0)
	end := time.Now()

	repo := repository.NewOverallRepository(db)

	mock.ExpectQuery("SELECT SUM\\(\\(r.rating \\* 1.0 / 5.0\\) \\* rc.weight\\)").
		WithArgs(start, end).
		WillReturnRows(sqlmock.NewRows([]string{"total_weighted_score", "total_weight", "rating_count"}).
			AddRow(0.0, 0.0, 10))

	score, count, err := repo.GetOverallScore(context.Background(), start, end)

	assert.NoError(t, err)
	assert.Equal(t, 10, count)
	assert.Equal(t, float64(0), score)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOverallScore_QueryError(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	start := time.Now().AddDate(0, -1, 0)
	end := time.Now()

	repo := repository.NewOverallRepository(db)

	mock.ExpectQuery("SELECT SUM\\(\\(r.rating \\* 1.0 / 5.0\\) \\* rc.weight\\)").
		WithArgs(start, end).
		WillReturnError(sql.ErrConnDone)

	score, count, err := repo.GetOverallScore(context.Background(), start, end)

	assert.Error(t, err)
	assert.Equal(t, float64(0), score)
	assert.Equal(t, 0, count)

	assert.NoError(t, mock.ExpectationsWereMet())
}
