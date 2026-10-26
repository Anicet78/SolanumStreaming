package service

import (
	"net/http"
	"strings"
	"time"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
)

func (s *StreamService) getFile(link string) (*torrent.File, error) {
	s.mu.Lock()
	entry, ok := s.torrents[link]
	if !ok {
		entry = &loadedTorrent{}
		s.torrents[link] = entry
	}
	entry.lastAccess = time.Now()
	s.mu.Unlock()

	entry.once.Do(func() {
		entry.t, entry.file, entry.err = s.load(link)
	})

	if entry.err != nil {
		s.mu.Lock()
		if s.torrents[link] == entry {
			delete(s.torrents, link)
		}
		s.mu.Unlock()
		return nil, entry.err
	}

	return entry.file, nil
}

func (s *StreamService) load(link string) (*torrent.Torrent, *torrent.File, error) {
	t, err := loadTorrent(s.client, link)
	if err != nil {
		return nil, nil, err
	}

	select {
	case <-t.GotInfo():
	case <-time.After(30 * time.Second):
		t.Drop()
		return nil, nil, domain.ErrTorrentLoadingTimeout
	}

	var biggest *torrent.File
	for _, f := range t.Files() {
		if biggest == nil || f.Length() > biggest.Length() {
			biggest = f
		}
	}
	if biggest == nil {
		t.Drop()
		return nil, nil, domain.ErrTorrentNoFile
	}

	return t, biggest, nil
}

func loadTorrent(client *torrent.Client, torrentLink string) (*torrent.Torrent, error) {
	var t *torrent.Torrent
	var err error

	if strings.HasPrefix(torrentLink, "magnet:") {
		t, err = client.AddMagnet(torrentLink)
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
				t, err = client.AddMagnet(location)
				if err != nil {
					return nil, err
				}
			}
		} else {
			mi, err := metainfo.Load(resp.Body)
			if err != nil {
				return nil, err
			}
			t, err = client.AddTorrent(mi)
			if err != nil {
				return nil, err
			}
		}
	}
	return t, nil
}
