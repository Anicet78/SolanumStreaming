package service

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os/exec"
	"sync"
	"time"

	"github.com/Anicet78/SolanumStreaming/stream/internal/domain"
)

const (
	tsOffset = 1.0

	cutTolerance = 0.010

	maxSkew = 0.5

	lookahead  = 3
	reachAhead = 4
	maxJobs    = 3

	maxCacheBytes  = 256 << 20
	jobIdleTimeout = 2 * time.Minute
)

var seekMargins = []float64{0.25, 0.6, 0.05}

type packager struct {
	svc *StreamService
	pl  *Playlist

	mu   sync.Mutex
	cond *sync.Cond

	init    []byte
	skew    float64
	skewSet bool

	cache      map[int][]byte
	cacheBytes int
	jobs       []*job
}

type job struct {
	start, next int
	lastReq     int
	touched     time.Time
	cancel      context.CancelFunc
	cancelled   bool
	done        bool
	err         error
}

func newPackager(svc *StreamService, pl *Playlist) *packager {
	p := &packager{svc: svc, pl: pl, cache: map[int][]byte{}}
	p.cond = sync.NewCond(&p.mu)
	return p
}

func (s *StreamService) GetInit(ctx context.Context, pl *Playlist) ([]byte, error) {
	return pl.pk.getInit(ctx)
}

func (s *StreamService) GetSegment(ctx context.Context, pl *Playlist, index int) ([]byte, error) {
	if index < 0 || index >= len(pl.Segments) {
		return nil, domain.ErrSegmentOutOfRange
	}
	return pl.pk.getSegment(ctx, index)
}

func (p *packager) wakeOn(ctx context.Context) func() bool {
	return context.AfterFunc(ctx, func() {
		p.mu.Lock()
		p.cond.Broadcast()
		p.mu.Unlock()
	})
}

func (p *packager) getInit(ctx context.Context) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	defer p.wakeOn(ctx)()

	for {
		if p.init != nil {
			return p.init, nil
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := p.takeErrLocked(-1); err != nil {
			return nil, err
		}
		if len(p.jobs) == 0 {
			p.startJobLocked(0)
		}
		p.cond.Wait()
	}
}

func (p *packager) getSegment(ctx context.Context, idx int) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	defer p.wakeOn(ctx)()

	for {
		p.touchLocked(idx)
		if b, ok := p.cache[idx]; ok {
			return b, nil
		}
		if err := ctx.Err(); err != nil {
			return nil, err
		}
		if err := p.takeErrLocked(idx); err != nil {
			return nil, err
		}
		if p.findJobLocked(idx) == nil {
			p.startJobLocked(idx)
		}
		p.cond.Wait()
	}
}

func (p *packager) findJobLocked(idx int) *job {
	for _, j := range p.jobs {
		if !j.done && j.err == nil && idx >= j.next && idx <= j.next+reachAhead {
			return j
		}
	}
	return nil
}

func (p *packager) takeErrLocked(idx int) error {
	for _, j := range p.jobs {
		if j.err != nil && (idx < 0 || (idx >= j.start && idx <= j.next+reachAhead)) {
			p.removeJobLocked(j)
			return j.err
		}
	}
	return nil
}

func (p *packager) touchLocked(idx int) {
	for _, j := range p.jobs {
		if j.done || j.err != nil || idx < j.start || idx > j.next+reachAhead {
			continue
		}
		j.touched = time.Now()
		if j.lastReq != idx {
			j.lastReq = idx
			p.cond.Broadcast()
		}
	}
}

func (p *packager) removeJobLocked(rm *job) {
	for i, j := range p.jobs {
		if j == rm {
			p.jobs = append(p.jobs[:i], p.jobs[i+1:]...)
			return
		}
	}
}

func (p *packager) cancelJobLocked(j *job) {
	j.cancelled = true
	j.cancel()
	p.removeJobLocked(j)
	p.cond.Broadcast()
}

func (p *packager) startJobLocked(idx int) {
	for len(p.jobs) >= maxJobs {
		oldest := p.jobs[0]
		for _, j := range p.jobs[1:] {
			if j.touched.Before(oldest.touched) {
				oldest = j
			}
		}
		p.cancelJobLocked(oldest)
	}
	ctx, cancel := context.WithCancel(context.Background())
	j := &job{start: idx, next: idx, lastReq: idx, touched: time.Now(), cancel: cancel}
	p.jobs = append(p.jobs, j)
	go p.run(ctx, j)
}

func (p *packager) storeLocked(idx int, data []byte) {
	if old, ok := p.cache[idx]; ok {
		p.cacheBytes -= len(old)
	}
	p.cache[idx] = data
	p.cacheBytes += len(data)

	for p.cacheBytes > maxCacheBytes && len(p.cache) > 1 {
		worst, worstDist := -1, -1
		for k := range p.cache {
			d := math.MaxInt
			for _, j := range p.jobs {
				d = min(d, abs(k-j.lastReq))
			}
			if len(p.jobs) == 0 {
				d = abs(k - idx)
			}
			if k != idx && d > worstDist {
				worst, worstDist = k, d
			}
		}
		if worst < 0 {
			break
		}
		p.cacheBytes -= len(p.cache[worst])
		delete(p.cache, worst)
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}


func (s *StreamService) ffmpegArgs(pl *Playlist, ss float64) []string {
	args := []string{
		"-v", "error", "-nostdin", "-hide_banner",
		"-seek_timestamp", "1",
		"-ss", fmt.Sprintf("%.3f", max(ss, 0)),
		"-noaccurate_seek",
		"-copyts",
		"-i", s.rawURL(pl.Link),
		"-map", "0:v:0",
	}
	if len(pl.Probe.Audio) > 0 {
		args = append(args, "-map", fmt.Sprintf("0:%d", pl.Probe.Audio[0].Index))
	}
	args = append(args, "-c:v", "copy")
	if pl.Probe.Video.Codec == "hevc" {
		args = append(args, "-tag:v", "hvc1")
	}
	if len(pl.Probe.Audio) > 0 {
		args = append(args, audioArgs(pl.Probe.Audio[0])...)
	}
	args = append(args,
		"-avoid_negative_ts", "disabled",
		"-use_editlist", "0",
		"-output_ts_offset", fmt.Sprintf("%.3f", tsOffset),
		"-movflags", "empty_moov+frag_keyframe+default_base_moof+frag_discont",
		"-f", "mp4",
		"pipe:1",
	)
	return args
}

func audioArgs(a AudioInfo) []string {
	switch a.Codec {
	case "aac", "opus":
		return []string{"-c:a", "copy"}
	default:
		return []string{"-c:a", "aac", "-b:a", "256k", "-ac", "2"}
	}
}

type tailBuffer struct {
	mu  sync.Mutex
	buf []byte
}

func (t *tailBuffer) Write(b []byte) (int, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.buf = append(t.buf, b...)
	if len(t.buf) > 4096 {
		t.buf = t.buf[len(t.buf)-4096:]
	}
	return len(b), nil
}

func (t *tailBuffer) String() string {
	t.mu.Lock()
	defer t.mu.Unlock()
	return string(t.buf)
}

func (p *packager) run(ctx context.Context, j *job) {
	err := p.stream(ctx, j)

	p.mu.Lock()
	defer p.mu.Unlock()
	j.done = true
	if err != nil && !j.cancelled {
		j.err = err
		slog.Error("ffmpeg job failed", "link", p.pl.ID, "start", j.start, "next", j.next, "error", err)
	}
	p.cond.Broadcast()
}

func (p *packager) stream(ctx context.Context, j *job) error {
	err := domain.ErrSeekMismatch
	for _, margin := range seekMargins {
		if err = p.streamOnce(ctx, j, margin); !errors.Is(err, domain.ErrSeekMismatch) {
			return err
		}
		slog.Warn("seek ffmpeg sur le mauvais keyframe, nouvel essai", "link", p.pl.ID, "segment", j.start, "margin", margin)
	}
	return err
}

func (p *packager) streamOnce(parent context.Context, j *job, margin float64) error {
	ctx, cancelAttempt := context.WithCancel(parent)
	defer cancelAttempt()

	segs := p.pl.Segments
	last := len(segs) - 1

	p.mu.Lock()
	skew := 0.0
	if p.skewSet {
		skew = p.skew
	}
	p.mu.Unlock()

	ss := segs[j.start].Start - skew + margin
	cmd := exec.CommandContext(ctx, "ffmpeg", p.svc.ffmpegArgs(p.pl, ss)...)
	stderr := &tailBuffer{}
	cmd.Stderr = stderr
	cmd.WaitDelay = 5 * time.Second
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("ffmpeg: %w", err)
	}
	defer func() {
		cancelAttempt()
		_ = cmd.Wait()
	}()

	wrap := func(err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		return fmt.Errorf("%w: %s", err, stderr.String())
	}

	fr := &fragReader{br: bufio.NewReaderSize(stdout, 1<<20)}
	initBytes, err := fr.readInit()
	if err != nil {
		return wrap(err)
	}
	videoID, videoTS, err := parseInit(initBytes)
	if err != nil {
		return wrap(err)
	}
	p.mu.Lock()
	if p.init == nil {
		p.init = initBytes
		p.cond.Broadcast()
	}
	p.mu.Unlock()

	var (
		cur       [][]byte
		curBytes  int
		idx       = j.start
		first     = true
		jobSkew   float64
		warnedNeg bool
	)

	flush := func() {
		data := make([]byte, 0, curBytes)
		for _, f := range cur {
			data = append(data, f...)
		}
		cur, curBytes = nil, 0

		p.mu.Lock()
		p.storeLocked(idx, data)
		idx++
		j.next = idx
		p.cond.Broadcast()
		p.pauseLocked(j)
		p.mu.Unlock()
	}

	for {
		frag, err := fr.next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return wrap(err)
		}

		pts, ok := fragmentVideoPTS(frag, videoID)
		sec := float64(pts)/float64(videoTS) - tsOffset
		if ok && pts < 0 && !warnedNeg {
			warnedNeg = true
			slog.Warn("DTS négatif en sortie : augmenter tsOffset", "link", p.pl.ID)
		}

		switch {
		case first:
			first = false
			if !ok {
				return wrap(errors.New("premier fragment sans vidéo"))
			}
			jobSkew = segs[idx].Start - sec
			p.mu.Lock()
			if (p.skewSet && math.Abs(jobSkew-p.skew) > 0.05) || math.Abs(jobSkew) > maxSkew {
				p.mu.Unlock()
				return domain.ErrSeekMismatch
			}
			if !p.skewSet {
				p.skew, p.skewSet = jobSkew, true
			}
			p.mu.Unlock()
		case ok && idx < last && sec >= segs[idx+1].Start-jobSkew-cutTolerance:
			flush()
			p.mu.Lock()
			cancelled := j.cancelled
			p.mu.Unlock()
			if cancelled {
				return nil
			}
		}

		setSequence(frag, uint32(idx)<<12|uint32(len(cur)+1))
		cur = append(cur, frag)
		curBytes += len(frag)
	}

	if err := cmd.Wait(); err != nil && ctx.Err() == nil {
		return fmt.Errorf("ffmpeg: %w: %s", err, stderr.String())
	}
	if idx != last {
		return fmt.Errorf("flux terminé trop tôt (segment %d sur %d)", idx, last)
	}
	if len(cur) > 0 {
		flush()
	}
	return nil
}

func (p *packager) pauseLocked(j *job) {
	if j.cancelled || j.next <= j.lastReq+lookahead {
		return
	}
	var t *time.Timer
	t = time.AfterFunc(jobIdleTimeout, func() {
		p.mu.Lock()
		defer p.mu.Unlock()
		if j.cancelled || j.done {
			return
		}
		if rest := jobIdleTimeout - time.Since(j.touched); rest > 0 {
			t.Reset(rest)
			return
		}
		p.cancelJobLocked(j)
	})
	defer t.Stop()
	for !j.cancelled && j.next > j.lastReq+lookahead {
		p.cond.Wait()
	}
}
