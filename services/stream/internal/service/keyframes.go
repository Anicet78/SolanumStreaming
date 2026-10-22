package service

import (
	"context"
	"io"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
)

type KeyframeExtractor interface {
	Keyframes(ctx context.Context, src io.ReadSeeker, duration float64) ([]float64, error)
}

func (s *StreamService) keyframes(ctx context.Context, torrentLink string, p *ProbeResult) ([]float64, error) {
	ex, ok := s.extractors[p.Container]
	if !ok {
		return nil, domain.ErrUnsupportedContainer
	}

	file, err := s.getFile(torrentLink)
	if err != nil {
		return nil, err
	}

	reader := file.NewReader()
	reader.SetResponsive()
	defer reader.Close()

	return ex.Keyframes(ctx, reader, p.Duration)
}
