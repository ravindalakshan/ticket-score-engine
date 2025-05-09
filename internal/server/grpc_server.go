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
