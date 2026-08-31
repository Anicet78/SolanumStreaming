package service

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/exec"
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

func (s *StreamService) Stream(ctx context.Context, responseWriter http.ResponseWriter, torrentLink string) error {
	var t *torrent.Torrent
	var err error

	if strings.HasPrefix(torrentLink, "magnet:") {
		t, err = s.client.AddMagnet(torrentLink)
		if err != nil {
			return err
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
			return err
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
			location := resp.Header.Get("Location")
			if strings.HasPrefix(location, "magnet:") {
				t, err = s.client.AddMagnet(location)
				if err != nil {
					return err
				}
			}
		} else {
			mi, err := metainfo.Load(resp.Body)
			if err != nil {
				return err
			}
			t, err = s.client.AddTorrent(mi)
			if err != nil {
				return err
			}
		}
	}

	select {
	case <-t.GotInfo():
	case <-time.After(30 * time.Second):
		return domain.ErrTorrentLoadingTimeout
	}

	var file *torrent.File
	for _, f := range t.Files() {
		if file == nil || f.Length() > file.Length() {
			file = f
		}
	}

	file.SetPriority(torrent.PiecePriorityNow)

	reader := file.NewReader()
	reader.SetReadahead(20 * 1024 * 1024)

	streamCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer reader.Close()
	defer t.Drop()

	responseWriter.Header().Set("Content-Type", "video/mp4")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.WriteHeader(http.StatusOK)
	if f, ok := responseWriter.(http.Flusher); ok {
		f.Flush()
	}

	cmd := exec.CommandContext(streamCtx, "ffmpeg",
		"-fflags", "nobuffer",
		"-probesize", "10M",
		"-analyzeduration", "0",
		"-i", "pipe:0",
		"-map", "0:v:0",
		"-map", "0:a:0",
		"-c:v", "copy",
		"-c:a", "copy",
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4",
		"pipe:1",
	)

	cmd.Stdin = reader
	cmd.Stdout = responseWriter
	cmd.Stderr = os.Stderr

	go func() {
		select {
		case <-ctx.Done():
			cancel()
		case <-streamCtx.Done():
		}
	}()

	if err := cmd.Run(); err != nil {
		if streamCtx.Err() != nil || ctx.Err() != nil {
			return nil
		}
		log.Printf("ffmpeg error: %v", err)
		return err
	}

	return nil
}
