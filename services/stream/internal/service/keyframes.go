package service

import (
	"context"
	"errors"
	"io"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
)

type KeyframeExtractor interface {
	Keyframes(ctx context.Context, src io.ReadSeeker, duration float64) ([]float64, error)
}

func (s *StreamService) keyframes(ctx context.Context, link string, p *ProbeResult) ([]float64, error) {
	ex, ok := s.extractors[p.Container]
	if !ok {
		return nil, domain.ErrUnsupportedContainer
	}

	file, err := s.getFile(link)
	if err != nil {
		return nil, err
	}
	reader := file.NewReader()
	reader.SetResponsive()
	defer reader.Close()

	return ex.Keyframes(ctx, reader, p.Duration)
}

type mp4Extractor struct{}

func (mp4Extractor) Keyframes(ctx context.Context, src io.ReadSeeker, duration float64) ([]float64, error) {
	return nil, errors.New("mp4: not implemented")
}
