package service

import (
	"context"
	"encoding/binary"
	"io"
	"sort"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
)

type mp4Extractor struct{}

func walkBoxes(data []byte, fn func(typ string, body []byte)) error {
	off := 0
	for off+8 <= len(data) {
		size := int64(binary.BigEndian.Uint32(data[off : off+4]))
		typ := string(data[off+4 : off+8])
		hdr := 8
		switch size {
		case 1:
			if off+16 > len(data) {
				return domain.ErrInvalidData
			}
			size = int64(binary.BigEndian.Uint64(data[off+8 : off+16]))
			hdr = 16
		case 0:
			size = int64(len(data) - off)
		}
		if size < int64(hdr) || int64(off)+size > int64(len(data)) {
			return domain.ErrInvalidData
		}
		fn(typ, data[off+hdr:off+int(size)])
		off += int(size)
	}
	return nil
}

type run struct {
	count uint32
	val   int64
}

func readRuns(b []byte) []run {
	if len(b) < 8 {
		return nil
	}
	n := int(binary.BigEndian.Uint32(b[4:8]))
	b = b[8:]
	if n > len(b)/8 {
		n = len(b) / 8
	}
	runs := make([]run, n)
	for i := 0; i < n; i++ {
		runs[i] = run{
			count: binary.BigEndian.Uint32(b[i*8:]),
			val:   int64(int32(binary.BigEndian.Uint32(b[i*8+4:]))),
		}
	}
	return runs
}

func valueAt(runs []run, samples []uint32, cumulative bool) []int64 {
	out := make([]int64, 0, len(samples))
	ri := 0
	var start uint32
	var acc int64

	for _, n := range samples {
		if n == 0 {
			continue
		}
		s := n - 1
		for ri < len(runs) && s >= start+runs[ri].count {
			acc += int64(runs[ri].count) * runs[ri].val
			start += runs[ri].count
			ri++
		}
		if ri >= len(runs) {
			if cumulative {
				return out
			}
			out = append(out, 0)
			continue
		}
		if cumulative {
			out = append(out, acc+int64(s-start)*runs[ri].val)
		} else {
			out = append(out, runs[ri].val)
		}
	}
	return out
}

func (mp4Extractor) Keyframes(ctx context.Context, src io.ReadSeeker, duration float64) ([]float64, error) {
	pos := int64(0)
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if _, err := src.Seek(pos, io.SeekStart); err != nil {
			return nil, err
		}

		var h [16]byte
		if _, err := io.ReadFull(src, h[:8]); err != nil {
			return nil, domain.ErrMoovNotFoundData
		}
		size := int64(binary.BigEndian.Uint32(h[:4]))
		typ := string(h[4:8])
		hdr := int64(8)
		if size == 1 {
			if _, err := io.ReadFull(src, h[8:16]); err != nil {
				return nil, err
			}
			size = int64(binary.BigEndian.Uint64(h[8:16]))
			hdr = 16
		}
		if size == 0 {
			return nil, domain.ErrBoxWithUnknownSize
		}
		if size < hdr {
			return nil, domain.ErrInvalidData
		}

		switch typ {
		case "moof":
			return nil, domain.ErrFmp4NotSupported
		case "moov":
			body, err := readBody(src, size-hdr)
			if err != nil {
				return nil, err
			}
			return parseMoov(body)
		}
		pos += size
	}
}

func parseMoov(body []byte) ([]float64, error) {
	var (
		result     []float64
		found      bool
		fragmented bool
		trakErr    error
	)

	err := walkBoxes(body, func(typ string, b []byte) {
		switch typ {
		case "mvex":
			fragmented = true
		case "trak":
			if found {
				return
			}
			times, isVideo, err := parseTrak(b)
			if err != nil {
				trakErr = err
				return
			}
			if isVideo {
				result, found = times, true
			}
		}
	})
	if err != nil {
		return nil, err
	}
	if fragmented {
		return nil, domain.ErrFmp4NotSupported
	}
	if trakErr != nil {
		return nil, trakErr
	}
	if !found {
		return nil, domain.ErrTrackNotFound
	}
	return result, nil
}

func parseTrak(trak []byte) (times []float64, isVideo bool, err error) {
	var (
		timescale uint32
		handler   string
		stbl      []byte
	)

	walkBoxes(trak, func(typ string, b []byte) {
		if typ != "mdia" {
			return
		}
		walkBoxes(b, func(typ string, b []byte) {
			switch typ {
			case "mdhd":
				if len(b) < 24 {
					return
				}
				if b[0] == 1 {
					timescale = binary.BigEndian.Uint32(b[20:24])
				} else {
					timescale = binary.BigEndian.Uint32(b[12:16])
				}
			case "hdlr":
				if len(b) >= 12 {
					handler = string(b[8:12])
				}
			case "minf":
				walkBoxes(b, func(typ string, b []byte) {
					if typ == "stbl" {
						stbl = b
					}
				})
			}
		})
	})

	if handler != "vide" {
		return nil, false, nil
	}
	if timescale == 0 || stbl == nil {
		return nil, true, domain.ErrIncompleteVideoTrack
	}

	var stts, ctts []run
	var stss []uint32
	hasStss := false

	walkBoxes(stbl, func(typ string, b []byte) {
		switch typ {
		case "stts":
			stts = readRuns(b)
		case "ctts":
			ctts = readRuns(b)
		case "stss":
			hasStss = true
			if len(b) < 8 {
				return
			}
			n := int(binary.BigEndian.Uint32(b[4:8]))
			b = b[8:]
			if n > len(b)/4 {
				n = len(b) / 4
			}
			stss = make([]uint32, n)
			for i := range stss {
				stss[i] = binary.BigEndian.Uint32(b[i*4:])
			}
		}
	})

	if len(stts) == 0 {
		return nil, true, domain.ErrNoStssTable
	}

	if !hasStss {
		var total uint32
		for _, r := range stts {
			total += r.count
		}
		stss = make([]uint32, total)
		for i := range stss {
			stss[i] = uint32(i + 1)
		}
	}
	if len(stss) == 0 {
		return nil, true, domain.ErrNoKeyframes
	}
	sort.Slice(stss, func(i, j int) bool { return stss[i] < stss[j] })

	dts := valueAt(stts, stss, true)
	var offsets []int64
	if len(ctts) > 0 {
		offsets = valueAt(ctts, stss, false)
	}

	for i, d := range dts {
		t := d
		if i < len(offsets) {
			t += offsets[i]
		}
		if t < 0 {
			t = 0
		}
		sec := float64(t) / float64(timescale)
		if len(times) == 0 || sec > times[len(times)-1] {
			times = append(times, sec)
		}
	}
	if len(times) == 0 {
		return nil, true, domain.ErrNoKeyframes
	}
	return times, true, nil
}
