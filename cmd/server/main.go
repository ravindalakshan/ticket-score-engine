package main

import (
	"database/sql"
	"log"
	"net"

	"ticket-score-engine/internal/server"

	pb "ticket-score-engine/generated" // generated proto package

	"google.golang.org/grpc"
	_ "modernc.org/sqlite"
)

func main() {
	log.Println("Starting Ticket Score Engine...")

	db, err := sql.Open("sqlite", "./database.db")
	if err != nil {
		log.Fatalf("Failed to open DB: %v", err)
	}
	defer db.Close()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterScoringServiceServer(grpcServer, server.NewTicketScoreServer(db))

	log.Println("gRPC server listening on :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
