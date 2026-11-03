package service

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
)

const maxBoxSize = 1 << 30

func readBox(r *bufio.Reader) (typ string, raw []byte, err error) {
	var h [16]byte
	if _, err = io.ReadFull(r, h[:8]); err != nil {
		return "", nil, err
	}
	size := uint64(binary.BigEndian.Uint32(h[:4]))
	typ = string(h[4:8])
	hdr := uint64(8)
	switch size {
	case 1:
		if _, err = io.ReadFull(r, h[8:16]); err != nil {
			return "", nil, err
		}
		size = binary.BigEndian.Uint64(h[8:16])
		hdr = 16
	case 0:
		return "", nil, domain.ErrBoxWithUnknownSize
	}
	if size < hdr || size > maxBoxSize {
		return "", nil, domain.ErrInvalidData
	}
	raw = make([]byte, size)
	copy(raw, h[:hdr])
	if _, err = io.ReadFull(r, raw[hdr:]); err != nil {
		if errors.Is(err, io.EOF) {
			err = io.ErrUnexpectedEOF
		}
		return "", nil, err
	}
	return typ, raw, nil
}

type fragReader struct {
	br      *bufio.Reader
	pending []byte
}

func (f *fragReader) readInit() ([]byte, error) {
	var init []byte
	haveMoov := false
	for {
		typ, raw, err := readBox(f.br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = domain.ErrMp4NoMoov
			}
			return nil, err
		}
		switch typ {
		case "ftyp", "moov":
			init = append(init, raw...)
			haveMoov = haveMoov || typ == "moov"
		case "moof":
			if !haveMoov {
				return nil, domain.ErrMp4NoMoov
			}
			f.pending = raw
			return init, nil
		}
	}
}

func (f *fragReader) next() ([]byte, error) {
	moof := f.pending
	f.pending = nil
	for moof == nil {
		typ, raw, err := readBox(f.br)
		if err != nil {
			return nil, err
		}
		if typ == "moof" {
			moof = raw
		}
	}
	for {
		typ, raw, err := readBox(f.br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				err = io.ErrUnexpectedEOF
			}
			return nil, err
		}
		if typ == "mdat" {
			frag := make([]byte, 0, len(moof)+len(raw))
			frag = append(frag, moof...)
			return append(frag, raw...), nil
		}
	}
}

func parseInit(init []byte) (videoID, videoTimescale uint32, err error) {
	walkBoxes(init, func(typ string, moov []byte) {
		if typ != "moov" {
			return
		}
		walkBoxes(moov, func(typ string, trak []byte) {
			if typ != "trak" || videoID != 0 {
				return
			}
			var id, ts uint32
			var handler string
			walkBoxes(trak, func(typ string, b []byte) {
				switch typ {
				case "tkhd":
					if len(b) >= 24 && b[0] == 1 {
						id = binary.BigEndian.Uint32(b[20:24])
					} else if len(b) >= 16 {
						id = binary.BigEndian.Uint32(b[12:16])
					}
				case "mdia":
					walkBoxes(b, func(typ string, m []byte) {
						switch typ {
						case "mdhd":
							if len(m) >= 24 && m[0] == 1 {
								ts = binary.BigEndian.Uint32(m[20:24])
							} else if len(m) >= 16 {
								ts = binary.BigEndian.Uint32(m[12:16])
							}
						case "hdlr":
							if len(m) >= 12 {
								handler = string(m[8:12])
							}
						}
					})
				}
			})
			if handler == "vide" && id != 0 && ts != 0 {
				videoID, videoTimescale = id, ts
			}
		})
	})
	if videoID == 0 {
		return 0, 0, domain.ErrTrackNotFound
	}
	return videoID, videoTimescale, nil
}

func fragmentVideoPTS(frag []byte, videoID uint32) (pts int64, ok bool) {
	walkBoxes(frag, func(typ string, moof []byte) {
		if typ != "moof" || ok {
			return
		}
		walkBoxes(moof, func(typ string, traf []byte) {
			if typ != "traf" || ok {
				return
			}
			var id uint32
			var dts, cto int64
			haveTfdt := false
			haveRun := false
			walkBoxes(traf, func(typ string, b []byte) {
				switch typ {
				case "tfhd":
					if len(b) >= 8 {
						id = binary.BigEndian.Uint32(b[4:8])
					}
				case "tfdt":
					if len(b) >= 8 && b[0] == 1 && len(b) >= 12 {
						dts = int64(binary.BigEndian.Uint64(b[4:12]))
						haveTfdt = true
					} else if len(b) >= 8 && b[0] == 0 {
						dts = int64(binary.BigEndian.Uint32(b[4:8]))
						haveTfdt = true
					}
				case "trun":
					if haveRun || len(b) < 8 {
						return
					}
					haveRun = true
					flags := binary.BigEndian.Uint32(b[0:4]) & 0xFFFFFF
					p := 8
					if flags&0x1 != 0 {
						p += 4
					}
					if flags&0x4 != 0 {
						p += 4
					}
					if flags&0x800 != 0 {
						q := p
						for _, bit := range []uint32{0x100, 0x200, 0x400} {
							if flags&bit != 0 {
								q += 4
							}
						}
						if q+4 <= len(b) {
							cto = int64(int32(binary.BigEndian.Uint32(b[q : q+4])))
						}
					}
				}
			})
			if id == videoID && haveTfdt {
				pts, ok = dts+cto, true
			}
		})
	})
	return pts, ok
}

func setSequence(frag []byte, seq uint32) {
	if len(frag) >= 24 && string(frag[4:8]) == "moof" && string(frag[12:16]) == "mfhd" {
		binary.BigEndian.PutUint32(frag[20:24], seq)
	}
}
