package app

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"

	"github.com/OskolockKoli/url_shortener/internal/db"
	"github.com/OskolockKoli/url_shortener/internal/models"
	"github.com/OskolockKoli/url_shortener/pkg/shortener"
	pb "github.com/OskolockKoli/url_shortener/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Server struct {
	db     db.Database
	shorty *shortener.Shortener
}

// NewServer creates a new instance of the Server with specified database type.
func NewServer(dbType string) (*Server, error) {
	var dbInstance db.Database
	switch dbType {
	case "memory":
		dbInstance = &db.MemoryDB{}
	case "postgres":
		connStr := os.Getenv("POSTGRES_URL")
		if connStr == "" {
			return nil, fmt.Errorf("POSTGRES_URL environment variable is not set")
		}
		pgDb, err := sql.Open("postgres", connStr)
		if err != nil {
			return nil, fmt.Errorf("failed to open postgres connection: %w", err)
		}
		dbInstance = &db.PostgreSQL{
			DB: pgDb,
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %s", dbType)
	}

	s := &Server{
		db:     dbInstance,
		shorty: shortener.New(),
	}

	return s, nil
}

// Close closes the underlying database connection.
func (s *Server) Close() {
	if closer, ok := s.db.(io.Closer); ok {
		closer.Close()
	}
}

// CreateShortLink implements the gRPC method for creating a short link.
func (s *Server) CreateShortLink(ctx context.Context, req *pb.CreateShortLinkRequest) (*pb.CreateShortLinkResponse, error) {
	url := req.Url
	if url == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty URL provided")
	}

	shortID := s.shorty.GenerateShortID()
	link := models.Link{
		ShortID: shortID,
		URL:     url,
	}

	err := s.db.Save(link)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to save link: %v", err)
	}

	return &pb.CreateShortLinkResponse{
		ShortLink: shortID,
	}, nil
}

// ResolveShortLink implements the gRPC method for resolving a short link.
func (s *Server) ResolveShortLink(ctx context.Context, req *pb.ResolveShortLinkRequest) (*pb.ResolveShortLinkResponse, error) {
	shortID := req.ShortLink
	if shortID == "" {
		return nil, status.Errorf(codes.InvalidArgument, "Empty short link provided")
	}

	link, err := s.db.GetByShortID(shortID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "Failed to resolve short link: %v", err)
	}

	return &pb.ResolveShortLinkResponse{
		Url: link.URL,
	}, nil
}
