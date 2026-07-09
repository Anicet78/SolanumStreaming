package service

import (
	"context"
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
		mi, err := metainfo.LoadFromFile(torrentLink)
		if err != nil {
			return nil, err
		}

		t, err = s.client.AddTorrent(mi)
		if err != nil {
			return nil, err
		}
	}

	select {
	case <-t.GotInfo():
	case <-time.After(30 * time.Second):
		return nil, domain.ErrTorrentLoadingTimeout
	}

	var file *torrent.File
	for _, f := range t.Files() {
		if file == nil || f.Length() > file.Length() {
			file = f
		}
	}

	file.SetPriority(torrent.PiecePriorityNow)

	reader := file.NewReader()
	reader.SetReadahead(file.Length() / 100)

	return reader, nil
}
