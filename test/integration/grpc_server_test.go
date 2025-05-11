package integration

import (
	"context"
	"database/sql"
	"log"
	"net"
	"testing"
	"time"

	pb "ticket-score-engine/generated"
	"ticket-score-engine/internal/server"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func startTestGRPCServer(t *testing.T, db *sql.DB) (pb.ScoringServiceClient, func()) {
	// Create listener
	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	srv := server.NewTicketScoreServer(db)
	pb.RegisterScoringServiceServer(grpcServer, srv)

	// Run server in background
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Dial the server
	conn, err := grpc.Dial(lis.Addr().String(), grpc.WithInsecure())
	require.NoError(t, err)

	client := pb.NewScoringServiceClient(conn)

	// Return client
	return client, func() {
		grpcServer.Stop()
		conn.Close()
		lis.Close()
	}
}

func TestGetOverallScore(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	// Set up expected SQL query and mock result
	start := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 5, 31, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT (.+) FROM ratings r").
		WithArgs(start, end).
		WillReturnRows(sqlmock.NewRows([]string{"total_weighted_score", "total_weight", "rating_count"}).
			AddRow(75.0, 100.0, 15))

	client, cleanup := startTestGRPCServer(t, db)
	defer cleanup()

	req := &pb.ScoreRequest{
		StartDate: "2024-05-01",
		EndDate:   "2024-05-31",
	}

	resp, err := client.GetOverallScore(context.Background(), req)
	require.NoError(t, err)
	require.Equal(t, float32(75.0), resp.Score)
	require.Equal(t, int32(15), resp.RatingCount)

	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetCategoryScores(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	start := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 5, 7, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT (.+) FROM ratings r").
		WithArgs(start, end).
		WillReturnRows(sqlmock.NewRows([]string{"category", "period", "count", "weighted_score", "total_weight"}).
			AddRow("GDPR", "2024-05-01", 10, 40.0, 50.0).
			AddRow("Spelling", "2024-05-01", 5, 20.0, 25.0))

	client, cleanup := startTestGRPCServer(t, db)
	defer cleanup()

	req := &pb.ScoreRequest{
		StartDate: "2024-05-01",
		EndDate:   "2024-05-07",
	}

	resp, err := client.GetCategoryScores(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, resp.Scores, 2)
	require.Equal(t, "GDPR", resp.Scores[0].CategoryName)
	require.Equal(t, float32(80.0), resp.Scores[0].Score) // 40/50 * 100
	require.Equal(t, int32(10), resp.Scores[0].RatingCount)
}

func TestGetTicketScores(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	start := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2024, 5, 2, 0, 0, 0, 0, time.UTC)

	// Mocking 4 columns: ticket_id, category_name, weighted_score, total_weight
	mock.ExpectQuery("SELECT (.+) FROM ratings r").
		WithArgs(start, end).
		WillReturnRows(sqlmock.NewRows([]string{
			"ticket_id", "category", "weighted_score", "total_weight",
		}).
			AddRow(101, "Spelling", 40.0, 50.0). // Score = (40 / 50) * 100 = 80
			AddRow(101, "Grammer", 30.0, 60.0).  // Score = (30 / 60) * 100 = 50
			AddRow(102, "GDPR", 90.0, 90.0))     // Score = 100

	client, cleanup := startTestGRPCServer(t, db)
	defer cleanup()

	req := &pb.ScoreRequest{
		StartDate: "2024-05-01",
		EndDate:   "2024-05-02",
	}

	resp, err := client.GetTicketScores(context.Background(), req)
	require.NoError(t, err)
	require.Len(t, resp.TicketScores, 2)

	ticketMap := make(map[int32]map[string]float32)
	for _, t := range resp.TicketScores {
		ticketMap[t.TicketId] = t.CategoryScores
	}

	require.Equal(t, float32(80.0), ticketMap[101]["Spelling"])
	require.Equal(t, float32(50.0), ticketMap[101]["Grammer"])
	require.Equal(t, float32(100.0), ticketMap[102]["GDPR"])
}

func TestGetPeriodComparison(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	defer db.Close()

	client, cleanup := startTestGRPCServer(t, db)
	defer cleanup()

	currentStart := time.Date(2024, 5, 1, 0, 0, 0, 0, time.UTC)
	currentEnd := time.Date(2024, 5, 2, 0, 0, 0, 0, time.UTC)
	previousStart := time.Date(2024, 4, 29, 0, 0, 0, 0, time.UTC)
	previousEnd := time.Date(2024, 4, 30, 0, 0, 0, 0, time.UTC)

	mock.ExpectQuery("SELECT (.+) FROM ratings r").
		WithArgs(currentStart, currentEnd).
		WillReturnRows(sqlmock.NewRows([]string{
			"total_weighted_score", "total_weight", "rating_count",
		}).AddRow(60.0, 100.0, 10))

	mock.ExpectQuery("SELECT (.+) FROM ratings r").
		WithArgs(previousStart, previousEnd).
		WillReturnRows(sqlmock.NewRows([]string{
			"total_weighted_score", "total_weight", "rating_count",
		}).AddRow(40.0, 100.0, 8))

	req := &pb.PeriodComparisonRequest{
		CurrentPeriod: &pb.ScoreRequest{
			StartDate: "2024-05-01",
			EndDate:   "2024-05-02",
		},
		PreviousPeriod: &pb.ScoreRequest{
			StartDate: "2024-04-29",
			EndDate:   "2024-04-30",
		},
	}

	resp, err := client.GetPeriodComparison(context.Background(), req)
	require.NoError(t, err)

	require.InDelta(t, 50.0, resp.PercentageChange, 0.01)
	require.InDelta(t, 60.0, resp.CurrentScore, 0.01)
	require.InDelta(t, 40.0, resp.PreviousScore, 0.01)
	require.Equal(t, int32(10), resp.CurrentCount)
	require.Equal(t, int32(8), resp.PreviousCount)
}
