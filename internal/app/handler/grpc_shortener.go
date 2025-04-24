package handler

import (
	"context"
	"errors"

	"github.com/bufbuild/protovalidate-go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/patraden/ya-practicum-go-shortly/api/shortener/v1"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/config"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/domain"
	e "github.com/patraden/ya-practicum-go-shortly/internal/app/domain/errors"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/middleware"
	"github.com/patraden/ya-practicum-go-shortly/internal/app/service/shortener"
)

// GRPCShortenerHandler provides gRPC request handling for URL shortening operations.
type GRPCShortenerHandler struct {
	service   shortener.URLShortener
	config    *config.Config
	log       *zerolog.Logger
	validator protovalidate.Validator
	pb.UnimplementedURLShortenerServiceServer
}

// NewGRPCURLShortenerHandler creates a new instance of GRPCShortenerHandler using service interface.
func NewGRPCURLShortenerHandler(
	service shortener.URLShortener,
	config *config.Config,
	log *zerolog.Logger,
) (*GRPCShortenerHandler, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &GRPCShortenerHandler{
		service:   service,
		config:    config,
		log:       log,
		validator: validator,
	}, nil
}

// NewGRPCShortenerHandler creates a new instance of GRPCShortenerHandler.
func NewGRPCShortenerHandler(
	config *config.Config,
	service *shortener.InsistentShortener,
	log *zerolog.Logger,
) (*GRPCShortenerHandler, error) {
	validator, err := protovalidate.New()
	if err != nil {
		return nil, err
	}

	return &GRPCShortenerHandler{
		service:   service,
		config:    config,
		log:       log,
		validator: validator,
	}, nil
}

// ShortenURL handles requests to shorten a given URL.
func (h *GRPCShortenerHandler) ShortenURL(
	ctx context.Context,
	r *pb.ShortenURLRequest,
) (*pb.ShortenURLResponse, error) {
	if err := h.validator.Validate(r); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad Request")
	}

	slug, err := h.service.ShortenURL(ctx, domain.OriginalURL(r.GetUrl()))
	if err != nil && !errors.Is(err, e.ErrOriginalExists) {
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	if errors.Is(err, e.ErrOriginalExists) {
		return nil, status.Error(codes.AlreadyExists, "Conflict")
	}

	return &pb.ShortenURLResponse{Slug: slug.WithBaseURL(h.config.BaseURL)}, nil
}

// GetOriginalURL handles requests to retrieve the original URL from a shortened slug.
func (h *GRPCShortenerHandler) GetOriginalURL(
	ctx context.Context,
	r *pb.GetOriginalURLRequest,
) (*pb.GetOriginalURLResponse, error) {
	if err := h.validator.Validate(r); err != nil {
		return nil, status.Error(codes.InvalidArgument, "Bad Request")
	}

	slug := domain.Slug(r.GetSlug())
	original, err := h.service.GetOriginalURL(ctx, slug)

	switch {
	case errors.Is(err, e.ErrSlugInvalid):
		return nil, status.Error(codes.InvalidArgument, "Bad Request")

	case errors.Is(err, e.ErrSlugNotFound):
		return nil, status.Error(codes.NotFound, "Not Found")

	case errors.Is(err, e.ErrSlugDeleted):
		return nil, status.Error(codes.NotFound, "Deleted")

	case errors.Is(err, e.ErrShortenerInternal) || err != nil:
		return nil, status.Error(codes.Internal, "Internal Server Error")
	}

	return &pb.GetOriginalURLResponse{Url: original.String()}, nil
}

// Interceptors returns interceptors that should be used with the handler.
func (h *GRPCShortenerHandler) Interceptors() []grpc.UnaryServerInterceptor {
	filter := func(method string) bool {
		switch method {
		case pb.URLShortenerService_ShortenURL_FullMethodName:
			return true
		case pb.URLShortenerService_GetOriginalURL_FullMethodName:
			return true
		default:
			return false
		}
	}

	return []grpc.UnaryServerInterceptor{
		middleware.AuthenticateGRPC(filter, h.log, h.config),
	}
}
