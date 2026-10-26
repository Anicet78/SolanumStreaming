package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"math/bits"
	"sort"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
)

// IDs EBML
const (
	idEBML              = 0x1A45DFA3
	idSegment           = 0x18538067
	idSeekHead          = 0x114D9B74
	idSeek              = 0x4DBB
	idSeekID            = 0x53AB
	idSeekPosition      = 0x53AC
	idInfo              = 0x1549A966
	idTimecodeScale     = 0x2AD7B1
	idTracks            = 0x1654AE6B
	idTrackEntry        = 0xAE
	idTrackNumber       = 0xD7
	idTrackType         = 0x83
	idCues              = 0x1C53BB6B
	idCuePoint          = 0xBB
	idCueTime           = 0xB3
	idCueTrackPositions = 0xB7
	idCueTrack          = 0xF7
	idCluster           = 0x1F43B675
)

const maxMasterSize = 64 << 20

func readHeader(r io.Reader) (id uint32, size int64, hdr int, err error) {
	var b [8]byte

	if _, err = io.ReadFull(r, b[:1]); err != nil {
		return
	}
	n := bits.LeadingZeros8(b[0]) + 1
	if n > 4 {
		err = domain.ErrBadEBML
		return
	}
	if n > 1 {
		if _, err = io.ReadFull(r, b[1:n]); err != nil {
			return
		}
	}
	for i := 0; i < n; i++ {
		id = id<<8 | uint32(b[i])
	}
	hdr = n

	if _, err = io.ReadFull(r, b[:1]); err != nil {
		return
	}
	m := bits.LeadingZeros8(b[0]) + 1
	if m > 8 {
		err = domain.ErrBadEBML
		return
	}
	mask := uint64(0xFF) >> m
	val := uint64(b[0]) & mask
	allOnes := val == mask
	if m > 1 {
		if _, err = io.ReadFull(r, b[1:m]); err != nil {
			return
		}
		for i := 1; i < m; i++ {
			val = val<<8 | uint64(b[i])
			allOnes = allOnes && b[i] == 0xFF
		}
	}
	hdr += m

	if allOnes {
		size = -1
	} else {
		size = int64(val)
	}
	return
}

func readUint(b []byte) uint64 {
	var v uint64
	for _, x := range b {
		v = v<<8 | uint64(x)
	}
	return v
}

func walk(data []byte, fn func(id uint32, body []byte)) error {
	r := bytes.NewReader(data)
	for r.Len() > 0 {
		id, size, _, err := readHeader(r)
		if err != nil {
			return err
		}
		if size < 0 || size > int64(r.Len()) {
			return domain.ErrBadEBML
		}
		pos := len(data) - r.Len()
		fn(id, data[pos:pos+int(size)])
		r.Seek(size, io.SeekCurrent)
	}
	return nil
}

func readBody(src io.Reader, size int64) ([]byte, error) {
	if size < 0 || size > maxMasterSize {
		return nil, fmt.Errorf("mkv: element too large (%d)", size)
	}
	buf := make([]byte, size)
	_, err := io.ReadFull(src, buf)
	return buf, err
}


type mkvExtractor struct{}

func (mkvExtractor) Keyframes(ctx context.Context, src io.ReadSeeker, duration float64) ([]float64, error) {
	id, size, _, err := readHeader(src)
	if err != nil || id != idEBML || size < 0 {
		return nil, domain.ErrWrongFile
	}
	if _, err := src.Seek(size, io.SeekCurrent); err != nil {
		return nil, err
	}

	id, _, _, err = readHeader(src)
	if err != nil || id != idSegment {
		return nil, domain.ErrSegmentNotFound
	}
	segStart, err := src.Seek(0, io.SeekCurrent)
	if err != nil {
		return nil, err
	}

	timecodeScale := uint64(1_000_000)
	var videoTrack uint64
	cuesPos := int64(-1)

	pos := segStart
scan:
	for {
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if _, err := src.Seek(pos, io.SeekStart); err != nil {
			return nil, err
		}
		id, size, hdr, err := readHeader(src)
		if err != nil {
			break
		}

		switch id {
		case idCluster:
			break scan
		case idCues:
			cuesPos = pos
		case idSeekHead, idInfo, idTracks:
			body, err := readBody(src, size)
			if err != nil {
				return nil, err
			}
			switch id {
			case idSeekHead:
				if p := cuesFromSeekHead(body); p >= 0 {
					cuesPos = segStart + p
				}
			case idInfo:
				walk(body, func(id uint32, b []byte) {
					if id == idTimecodeScale {
						timecodeScale = readUint(b)
					}
				})
			case idTracks:
				if videoTrack == 0 {
					videoTrack = firstVideoTrack(body)
				}
			}
		}

		if size < 0 {
			break
		}
		pos += int64(hdr) + size
	}

	if videoTrack == 0 {
		return nil, domain.ErrTrackNotFound
	}
	if cuesPos < 0 {
		return nil, domain.ErrNoCues
	}

	if _, err := src.Seek(cuesPos, io.SeekStart); err != nil {
		return nil, err
	}
	id, size, _, err = readHeader(src)
	if err != nil || id != idCues {
		return nil, domain.ErrCuesBadPosition
	}
	body, err := readBody(src, size)
	if err != nil {
		return nil, err
	}

	return parseCues(body, videoTrack, timecodeScale)
}

func cuesFromSeekHead(body []byte) int64 {
	result := int64(-1)
	walk(body, func(id uint32, b []byte) {
		if id != idSeek {
			return
		}
		var seekID uint64
		var seekPos int64 = -1
		walk(b, func(id uint32, v []byte) {
			switch id {
			case idSeekID:
				seekID = readUint(v)
			case idSeekPosition:
				seekPos = int64(readUint(v))
			}
		})
		if seekID == idCues && seekPos >= 0 {
			result = seekPos
		}
	})
	return result
}

func firstVideoTrack(body []byte) uint64 {
	var found uint64
	walk(body, func(id uint32, b []byte) {
		if id != idTrackEntry || found != 0 {
			return
		}
		var number, typ uint64
		walk(b, func(id uint32, v []byte) {
			switch id {
			case idTrackNumber:
				number = readUint(v)
			case idTrackType:
				typ = readUint(v)
			}
		})
		if typ == 1 {
			found = number
		}
	})
	return found
}

func parseCues(body []byte, videoTrack, timecodeScale uint64) ([]float64, error) {
	var times []float64

	err := walk(body, func(id uint32, b []byte) {
		if id != idCuePoint {
			return
		}
		var cueTime uint64
		isVideo := false
		walk(b, func(id uint32, v []byte) {
			switch id {
			case idCueTime:
				cueTime = readUint(v)
			case idCueTrackPositions:
				walk(v, func(id uint32, x []byte) {
					if id == idCueTrack && readUint(x) == videoTrack {
						isVideo = true
					}
				})
			}
		})
		if isVideo {
			times = append(times, float64(cueTime)*float64(timecodeScale)/1e9)
		}
	})
	if err != nil {
		return nil, err
	}
	if len(times) == 0 {
		return nil, domain.ErrCuesEmpty
	}

	sort.Float64s(times)
	out := times[:1]
	for _, t := range times[1:] {
		if t != out[len(out)-1] {
			out = append(out, t)
		}
	}
	return out, nil
}
