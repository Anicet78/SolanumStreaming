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

	defer reader.Close()
	defer t.Drop()

	responseWriter.Header().Set("Content-Type", "video/mp4")
	responseWriter.Header().Set("Cache-Control", "no-cache")
	responseWriter.Header().Del("Accept-Ranges")
	responseWriter.WriteHeader(http.StatusOK)
	if f, ok := responseWriter.(http.Flusher); ok {
		f.Flush()
	}

	// log.Println("waiting for initial pieces...")
	// for {
	// 	stats := t.Stats()
	// 	downloaded := stats.BytesReadUsefulData.Int64()
	// 	log.Printf("downloaded: %d bytes", downloaded)
	// 	if downloaded > 5*1024*1024 { // attend 5MB
	// 		break
	// 	}
	// 	time.Sleep(500 * time.Millisecond)
	// }
	// log.Println("enough data, starting ffmpeg")

	cmd := exec.CommandContext(ctx, "ffmpeg",
		"-fflags", "nobuffer",
		"-flags", "low_delay",
		"-probesize", "10M",
		"-analyzeduration", "0",
		"-i", "pipe:0",
		"-c:v", "libx264",
		"-preset", "ultrafast",
		"-tune", "zerolatency",
		"-vf", "format=yuv420p",
		"-profile:v", "high",
		"-level", "4.1",
		"-c:a", "aac",
		"-b:a", "128k",
		"-ac", "2",
		"-movflags", "frag_keyframe+empty_moov",
		"-f", "mp4",
		"pipe:1",
	)

	cmd.Stdin = reader
	cmd.Stdout = responseWriter
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return nil
		}
		log.Printf("ffmpeg error: %v", err)
		return err
	}

	return nil
}
