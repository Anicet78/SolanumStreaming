package service

import (
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
	client   *torrent.Client
	mu       sync.Mutex
	torrents map[string]*loadedTorrent
}

func NewStreamService(client *torrent.Client) *StreamService {
	return &StreamService{
		client:   client,
		torrents: map[string]*loadedTorrent{},
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
