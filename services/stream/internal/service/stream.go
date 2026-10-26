package service

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
	"github.com/anacrolix/torrent"
	"golang.org/x/sync/singleflight"
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
	playlists     map[string]*Playlist
	group         singleflight.Group
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
		playlists:     map[string]*Playlist{},
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

func (s *StreamService) CreatePlaylist(ctx context.Context, link string) (string, error) {
	id := makeID(link)

	s.mu.Lock()
	_, ok := s.playlists[id]
	s.mu.Unlock()
	if ok {
		return id, nil
	}

	_, err, _ := s.group.Do(id, func() (any, error) {
		bctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 2*time.Minute)
		defer cancel()

		probe, err := s.probe(bctx, link)
		if err != nil {
			return nil, err
		}

		kf, err := s.keyframes(bctx, link, probe)
		if err != nil {
			return nil, err
		}

		segs := buildSegments(kf, probe.Duration, 8)
		if len(segs) == 0 {
			return nil, domain.ErrNoSegments
		}

		s.mu.Lock()
		s.playlists[id] = &Playlist{
			ID: id, Link: link, Text: buildPlaylist(segs),
			Segments: segs, Probe: probe,
		}
		s.mu.Unlock()
		return nil, nil
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *StreamService) GetPlaylist(id string) (*Playlist, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.playlists[id]
	return p, ok
}
