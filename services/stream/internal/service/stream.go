package service

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

type StreamService struct {
	client *torrent.Client
}

func NewStreamService(client *torrent.Client) *StreamService {
	return &StreamService{client: client}
}

func (s *StreamService) Stream(ctx context.Context, torrentLink string) (torrent.Reader, error) {
	var t *torrent.Torrent
	var err error

	if strings.HasPrefix(torrentLink, "magnet:") {
		t, err = s.client.AddMagnet(torrentLink)
		if err != nil {
			return nil, err
		}
	} else {
		httpClient := &http.Client{
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if strings.HasPrefix(req.URL.String(), "magnet:") {
					return http.ErrUseLastResponse
				}
				return nil
			},
		}

		resp, err := httpClient.Get(torrentLink)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
			location := resp.Header.Get("Location")
			if strings.HasPrefix(location, "magnet:") {
				t, err = s.client.AddMagnet(location)
				if err != nil {
					return nil, err
				}
			}
		} else {
			mi, err := metainfo.Load(resp.Body)
			if err != nil {
				return nil, err
			}
			t, err = s.client.AddTorrent(mi)
			if err != nil {
				return nil, err
			}
		}
	}

	select {
	case <-t.GotInfo():
	case <-time.After(30 * time.Second):
		return nil, domain.ErrTorrentLoadingTimeout
	}

	go func() {
		for {
			time.Sleep(2 * time.Second)
			stats := t.Stats()
			log.Printf("peers: %d, downloaded: %d bytes",
				stats.ActivePeers,
				stats.BytesReadUsefulData.Int64())
		}
	}()

	var file *torrent.File
	for _, f := range t.Files() {
		if file == nil || f.Length() > file.Length() {
			file = f
		}
	}

	log.Println("got info, files:", len(t.Files()))
	log.Println("selected file:", file.DisplayPath(), "size:", file.Length())

	file.SetPriority(torrent.PiecePriorityNow)

	reader := file.NewReader()
	reader.SetReadahead(file.Length() / 100)

	return reader, nil
}
