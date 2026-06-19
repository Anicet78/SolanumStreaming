package service

import (
	"context"
)

type StreamService struct{}

func NewStreamService() *StreamService {
	return &StreamService{}
}

func (s *StreamService) Stream(ctx context.Context) error {

	return nil
}
