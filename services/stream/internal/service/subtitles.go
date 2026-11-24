package service

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os/exec"
	"path"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

const (
	subWindow      = 60.0
	subPreroll     = 10.0
	maxSidecarSize = 5 << 20
	maxSidecars    = 30
	subCacheMax    = 128
	subTimeout     = 60 * time.Second
)

var ErrSubtitleNotFound = errors.New("subtitle track or segment not found")

type cue struct {
	start, end float64
	settings   string
	text       string
}

type sidecarFile struct {
	Name string
	Size int64
	Open func() (io.ReadSeekCloser, error)
	Fetch func(ctx context.Context) ([]byte, error)
	Label string
	Lang  string
}

func (sc *sidecarFile) external() bool { return sc.Label != "" }

type subTrack struct {
	ID     int
	Name   string
	Lang   string
	Forced bool

	stream  int
	sidecar *sidecarFile

	mu     sync.Mutex
	loaded bool
	cues   []cue
}

type vttEntry struct {
	done chan struct{}
	data []byte
	err  error
}

type subtitles struct {
	svc    *StreamService
	pl     *Playlist
	tracks []*subTrack
	sem    chan struct{}

	mu    sync.Mutex
	cache map[string]*vttEntry
	order []string
}

func newSubtitles(s *StreamService, pl *Playlist, sidecars []sidecarFile) *subtitles {
	st := &subtitles{svc: s, pl: pl, sem: make(chan struct{}, 2), cache: map[string]*vttEntry{}}

	for _, si := range pl.Probe.Subs {
		iso, lname := lookupLang(si.Language)
		name := si.Title
		if name == "" {
			name = lname
		}
		if name == "" {
			name = fmt.Sprintf("Sous-titres %d", len(st.tracks)+1)
		}
		if si.Forced && !strings.Contains(strings.ToLower(name), "forc") {
			name += " (forcés)"
		}
		st.tracks = append(st.tracks, &subTrack{
			Name: name, Lang: langTag(si.Language, iso), Forced: si.Forced, stream: si.Index,
		})
	}

	sort.SliceStable(sidecars, func(i, j int) bool {
		ei, ej := sidecars[i].external(), sidecars[j].external()
		if ei != ej {
			return !ei
		}
		return !ei && sidecars[i].Name < sidecars[j].Name
	})
	nTorrent := 0
	for i := range sidecars {
		sc := sidecars[i]
		if sc.external() {
			st.tracks = append(st.tracks, &subTrack{Name: sc.Label, Lang: sc.Lang, stream: -1, sidecar: &sc})
			continue
		}
		if nTorrent++; nTorrent > maxSidecars {
			continue
		}
		iso, lname := detectLangFromPath(sc.Name)
		name := lname
		if name == "" {
			name = strings.TrimSuffix(path.Base(sc.Name), path.Ext(sc.Name))
		}
		st.tracks = append(st.tracks, &subTrack{Name: name, Lang: iso, stream: -1, sidecar: &sc})
	}

	seen := map[string]int{}
	for i, t := range st.tracks {
		t.ID = i
		seen[t.Name]++
		if n := seen[t.Name]; n > 1 {
			t.Name = fmt.Sprintf("%s (%d)", t.Name, n)
		}
	}
	return st
}

func (s *StreamService) SubtitlePlaylist(pl *Playlist, id int) (string, bool) {
	return pl.subs.mediaPlaylist(id)
}

func (s *StreamService) SubtitleSegment(ctx context.Context, pl *Playlist, id, k int) ([]byte, error) {
	return pl.subs.segment(ctx, id, k)
}

func hlsAttr(s string) string { return strings.NewReplacer(`"`, "'", "\n", " ", "\r", " ").Replace(s) }

func (st *subtitles) master() string {
	if len(st.tracks) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("#EXTM3U\n#EXT-X-VERSION:7\n#EXT-X-INDEPENDENT-SEGMENTS\n")
	for _, t := range st.tracks {
		fmt.Fprintf(&b, "#EXT-X-MEDIA:TYPE=SUBTITLES,GROUP-ID=\"subs\",NAME=\"%s\",", hlsAttr(t.Name))
		if t.Lang != "" {
			fmt.Fprintf(&b, "LANGUAGE=\"%s\",", hlsAttr(t.Lang))
		}
		forced := "NO"
		if t.Forced {
			forced = "YES"
		}
		fmt.Fprintf(&b, "DEFAULT=NO,AUTOSELECT=NO,FORCED=%s,URI=\"sub_%d.m3u8\"\n", forced, t.ID)
	}
	bw := st.pl.Probe.Bitrate
	if bw <= 0 {
		bw = 8_000_000
	}
	fmt.Fprintf(&b, "#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%dx%d,SUBTITLES=\"subs\"\nvideo.m3u8\n",
		bw, st.pl.Probe.Video.Width, st.pl.Probe.Video.Height)
	return b.String()
}

func (st *subtitles) windows() int {
	return max(1, int(math.Ceil(st.pl.Probe.Duration/subWindow)))
}

func (st *subtitles) mediaPlaylist(id int) (string, bool) {
	if id < 0 || id >= len(st.tracks) {
		return "", false
	}
	n := st.windows()
	total := st.pl.Probe.Duration

	var b strings.Builder
	b.WriteString("#EXTM3U\n#EXT-X-VERSION:7\n")
	fmt.Fprintf(&b, "#EXT-X-TARGETDURATION:%d\n", int(subWindow))
	b.WriteString("#EXT-X-PLAYLIST-TYPE:VOD\n#EXT-X-MEDIA-SEQUENCE:0\n")
	for k := 0; k < n; k++ {
		d := subWindow
		if k == n-1 && total > 0 {
			d = total - float64(k)*subWindow
		}
		fmt.Fprintf(&b, "#EXTINF:%.3f,\nsub_%d_%d.vtt\n", d, id, k)
	}
	b.WriteString("#EXT-X-ENDLIST\n")
	return b.String(), true
}

func (st *subtitles) segment(ctx context.Context, id, k int) ([]byte, error) {
	if id < 0 || id >= len(st.tracks) || k < 0 || k >= st.windows() {
		return nil, ErrSubtitleNotFound
	}
	tr := st.tracks[id]
	ws := float64(k) * subWindow
	we := ws + subWindow

	return st.cached(ctx, fmt.Sprintf("%d:%d", id, k), func(ctx context.Context) ([]byte, error) {
		var all []cue
		var err error
		if tr.sidecar != nil {
			all, err = st.loadSidecar(ctx, tr)
		} else {
			all, err = st.extract(ctx, tr, ws, we)
		}
		if err != nil {
			return nil, err
		}
		var in []cue
		for _, c := range all {
			if c.end > ws && c.start < we {
				in = append(in, c)
			}
		}
		return renderVTT(in), nil
	})
}

func (st *subtitles) cached(ctx context.Context, key string, fn func(context.Context) ([]byte, error)) ([]byte, error) {
	st.mu.Lock()
	e, ok := st.cache[key]
	if !ok {
		e = &vttEntry{done: make(chan struct{})}
		st.cache[key] = e
		st.order = append(st.order, key)
		for len(st.order) > subCacheMax {
			delete(st.cache, st.order[0])
			st.order = st.order[1:]
		}
		go func() {
			cctx, cancel := context.WithTimeout(context.Background(), subTimeout)
			defer cancel()
			e.data, e.err = fn(cctx)
			if e.err != nil {
				slog.Error("subtitle segment failed", "link", st.pl.ID, "key", key, "error", e.err)
				st.mu.Lock()
				if st.cache[key] == e {
					delete(st.cache, key)
				}
				st.mu.Unlock()
			}
			close(e.done)
		}()
	}
	st.mu.Unlock()

	select {
	case <-e.done:
		return e.data, e.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (st *subtitles) extract(ctx context.Context, tr *subTrack, ws, we float64) ([]cue, error) {
	select {
	case st.sem <- struct{}{}:
		defer func() { <-st.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	args := []string{
		"-v", "error", "-nostdin", "-hide_banner",
		"-seek_timestamp", "1",
		"-ss", fmt.Sprintf("%.3f", max(ws-subPreroll, 0)),
		"-copyts",
		"-i", st.svc.rawURL(st.pl.Link),
		"-to", fmt.Sprintf("%.3f", we),
		"-map", fmt.Sprintf("0:%d", tr.stream),
		"-f", "webvtt", "pipe:1",
	}
	out, err := runFFmpeg(ctx, nil, args)
	if err != nil {
		return nil, err
	}
	return parseVTT(out), nil
}

func runFFmpeg(ctx context.Context, stdin []byte, args []string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout, cmd.Stderr = &stdout, &stderr
	if stdin != nil {
		cmd.Stdin = bytes.NewReader(stdin)
	}
	if err := cmd.Run(); err != nil {
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
		return nil, fmt.Errorf("ffmpeg: %w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return stdout.Bytes(), nil
}

func sidecarFormat(name string) (format string, ok bool) {
	switch strings.ToLower(path.Ext(name)) {
	case ".srt":
		return "srt", true
	case ".ass", ".ssa":
		return "ass", true
	case ".vtt":
		return "webvtt", true
	}
	return "", false
}

func isSidecarName(name string) bool { _, ok := sidecarFormat(name); return ok }

func (st *subtitles) loadSidecar(ctx context.Context, tr *subTrack) ([]cue, error) {
	tr.mu.Lock()
	defer tr.mu.Unlock()
	if tr.loaded {
		return tr.cues, nil
	}

	data, err := readLimited(ctx, tr.sidecar)
	if err != nil {
		return nil, err
	}
	format, _ := sidecarFormat(tr.sidecar.Name)

	var cues []cue
	if format == "webvtt" {
		cues = parseVTT(decodeText(data))
	} else {
		args := []string{"-v", "error", "-nostdin", "-hide_banner"}
		if enc := charsetFlag(data); enc != "" {
			args = append(args, "-sub_charenc", enc)
		}
		args = append(args, "-f", format, "-i", "pipe:0", "-f", "webvtt", "pipe:1")
		out, err := runFFmpeg(ctx, data, args)
		if err != nil {
			return nil, err
		}
		cues = parseVTT(out)
	}
	tr.cues, tr.loaded = cues, true
	return cues, nil
}

func readLimited(ctx context.Context, sc *sidecarFile) ([]byte, error) {
	type result struct {
		data []byte
		err  error
	}
	ch := make(chan result, 1)
	go func() {
		if sc.Fetch != nil {
			data, err := sc.Fetch(ctx)
			if err == nil && len(data) > maxSidecarSize {
				err = fmt.Errorf("fichier de sous-titres trop gros (%s)", sc.Name)
			}
			ch <- result{data, err}
			return
		}
		r, err := sc.Open()
		if err != nil {
			ch <- result{nil, err}
			return
		}
		defer r.Close()
		data, err := io.ReadAll(io.LimitReader(r, maxSidecarSize+1))
		if err == nil && len(data) > maxSidecarSize {
			err = fmt.Errorf("fichier de sous-titres trop gros (%s)", sc.Name)
		}
		ch <- result{data, err}
	}()
	select {
	case r := <-ch:
		return r.data, r.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func charsetFlag(data []byte) string {
	if bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) || bytes.HasPrefix(data, []byte{0xFF, 0xFE}) || bytes.HasPrefix(data, []byte{0xFE, 0xFF}) {
		return ""
	}
	if utf8.Valid(data) {
		return ""
	}
	return "CP1252"
}

func decodeText(data []byte) []byte {
	data = bytes.TrimPrefix(data, []byte{0xEF, 0xBB, 0xBF})
	if utf8.Valid(data) {
		return data
	}
	cp1252 := [32]rune{0x20AC, 0x81, 0x201A, 0x0192, 0x201E, 0x2026, 0x2020, 0x2021, 0x02C6, 0x2030, 0x0160, 0x2039, 0x0152, 0x8D, 0x017D, 0x8F,
		0x90, 0x2018, 0x2019, 0x201C, 0x201D, 0x2022, 0x2013, 0x2014, 0x02DC, 0x2122, 0x0161, 0x203A, 0x0153, 0x9D, 0x017E, 0x0178}
	var b strings.Builder
	for _, c := range data {
		switch {
		case c >= 0x80 && c < 0xA0:
			b.WriteRune(cp1252[c-0x80])
		default:
			b.WriteRune(rune(c))
		}
	}
	return []byte(b.String())
}

func parseVTTTime(s string) (float64, bool) {
	parts := strings.Split(strings.TrimSpace(s), ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}
	var sec float64
	for _, p := range parts {
		v, err := strconv.ParseFloat(strings.ReplaceAll(p, ",", "."), 64)
		if err != nil || v < 0 {
			return 0, false
		}
		sec = sec*60 + v
	}
	return sec, true
}

func parseVTT(data []byte) []cue {
	text := strings.ReplaceAll(string(data), "\r\n", "\n")
	text = strings.ReplaceAll(text, "\r", "\n")

	var cues []cue
	for _, block := range strings.Split(text, "\n\n") {
		lines := strings.Split(strings.Trim(block, "\n"), "\n")
		for i, line := range lines {
			if !strings.Contains(line, "-->") {
				continue
			}
			left, right, _ := strings.Cut(line, "-->")
			fields := strings.Fields(right)
			if len(fields) == 0 {
				break
			}
			start, ok1 := parseVTTTime(left)
			end, ok2 := parseVTTTime(fields[0])
			if !ok1 || !ok2 {
				break
			}
			body := strings.TrimSpace(strings.Join(lines[i+1:], "\n"))
			if body != "" && end > start {
				cues = append(cues, cue{start: start, end: end, settings: strings.Join(fields[1:], " "), text: body})
			}
			break
		}
	}
	return cues
}

func fmtVTTTime(sec float64) string {
	ms := int64(math.Round(max(sec, 0) * 1000))
	return fmt.Sprintf("%02d:%02d:%02d.%03d", ms/3600000, ms/60000%60, ms/1000%60, ms%1000)
}

func renderVTT(cues []cue) []byte {
	var b strings.Builder
	fmt.Fprintf(&b, "WEBVTT\nX-TIMESTAMP-MAP=MPEGTS:%d,LOCAL:00:00:00.000\n\n", int(math.Round(tsOffset*90000)))
	for _, c := range cues {
		fmt.Fprintf(&b, "%s --> %s", fmtVTTTime(c.start), fmtVTTTime(c.end))
		if c.settings != "" {
			b.WriteString(" " + c.settings)
		}
		b.WriteString("\n" + c.text + "\n\n")
	}
	return []byte(b.String())
}

type langEntry struct {
	iso1, name string
	aliases    []string
}

var languages = []langEntry{
	{"en", "English", []string{"eng", "english", "anglais"}},
	{"fr", "Français", []string{"fre", "fra", "french", "francais", "français"}},
	{"es", "Español", []string{"spa", "spanish", "espanol", "español", "castellano", "latino"}},
	{"de", "Deutsch", []string{"ger", "deu", "german", "deutsch", "allemand"}},
	{"it", "Italiano", []string{"ita", "italian", "italiano"}},
	{"pt", "Português", []string{"por", "portuguese", "português", "portugues", "brazilian", "ptbr"}},
	{"nl", "Nederlands", []string{"dut", "nld", "dutch", "nederlands"}},
	{"ru", "Русский", []string{"rus", "russian"}},
	{"ja", "日本語", []string{"jpn", "japanese"}},
	{"ko", "한국어", []string{"kor", "korean"}},
	{"zh", "中文", []string{"chi", "zho", "chinese", "mandarin", "cantonese"}},
	{"ar", "العربية", []string{"ara", "arabic"}},
	{"tr", "Türkçe", []string{"tur", "turkish"}},
	{"pl", "Polski", []string{"pol", "polish", "polski"}},
	{"sv", "Svenska", []string{"swe", "swedish", "svenska"}},
	{"da", "Dansk", []string{"dan", "danish", "dansk"}},
	{"no", "Norsk", []string{"nor", "nob", "norwegian", "norsk"}},
	{"fi", "Suomi", []string{"fin", "finnish", "suomi"}},
	{"cs", "Čeština", []string{"cze", "ces", "czech"}},
	{"el", "Ελληνικά", []string{"gre", "ell", "greek"}},
	{"he", "עברית", []string{"heb", "hebrew"}},
	{"hu", "Magyar", []string{"hun", "hungarian"}},
	{"ro", "Română", []string{"rum", "ron", "romanian"}},
	{"uk", "Українська", []string{"ukr", "ukrainian"}},
	{"vi", "Tiếng Việt", []string{"vie", "vietnamese"}},
	{"th", "ไทย", []string{"tha", "thai"}},
	{"id", "Bahasa Indonesia", []string{"ind", "indonesian"}},
	{"hi", "हिन्दी", []string{"hin", "hindi"}},
}

var langIndex = func() map[string]*langEntry {
	m := map[string]*langEntry{}
	for i := range languages {
		l := &languages[i]
		m[l.iso1] = l
		for _, a := range l.aliases {
			m[a] = l
		}
	}
	return m
}()

func lookupLang(code string) (iso1, name string) {
	if l, ok := langIndex[strings.ToLower(strings.TrimSpace(code))]; ok {
		return l.iso1, l.name
	}
	return "", ""
}

func langTag(raw, iso1 string) string {
	if iso1 != "" {
		return iso1
	}
	raw = strings.ToLower(strings.TrimSpace(raw))
	if len(raw) == 3 && raw != "und" && raw != "mul" && raw != "zxx" {
		return raw
	}
	return ""
}

var reLangRegion = regexp.MustCompile(`(?i)(?:^|[._ ])([a-z]{2})[-_][a-z]{2}$`)

func isTrailingCode(toks []string, i int) bool {
	if i == len(toks)-1 {
		return true
	}
	if i == len(toks)-2 {
		switch toks[i+1] {
		case "forced", "sdh", "cc", "hi", "full", "sub", "subs", "subtitle", "subtitles":
			return true
		}
	}
	return false
}

func detectLangFromPath(p string) (iso1, name string) {
	base := strings.TrimSuffix(path.Base(p), path.Ext(p))
	split := func(s string) []string {
		return strings.FieldsFunc(strings.ToLower(s), func(r rune) bool { return !unicode.IsLetter(r) })
	}
	if m := reLangRegion.FindStringSubmatch(base); m != nil {
		if l, ok := langIndex[strings.ToLower(m[1])]; ok {
			return l.iso1, l.name
		}
	}
	toks := split(base)
	for i, t := range toks {
		if len(t) == 2 && !isTrailingCode(toks, i) {
			continue
		}
		if l, ok := langIndex[t]; ok {
			return langTag("", l.iso1), l.name
		}
	}
	for _, t := range split(path.Dir(p)) {
		if len(t) > 2 {
			if l, ok := langIndex[t]; ok {
				return l.iso1, l.name
			}
		}
	}
	return "", ""
}
