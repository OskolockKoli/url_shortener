package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/OskolockKoli/url_shortener/internal/app"
	pb "github.com/OskolockKoli/url_shortener/proto"
)

func main() {
	port := flag.String("p", "50051", "Port to listen on")
	dbType := flag.String("d", "memory", "Database type: 'memory' or 'postgres'")
	flag.Parse()

	server, err := app.NewServer(*dbType)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}
	defer server.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterShortenerServiceServer(s, server)

	log.Printf("Starting gRPC server on port %s\n", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
