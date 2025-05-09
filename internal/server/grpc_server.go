package server

import (
	"context"
	"database/sql"
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
