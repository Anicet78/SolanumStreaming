package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
)

func splitMP4(data []byte) (init, media []byte, err error) {
	off := 0
	for off+8 <= len(data) {
		size := int(binary.BigEndian.Uint32(data[off : off+4]))
		typ := string(data[off+4 : off+8])
		if size == 1 && off+16 <= len(data) {
			size = int(binary.BigEndian.Uint64(data[off+8 : off+16]))
		}
		if size < 8 || off+size > len(data) {
			return nil, nil, domain.ErrMp4InvalidBox
		}
		if typ == "ftyp" || typ == "moov" {
			init = append(init, data[off:off+size]...)
		} else {
			media = append(media, data[off:off+size]...)
		}
		off += size
	}
	if len(init) == 0 {
		return nil, nil, domain.ErrMp4NoMoov
	}
	return init, media, nil
}

func audioArgs(a AudioInfo) []string {
	switch a.Codec {
	case "aac", "opus":
		return []string{"-c:a", "copy"}
	default:
		return []string{"-c:a", "aac", "-b:a", "256k", "-ac", "2"}
	}
}

func (s *StreamService) generate(ctx context.Context, pl *Playlist, seg Segment) (init, media []byte, err error) {
	if len(pl.Probe.Audio) == 0 {
		return nil, nil, domain.ErrNoAudioTrack
	}
	audio := pl.Probe.Audio[0]

	args := []string{
		"-v", "error",
		"-ss", fmt.Sprintf("%.6f", seg.Start+0.002),
		"-i", s.rawURL(pl.Link),
		"-t",  fmt.Sprintf("%.6f", seg.Duration),
		"-map", "0:v:0",
		"-map", fmt.Sprintf("0:%d", audio.Index),
		"-c:v", "copy",
	}
	if pl.Probe.Video.Codec == "hevc" {
		args = append(args, "-tag:v", "hvc1")
	}
	args = append(args, audioArgs(audio)...)
	args = append(args,
		"-avoid_negative_ts", "make_zero",
		"-output_ts_offset", fmt.Sprintf("%.3f", seg.Start),
		"-movflags", "empty_moov+frag_keyframe+default_base_moof",
		"-f", "mp4",
		"pipe:1",
	)

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return nil, nil, fmt.Errorf("ffmpeg: %w: %s", err, stderr.String())
	}
	_ = os.Stderr
	return splitMP4(out.Bytes())
}

func (s *StreamService) GetInit(ctx context.Context, pl *Playlist) ([]byte, error) {
	init, _, err := s.generate(ctx, pl, pl.Segments[0])
	return init, err
}

func (s *StreamService) GetSegment(ctx context.Context, pl *Playlist, index int) ([]byte, error) {
	if index < 0 || index >= len(pl.Segments) {
		return nil, domain.ErrSegmentOutOfRange
	}
	_, media, err := s.generate(ctx, pl, pl.Segments[index])
	return media, err
}
