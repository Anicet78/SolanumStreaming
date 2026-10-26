package service

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"math"
	"strings"
)

type Segment struct {
	Index    int
	Start    float64
	Duration float64
}

type Playlist struct {
	ID       string
	Link     string
	Text     string
	Segments []Segment
	Probe    *ProbeResult
}

func buildSegments(keyframes []float64, total, target float64) []Segment {
	if len(keyframes) == 0 || total <= 0 {
		return nil
	}

	cuts := []float64{0}
	last := 0.0
	for _, k := range keyframes {
		if k-last >= target {
			cuts = append(cuts, k)
			last = k
		}
	}

	if len(cuts) > 1 && total-cuts[len(cuts)-1] < target/2 {
		cuts = cuts[:len(cuts)-1]
	}

	segs := make([]Segment, len(cuts))
	for i, start := range cuts {
		end := total
		if i+1 < len(cuts) {
			end = cuts[i+1]
		}
		segs[i] = Segment{Index: i, Start: start, Duration: end - start}
	}
	return segs
}

func buildPlaylist(segs []Segment) string {
	maxDur := 0.0
	for _, s := range segs {
		maxDur = max(maxDur, s.Duration)
	}

	var b strings.Builder
	b.WriteString("#EXTM3U\n#EXT-X-VERSION:7\n")
	fmt.Fprintf(&b, "#EXT-X-TARGETDURATION:%d\n", int(math.Ceil(maxDur)))
	b.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n")
	b.WriteString("#EXT-X-INDEPENDENT-SEGMENTS\n")
	b.WriteString("#EXT-X-MAP:URI=\"init.mp4\"\n")
	for _, s := range segs {
		fmt.Fprintf(&b, "#EXTINF:%.3f,\nseg_%05d.m4s\n", s.Duration, s.Index)
	}
	b.WriteString("#EXT-X-ENDLIST\n")
	return b.String()
}

func makeID(link string) string {
	h := sha1.Sum([]byte(link))
	return hex.EncodeToString(h[:])[:16]
}
