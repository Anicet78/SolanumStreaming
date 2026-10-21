package service

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/anacrolix/torrent"
)

type loadedTorrent struct {
	once       sync.Once
	t          *torrent.Torrent
	file       *torrent.File
	err        error
	lastAccess time.Time
}

type StreamService struct {
	client        *torrent.Client
	mu            sync.Mutex
	torrents      map[string]*loadedTorrent
	extractors    map[Container]KeyframeExtractor
	InternalBase  string
	InternalToken string
}

func NewStreamService(client *torrent.Client, internalBase string, internalToken string) *StreamService {
	return &StreamService{
		client:   client,
		torrents: map[string]*loadedTorrent{},
		extractors: map[Container]KeyframeExtractor{
			ContainerMKV: mkvExtractor{},
			ContainerMP4: mp4Extractor{},
		},
		InternalToken: internalToken,
		InternalBase:  internalBase,
	}
}

func (s *StreamService) ServeRaw(w http.ResponseWriter, r *http.Request, link string) error {
	file, err := s.getFile(link)
	if err != nil {
		return err
	}

	reader := file.NewReader()
	reader.SetResponsive()
	reader.SetReadahead(20 * 1024 * 1024)
	defer reader.Close()

	http.ServeContent(w, r, file.DisplayPath(), time.Time{}, reader)
	return nil
}

func (s *StreamService) CreatePlaylist(ctx context.Context, link string) (*ProbeResult, error) {
	file, err := s.getFile(link)
	if err != nil {
		return nil, err
	}

	res, err := s.probe(ctx, link)
	if err != nil {
		return nil, err
	}

	kf, err := s.keyframes(ctx, file, res)
	if err != nil {
		return nil, err
	}
	slog.Info("keyframes", "count", len(kf), "first", kf[:min(10, len(kf))], "last", kf[len(kf)-1])

	return res, nil
}
