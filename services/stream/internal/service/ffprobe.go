package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
)

func (s *StreamService) rawURL(link string) string {
	return fmt.Sprintf("%s/raw?link=%s&internal_token=%s", s.InternalBase, url.QueryEscape(link), s.InternalToken)
}

type VideoInfo struct {
	Codec         string
	Width, Height int
	PixFmt        string // yuv420p10le = 10 bits
	ColorTransfer string // smpte2084 = HDR10, arib-std-b67 = HLG
}

type AudioInfo struct {
	Index    int
	Codec    string
	Channels int
	Language string
	Title    string
}

type Container string

const (
	ContainerMKV     Container = "mkv"
	ContainerMP4     Container = "mp4"
	ContainerUnknown Container = "unknown"
)

func containerFromFormatName(name string) Container {
	for _, n := range strings.Split(name, ",") {
		switch n {
		case "matroska", "webm":
			return ContainerMKV
		case "mov", "mp4":
			return ContainerMP4
		}
	}
	return ContainerUnknown
}

type ProbeResult struct {
	Duration  float64
	Container Container
	Video     *VideoInfo
	Audio     []AudioInfo
}

type ffprobeOut struct {
	Format struct {
		Duration   string `json:"duration"`
		FormatName string `json:"format_name"`
	} `json:"format"`
	Streams []struct {
		Index         int               `json:"index"`
		CodecType     string            `json:"codec_type"`
		CodecName     string            `json:"codec_name"`
		Width         int               `json:"width"`
		Height        int               `json:"height"`
		PixFmt        string            `json:"pix_fmt"`
		ColorTransfer string            `json:"color_transfer"`
		Channels      int               `json:"channels"`
		Tags          map[string]string `json:"tags"`
	} `json:"streams"`
}

func (s *StreamService) probe(ctx context.Context, link string) (*ProbeResult, error) {
	ctx, cancel := context.WithTimeout(ctx, 90*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "ffprobe",
		"-v", "error",
		"-print_format", "json",
		"-show_format",
		"-show_streams",
		s.rawURL(link),
	)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe: %w: %s", err, stderr.String())
	}

	var raw ffprobeOut
	if err := json.Unmarshal(out, &raw); err != nil {
		return nil, err
	}

	res := &ProbeResult{}
	res.Duration, _ = strconv.ParseFloat(raw.Format.Duration, 64)
	res.Container = containerFromFormatName(raw.Format.FormatName)

	for _, st := range raw.Streams {
		switch st.CodecType {
		case "video":
			if res.Video == nil { // on garde la première piste vidéo
				res.Video = &VideoInfo{st.CodecName, st.Width, st.Height, st.PixFmt, st.ColorTransfer}
			}
		case "audio":
			res.Audio = append(res.Audio, AudioInfo{
				Index: st.Index, Codec: st.CodecName, Channels: st.Channels,
				Language: st.Tags["language"], Title: st.Tags["title"],
			})
		}
	}
	if res.Video == nil {
		return nil, domain.ErrNoVideo
	}
	return res, nil
}
