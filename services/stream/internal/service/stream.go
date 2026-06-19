package service

import (
	"context"
	"shared/tmdb"
)

type StreamService struct{}

func NewStreamService() *StreamService {
	return &StreamService{}
}

func (s *StreamService) Stream(ctx context.Context) (tmdb.SearchResponse, error) {
	var result tmdb.SearchResponse

	return result, nil
}
