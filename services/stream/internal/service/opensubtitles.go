package service

import (
	"bytes"
	"context"
	"crypto/sha1"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	defaultOSBase      = "https://api.opensubtitles.com/api/v1"
	osSearchTTL        = 6 * time.Hour
	osDiskPickTTL      = 7 * 24 * time.Hour
	osTokenTTL         = 20 * time.Hour
	osMaxMemSubs       = 64
	osLookupBudget     = 12 * time.Second
	osHashBudget       = 6 * time.Second
	maxOSResponseBytes = 4 << 20
)

var ErrSubtitleQuota = errors.New("opensubtitles: downloading quota reached")

type MovieMeta struct {
	TmdbID int
	ImdbID string
}

var reImdb = regexp.MustCompile(`^(?:tt)?0*([0-9]{1,10})$`)

func ParseMovieMeta(tmdb, imdb string) MovieMeta {
	var m MovieMeta
	if n, err := strconv.Atoi(strings.TrimSpace(tmdb)); err == nil && n > 0 {
		m.TmdbID = n
	}
	if sm := reImdb.FindStringSubmatch(strings.TrimSpace(imdb)); sm != nil {
		m.ImdbID = sm[1]
	}
	return m
}

func (m MovieMeta) key() string {
	var parts []string
	if m.TmdbID > 0 {
		parts = append(parts, "tmdb="+strconv.Itoa(m.TmdbID))
	}
	if m.ImdbID != "" {
		parts = append(parts, "imdb="+m.ImdbID)
	}
	if len(parts) == 0 {
		return ""
	}
	return "|" + strings.Join(parts, "&")
}

type OpenSubtitlesConfig struct {
	APIKey    string
	Username  string
	Password  string
	UserAgent string
	Languages []string
	CacheDir  string
	BaseURL   string
}

func OpenSubtitlesConfigFromEnv() (OpenSubtitlesConfig, bool) {
	cfg := OpenSubtitlesConfig{
		APIKey:    strings.TrimSpace(os.Getenv("OPENSUBTITLES_API_KEY")),
		Username:  strings.TrimSpace(os.Getenv("OPENSUBTITLES_USERNAME")),
		Password:  os.Getenv("OPENSUBTITLES_PASSWORD"),
		UserAgent: strings.TrimSpace(os.Getenv("OPENSUBTITLES_USER_AGENT")),
		CacheDir:  strings.TrimSpace(os.Getenv("SUBTITLE_CACHE_DIR")),
	}
	if cfg.APIKey == "" {
		return cfg, false
	}
	if cfg.UserAgent == "" {
		cfg.UserAgent = "SolanumStreaming v1.0"
	}
	if cfg.CacheDir == "" {
		cfg.CacheDir = filepath.Join(os.TempDir(), "solanum-subtitles")
	}
	langs := os.Getenv("SUBTITLE_LANGS")
	if strings.TrimSpace(langs) == "" {
		langs = "fr"
	}
	for _, l := range strings.Split(langs, ",") {
		if l = strings.ToLower(strings.TrimSpace(l)); l != "" {
			cfg.Languages = append(cfg.Languages, l)
		}
	}
	return cfg, true
}

type openSubtitles struct {
	cfg  OpenSubtitlesConfig
	http *http.Client

	dl sync.Mutex

	mu       sync.Mutex
	token    string
	tokenAt  time.Time
	apiBase  string
	blocked  time.Time
	mem      map[int][]byte
	memOrder []int
	picks    map[string]cachedPicks
}

type cachedPicks struct {
	at    time.Time
	picks []osPick
}

func newOpenSubtitles(cfg OpenSubtitlesConfig) *openSubtitles {
	if cfg.UserAgent == "" {
		cfg.UserAgent = "SolanumStreaming v1.0"
	}
	return &openSubtitles{
		cfg:   cfg,
		http:  &http.Client{Timeout: 20 * time.Second},
		mem:   map[int][]byte{},
		picks: map[string]cachedPicks{},
	}
}

func (o *openSubtitles) base() string {
	o.mu.Lock()
	defer o.mu.Unlock()
	if o.apiBase != "" {
		return o.apiBase
	}
	if o.cfg.BaseURL != "" {
		return strings.TrimRight(o.cfg.BaseURL, "/")
	}
	return defaultOSBase
}

func (o *openSubtitles) request(ctx context.Context, method, rawURL string, body any, withToken bool) (int, []byte, error) {
	var rd io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return 0, nil, err
		}
		rd = bytes.NewReader(b)
	}
	req, err := http.NewRequestWithContext(ctx, method, rawURL, rd)
	if err != nil {
		return 0, nil, err
	}
	req.Header.Set("Api-Key", o.cfg.APIKey)
	req.Header.Set("User-Agent", o.cfg.UserAgent)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if withToken {
		o.mu.Lock()
		tok := o.token
		o.mu.Unlock()
		if tok != "" {
			req.Header.Set("Authorization", "Bearer "+tok)
		}
	}
	resp, err := o.http.Do(req)
	if err != nil {
		return 0, nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(io.LimitReader(resp.Body, maxOSResponseBytes))
	return resp.StatusCode, data, err
}

func (o *openSubtitles) login(ctx context.Context, force bool) error {
	if o.cfg.Username == "" || o.cfg.Password == "" {
		return nil
	}
	o.mu.Lock()
	ok := !force && o.token != "" && time.Since(o.tokenAt) < osTokenTTL
	o.mu.Unlock()
	if ok {
		return nil
	}
	status, data, err := o.request(ctx, http.MethodPost, o.base()+"/login",
		map[string]string{"username": o.cfg.Username, "password": o.cfg.Password}, false)
	if err != nil {
		return err
	}
	if status != http.StatusOK {
		return fmt.Errorf("opensubtitles login: HTTP %d: %s", status, truncate(data))
	}
	var r struct {
		Token   string `json:"token"`
		BaseURL string `json:"base_url"`
	}
	if err := json.Unmarshal(data, &r); err != nil || r.Token == "" {
		return fmt.Errorf("opensubtitles login: unexpected reply: %s", truncate(data))
	}
	o.mu.Lock()
	o.token, o.tokenAt = r.Token, time.Now()
	if r.BaseURL != "" && o.cfg.BaseURL == "" {
		o.apiBase = "https://" + strings.TrimPrefix(strings.TrimPrefix(r.BaseURL, "https://"), "http://") + "/api/v1"
	}
	o.mu.Unlock()
	return nil
}

func truncate(b []byte) string {
	s := strings.TrimSpace(string(b))
	if len(s) > 200 {
		s = s[:200] + "…"
	}
	return s
}

type osPick struct {
	FileID    int    `json:"file_id"`
	Lang      string `json:"lang"`
	Release   string `json:"release"`
	HashMatch bool   `json:"hash_match"`
}

type osAttrs struct {
	Language          string  `json:"language"`
	DownloadCount     int     `json:"download_count"`
	HearingImpaired   bool    `json:"hearing_impaired"`
	FromTrusted       bool    `json:"from_trusted"`
	ForeignPartsOnly  bool    `json:"foreign_parts_only"`
	AITranslated      bool    `json:"ai_translated"`
	MachineTranslated bool    `json:"machine_translated"`
	MovieHashMatch    bool    `json:"moviehash_match"`
	Ratings           float64 `json:"ratings"`
	Release           string  `json:"release"`
	Files             []struct {
		FileID int `json:"file_id"`
	} `json:"files"`
}

type osSearchResponse struct {
	Data []struct {
		Attributes osAttrs `json:"attributes"`
	} `json:"data"`
}

func baseLang(l string) string {
	l = strings.ToLower(strings.TrimSpace(l))
	if i := strings.IndexAny(l, "-_"); i > 0 {
		l = l[:i]
	}
	return l
}

func scoreCandidate(a osAttrs) float64 {
	s := math.Log10(1 + float64(max(a.DownloadCount, 0)))
	if a.MovieHashMatch {
		s += 100
	}
	if a.FromTrusted {
		s += 3
	}
	if !a.HearingImpaired {
		s += 2
	}
	if a.ForeignPartsOnly {
		s -= 50
	}
	if a.AITranslated || a.MachineTranslated {
		s -= 20
	}
	s += math.Min(math.Max(a.Ratings, 0), 10) / 10
	return s
}

func pickBest(resp osSearchResponse, langs []string) []osPick {
	var out []osPick
	for _, want := range langs {
		bestScore, best := math.Inf(-1), osPick{}
		for _, d := range resp.Data {
			a := d.Attributes
			if baseLang(a.Language) != want || len(a.Files) != 1 || a.Files[0].FileID <= 0 {
				continue
			}
			if s := scoreCandidate(a); s > bestScore {
				bestScore = s
				best = osPick{FileID: a.Files[0].FileID, Lang: want, Release: a.Release, HashMatch: a.MovieHashMatch}
			}
		}
		if best.FileID > 0 {
			out = append(out, best)
		}
	}
	return out
}

func (o *openSubtitles) searchOnce(ctx context.Context, meta MovieMeta, hash string, langs []string, strict bool) (osSearchResponse, int, error) {
	q := url.Values{}
	q.Set("languages", strings.Join(langs, ","))
	q.Set("type", "movie")
	if meta.TmdbID > 0 {
		q.Set("tmdb_id", strconv.Itoa(meta.TmdbID))
	}
	if meta.ImdbID != "" {
		q.Set("imdb_id", meta.ImdbID)
	}
	if hash != "" {
		q.Set("moviehash", hash)
	}
	if strict {
		q.Set("ai_translated", "exclude")
		q.Set("machine_translated", "exclude")
		q.Set("foreign_parts_only", "exclude")
		q.Set("order_by", "download_count")
		q.Set("order_direction", "desc")
	}
	status, data, err := o.request(ctx, http.MethodGet, o.base()+"/subtitles?"+q.Encode(), nil, false)
	var resp osSearchResponse
	if err != nil {
		return resp, 0, err
	}
	if status != http.StatusOK {
		return resp, status, fmt.Errorf("opensubtitles search: HTTP %d: %s", status, truncate(data))
	}
	return resp, status, json.Unmarshal(data, &resp)
}

func (o *openSubtitles) find(ctx context.Context, meta MovieMeta, hash string, langs []string) ([]osPick, error) {
	if len(langs) == 0 || (meta.key() == "" && hash == "") {
		return nil, nil
	}
	key := meta.key() + "|" + hash + "|" + strings.Join(langs, ",")

	o.mu.Lock()
	c, ok := o.picks[key]
	o.mu.Unlock()
	if ok && time.Since(c.at) < osSearchTTL {
		return c.picks, nil
	}

	resp, status, err := o.searchOnce(ctx, meta, hash, langs, true)
	if err != nil && (status == http.StatusBadRequest || status == http.StatusUnprocessableEntity) {
		resp, status, err = o.searchOnce(ctx, meta, hash, langs, false)
	}
	if err != nil {
		if p, ok := o.loadPicks(key, 0); ok {
			return p, nil
		}
		return nil, err
	}
	picks := pickBest(resp, langs)
	o.mu.Lock()
	o.picks[key] = cachedPicks{at: time.Now(), picks: picks}
	o.mu.Unlock()
	o.savePicks(key, picks)
	return picks, nil
}

type queryTime struct {
	At    time.Time `json:"at"`
	Picks []osPick  `json:"picks"`
}

func (o *openSubtitles) path(kind, key string) string {
	h := sha1.Sum([]byte(key))
	return filepath.Join(o.cfg.CacheDir, kind+"_"+hex.EncodeToString(h[:10]))
}

func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func (o *openSubtitles) savePicks(key string, picks []osPick) {
	if o.cfg.CacheDir == "" {
		return
	}
	b, _ := json.Marshal(queryTime{At: time.Now(), Picks: picks})
	if err := writeAtomic(o.path("pick", key), b); err != nil {
		slog.Warn("opensubtitles: unavailable disk cache", "error", err)
	}
}

func (o *openSubtitles) loadPicks(key string, maxAge time.Duration) ([]osPick, bool) {
	if o.cfg.CacheDir == "" {
		return nil, false
	}
	b, err := os.ReadFile(o.path("pick", key))
	if err != nil {
		return nil, false
	}
	var q queryTime
	if json.Unmarshal(b, &q) != nil || (maxAge > 0 && time.Since(q.At) > maxAge) {
		return nil, false
	}
	return q.Picks, true
}

func (o *openSubtitles) cachedFile(id int) ([]byte, bool) {
	o.mu.Lock()
	b, ok := o.mem[id]
	o.mu.Unlock()
	if ok {
		return b, true
	}
	if o.cfg.CacheDir != "" {
		if b, err := os.ReadFile(o.path("sub", strconv.Itoa(id))); err == nil && len(b) > 0 {
			o.remember(id, b)
			return b, true
		}
	}
	return nil, false
}

func (o *openSubtitles) remember(id int, b []byte) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if _, ok := o.mem[id]; !ok {
		o.memOrder = append(o.memOrder, id)
		for len(o.memOrder) > osMaxMemSubs {
			delete(o.mem, o.memOrder[0])
			o.memOrder = o.memOrder[1:]
		}
	}
	o.mem[id] = b
}

func (o *openSubtitles) fetch(ctx context.Context, fileID int) ([]byte, error) {
	if b, ok := o.cachedFile(fileID); ok {
		return b, nil
	}

	o.dl.Lock()
	defer o.dl.Unlock()
	if b, ok := o.cachedFile(fileID); ok {
		return b, nil
	}

	o.mu.Lock()
	until := o.blocked
	o.mu.Unlock()
	if time.Now().Before(until) {
		return nil, fmt.Errorf("%w (resumption at %s)", ErrSubtitleQuota, until.Format(time.RFC3339))
	}

	if err := o.login(ctx, false); err != nil {
		slog.Warn("opensubtitles: connection failed, downloading without account", "error", err)
	}

	var status int
	var data []byte
	var err error
	for attempt := 0; attempt < 2; attempt++ {
		status, data, err = o.request(ctx, http.MethodPost, o.base()+"/download", map[string]int{"file_id": fileID}, true)
		if err != nil {
			return nil, err
		}
		if status == http.StatusUnauthorized && attempt == 0 && o.login(ctx, true) == nil {
			continue
		}
		break
	}

	var r struct {
		Link      string `json:"link"`
		Remaining int    `json:"remaining"`
		ResetUTC  string `json:"reset_time_utc"`
		Message   string `json:"message"`
	}
	_ = json.Unmarshal(data, &r)

	switch {
	case status == http.StatusOK && r.Link != "":
	case status == http.StatusNotAcceptable || (status == http.StatusTooManyRequests && r.Remaining <= 0 && r.ResetUTC != ""):
		until, perr := time.Parse(time.RFC3339, r.ResetUTC)
		if perr != nil {
			until = time.Now().Add(time.Hour)
		}
		o.mu.Lock()
		o.blocked = until
		o.mu.Unlock()
		slog.Warn("opensubtitles: downloading quota reached", "reset", until, "message", r.Message)
		return nil, fmt.Errorf("%w (reprise à %s)", ErrSubtitleQuota, until.Format(time.RFC3339))
	default:
		return nil, fmt.Errorf("opensubtitles download: HTTP %d: %s", status, truncate(data))
	}
	if r.Remaining >= 0 && r.Remaining < 10 {
		slog.Warn("opensubtitles: soon out of downloads", "remaining", r.Remaining, "reset", r.ResetUTC)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, r.Link, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", o.cfg.UserAgent)
	resp, err := o.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opensubtitles: downloading of the file: HTTP %d", resp.StatusCode)
	}
	b, err := io.ReadAll(io.LimitReader(resp.Body, maxSidecarSize+1))
	if err != nil {
		return nil, err
	}
	if len(b) > maxSidecarSize || len(b) == 0 {
		return nil, fmt.Errorf("opensubtitles: unsuable subtitle file (%d bytes)", len(b))
	}

	o.remember(fileID, b)
	if o.cfg.CacheDir != "" {
		if err := writeAtomic(o.path("sub", strconv.Itoa(fileID)), b); err != nil {
			slog.Warn("opensubtitles: cache disque indisponible", "error", err)
		}
	}
	return b, nil
}

func opensubtitlesHash(r io.ReadSeeker, size int64) (string, error) {
	const chunk = 64 * 1024
	if size < 2*chunk {
		return "", errors.New("file too small for the hash")
	}
	hash := uint64(size)
	buf := make([]byte, chunk)
	for _, off := range []int64{0, size - chunk} {
		if _, err := r.Seek(off, io.SeekStart); err != nil {
			return "", err
		}
		if _, err := io.ReadFull(r, buf); err != nil {
			return "", err
		}
		for i := 0; i < chunk; i += 8 {
			hash += binary.LittleEndian.Uint64(buf[i : i+8])
		}
	}
	return fmt.Sprintf("%016x", hash), nil
}

func missingLanguages(want []string, probe *ProbeResult, torrentSubs []sidecarFile) []string {
	have := map[string]bool{}
	for _, si := range probe.Subs {
		if iso, _ := lookupLang(si.Language); iso != "" && !si.Forced && !strings.Contains(strings.ToLower(si.Title), "forc") {
			have[iso] = true
		}
	}
	for _, sc := range torrentSubs {
		if iso, _ := detectLangFromPath(sc.Name); iso != "" && !strings.Contains(strings.ToLower(sc.Name), "forc") {
			have[iso] = true
		}
	}
	var missing []string
	for _, l := range want {
		if !have[l] {
			missing = append(missing, l)
		}
	}
	sort.Strings(missing)
	return missing
}

func (s *StreamService) openSubtitleTracks(ctx context.Context, link string, probe *ProbeResult, torrentSubs []sidecarFile, meta MovieMeta) []sidecarFile {
	o := s.opensubs
	if o == nil {
		return nil
	}
	missing := missingLanguages(o.cfg.Languages, probe, torrentSubs)
	if len(missing) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(ctx, osLookupBudget)
	defer cancel()

	hash := ""
	hctx, hcancel := context.WithTimeout(ctx, osHashBudget)
	if h, err := s.movieHash(hctx, link); err == nil {
		hash = h
	} else {
		slog.Debug("opensubtitles: hash unavailable", "error", err)
	}
	hcancel()

	picks, err := o.find(ctx, meta, hash, missing)
	if err != nil {
		slog.Warn("opensubtitles: impossible search", "error", err)
		return nil
	}

	var out []sidecarFile
	for _, p := range picks {
		p := p
		_, lname := lookupLang(p.Lang)
		if lname == "" {
			lname = strings.ToUpper(p.Lang)
		}
		out = append(out, sidecarFile{
			Name:  fmt.Sprintf("opensubtitles/%d.srt", p.FileID),
			Label: lname + " (OpenSubtitles)",
			Lang:  p.Lang,
			Fetch: func(ctx context.Context) ([]byte, error) { return o.fetch(ctx, p.FileID) },
		})
	}
	return out
}
