package server

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pb "ticket-score-engine/generated"
	"ticket-score-engine/internal/scoring"
)

type ticketScoreServer struct {
	pb.UnimplementedScoringServiceServer
	db *sql.DB
}

func NewTicketScoreServer(db *sql.DB) pb.ScoringServiceServer {
	return &ticketScoreServer{db: db}
}

func (s *ticketScoreServer) GetCategoryScores(ctx context.Context, req *pb.ScoreRequest) (*pb.ScoreResponse, error) {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}

	scores, err := scoring.GetCategoryScores(s.db, start, end)
	if err != nil {
		return nil, err
	}

	var resp pb.ScoreResponse
	for _, s := range scores {
		resp.Scores = append(resp.Scores, &pb.CategoryScore{
			CategoryName: s.CategoryName,
			Date:         s.Date,
			Score:        float32(s.Score),
			RatingCount:  int32(s.RatingCount),
		})
	}

	return &resp, nil
}

func (s *ticketScoreServer) GetTicketScores(ctx context.Context, req *pb.ScoreRequest) (*pb.TicketScoreResponse, error) {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, err
	}

	ticketCategoryScores, err := scoring.GetScoresByTicket(s.db, start, end)
	if err != nil {
		return nil, err
	}

	// Group by TicketID
	ticketMap := make(map[int64]map[string]float32)
	for _, score := range ticketCategoryScores {
		ticketID := int64(score.TicketID)
		if _, exists := ticketMap[ticketID]; !exists {
			ticketMap[ticketID] = make(map[string]float32)
		}
		ticketMap[ticketID][score.CategoryName] = float32(score.Score)
	}

	// Build gRPC response
	var grpcTicketScores []*pb.TicketScore
	for ticketID, categoryScores := range ticketMap {
		grpcTicketScores = append(grpcTicketScores, &pb.TicketScore{
			TicketId:       int32(ticketID),
			CategoryScores: categoryScores,
		})
	}

	return &pb.TicketScoreResponse{
		TicketScores: grpcTicketScores,
	}, nil
}

//

func (s *ticketScoreServer) GetOverallScore(ctx context.Context, req *pb.ScoreRequest) (*pb.OverallScoreResponse, error) {
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid start date: %w", err)
	}
	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid end date: %w", err)
	}

	result, err := scoring.GetOverallScore(s.db, start, end)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate overall score: %w", err)
	}

	return &pb.OverallScoreResponse{
		Score:       float32(result.Score),
		RatingCount: int32(result.RatingCount),
	}, nil
}

func (s *ticketScoreServer) GetPeriodComparison(ctx context.Context, req *pb.PeriodComparisonRequest) (*pb.PeriodComparisonResponse, error) {
	// Parse current period dates
	currentStart, err := time.Parse("2006-01-02", req.CurrentPeriod.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid current period start date: %w", err)
	}
	currentEnd, err := time.Parse("2006-01-02", req.CurrentPeriod.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid current period end date: %w", err)
	}

	// Parse previous period dates
	previousStart, err := time.Parse("2006-01-02", req.PreviousPeriod.StartDate)
	if err != nil {
		return nil, fmt.Errorf("invalid previous period start date: %w", err)
	}
	previousEnd, err := time.Parse("2006-01-02", req.PreviousPeriod.EndDate)
	if err != nil {
		return nil, fmt.Errorf("invalid previous period end date: %w", err)
	}

	// Get comparison results
	result, err := scoring.GetPeriodComparison(s.db, currentStart, currentEnd, previousStart, previousEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to compare periods: %w", err)
	}

	return &pb.PeriodComparisonResponse{
		PercentageChange: float32(result.PercentageChange),
		CurrentScore:     float32(result.CurrentScore),
		PreviousScore:    float32(result.PreviousScore),
		CurrentCount:     int32(result.CurrentCount),
		PreviousCount:    int32(result.PreviousCount),
	}, nil
}
