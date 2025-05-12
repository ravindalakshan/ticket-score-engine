package repository_test

import (
	"context"
	"testing"
	"time"

	"ticket-score-engine/internal/repository"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetScoresByTicket(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := repository.NewTicketRepository(db)

	start := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 5, 31, 23, 59, 59, 0, time.UTC)

	rows := sqlmock.NewRows([]string{
		"ticket_id", "category", "weighted_score", "total_weight",
	}).
		AddRow(1, "Grammer", 40.0, 50.0). // 80%
		AddRow(2, "GDPR", 25.0, 50.0)     // 50%

	mock.ExpectQuery("SELECT (.+) FROM ratings r").
		WithArgs(start, end).
		WillReturnRows(rows)

	result, err := repo.GetScoresByTicket(context.Background(), start, end)
	assert.NoError(t, err)
	assert.Len(t, result, 2)

	assert.Equal(t, 1, result[0].TicketID)
	assert.Equal(t, "Grammer", result[0].CategoryName)
	assert.InDelta(t, 80.0, result[0].Score, 0.01)

	assert.Equal(t, 2, result[1].TicketID)
	assert.Equal(t, "GDPR", result[1].CategoryName)
	assert.InDelta(t, 50.0, result[1].Score, 0.01)

	assert.NoError(t, mock.ExpectationsWereMet())
}
