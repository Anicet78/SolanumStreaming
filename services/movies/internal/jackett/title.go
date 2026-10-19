package jackett

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"regexp"
	"slices"
	"strings"

	"github.com/Anicet78/SolanumStreaming/movies/internal/tmdb"
)

/* func approxNameScore(targetNorm string, targetOrigNorm string, torrentNorm string) int {
	score := 0

	targetWords := strings.Fields(targetNorm)
	torrentWords := strings.Fields(torrentNorm)

	lengthRatio := float64(len(torrentWords)) / float64(max(len(targetWords), 1))
	if lengthRatio > 2.0 {
		score -= 200
	} else if lengthRatio > 1.4 {
		score -= 50
	}

	prefixMatch := true
	for i, w := range targetWords {
		if i >= len(torrentWords) || torrentWords[i] != w {
			prefixMatch = false
			break
		}
	}
	if !prefixMatch {
		score -= 80
	}

	targetSet := toSet(strings.Fields(targetNorm))
	torrentSet := toSet(strings.Fields(torrentNorm))

	intersection := 0
	for w := range targetSet {
		if torrentSet[w] {
			intersection++
		}
	}
	union := len(targetSet) + len(torrentSet) - intersection
	jaccard := float64(intersection) / float64(max(union, 1))
	score += int(jaccard * 60)

	dist := levenshtein(torrentNorm, targetNorm)
	maxLen := max(len(torrentNorm), len(targetNorm))
	similarity := 1.0 - float64(dist)/float64(max(maxLen, 1))
	if similarity > 0.85 {
		score += 40
	} else if similarity > 0.65 {
		score += 20
	} else if similarity < 0.3 {
		score -= 30
	}

	if strings.HasPrefix(torrentNorm, targetNorm) || strings.HasPrefix(torrentNorm, targetOrigNorm) {
		score += 20
	}

	return score
} */

func titleScore(torrentTitle string, target tmdb.Movie) int {
	score := 0

	cleanTorrent := func(s string) string {
		re := regexp.MustCompile(`(?i)\b(1080p?|720p?|2160p?|4k|bluray|blu-ray|webrip|web-dl|webdl|hdtv|dvdrip|dvdscr|bdrip|remux|x264|x265|h264|h265|hevc|avc|aac|dts|ac3|mp3|atmos|truehd|hdr|hdr10|dolby|french|english|vostfr|multi|proper|repack|extended|theatrical|directors\.cut|\d{4})\b.*`)
		return strings.TrimSpace(re.ReplaceAllString(s, ""))
	}

	targetNorm := normalize(target.Title)
	targetOrigNorm := normalize(target.OriginalTitle)
	torrentClean := cleanTorrent(torrentTitle)
	torrentNorm := normalize(torrentClean)

	if torrentNorm == targetNorm || torrentNorm == targetOrigNorm {
		score += 70
	} else {
		// score += approxNameScore(targetNorm, targetOrigNorm, torrentNorm)
		return 0
	}

	yearRegex := regexp.MustCompile(`\b(19|20)\d{2}\b`)
	torrentYear := yearRegex.FindString(torrentTitle)

	var targetYear string
	if len(target.ReleaseDate) >= 4 {
		targetYear = target.ReleaseDate[:4]
	}

	if targetYear == "" || (torrentYear != "" && torrentYear == targetYear) {
		return score + 30
	} else {
		return 0
	}
}

type Tier string

var tierRank = map[Tier]int{"sd": 0, "hd": 1, "fhd": 2, "qhd": 3, "uhd": 4, "dci": 5}

type Display struct {
	Tier   Tier `json:"tier"`
	Width  int  `json:"width"`
	Height int  `json:"height"`
	HDR    bool `json:"hdr"`
}

type CodecSupport struct {
	MaxTier       *Tier `json:"maxTier"`
	SmoothMaxTier *Tier `json:"smoothMaxTier"`
	HW            bool  `json:"hw"`
}

type ClientProfile struct {
	Display Display                 `json:"display"`
	Video   map[string]CodecSupport `json:"video"`
	Audio   []string                `json:"audio"`
}

func ParseClientProfile(param string) (*ClientProfile, error) {
	raw, err := base64.StdEncoding.DecodeString(param)
	if err != nil {
		return nil, err
	}
	var p ClientProfile
	if err := json.Unmarshal(raw, &p); err != nil {
		return nil, err
	}
	if _, ok := tierRank[p.Display.Tier]; !ok {
		return nil, errors.New("invalid display tier")
	}
	return &p, nil
}

func defaultProfile() *ClientProfile {
	t := Tier("fhd")
	return &ClientProfile{
		Display: Display{Tier: "fhd", Width: 1920, Height: 1080},
		Video:   map[string]CodecSupport{"h264": {MaxTier: &t, SmoothMaxTier: &t, HW: true}},
		Audio:   []string{"aac"},
	}
}

var (
	// Source
	reRemux  = regexp.MustCompile(`\bremux\b`)
	reBluRay = regexp.MustCompile(`\bblu-?ray\b`)
	reWebDL  = regexp.MustCompile(`\bweb-?dl\b`)
	reWebRip = regexp.MustCompile(`\bwebrip\b`)
	reHDTV   = regexp.MustCompile(`\bhdtv\b`)
	reDVDRip = regexp.MustCompile(`\b(dvdrip|dvdscr)\b`)
	reCam    = regexp.MustCompile(`\b(cam|ts|telecine|telesync)\b`)

	// Resolution
	reDS4K   = regexp.MustCompile(`\bds4k\b`)
	re4K     = regexp.MustCompile(`\b(4k|2160p|uhd)\b`)
	re1440p  = regexp.MustCompile(`\b(2k|1440p|qhd)\b`)
	re1080p  = regexp.MustCompile(`\b1080p\b`)
	reM1080p = regexp.MustCompile(`\bm1080p\b`)
	re720p   = regexp.MustCompile(`\b720p\b`)
	reM720p  = regexp.MustCompile(`\bm720p\b`)
	re480p   = regexp.MustCompile(`\b(480p|576p)\b`)

	// Video codec
	reAV1  = regexp.MustCompile(`\bav1\b`)
	reX264 = regexp.MustCompile(`\b(x264|h\.?264|avc)\b`)
	reX265 = regexp.MustCompile(`\b(x265|h\.?265|hevc)\b`)

	// HDR
	reHDR         = regexp.MustCompile(`\bhdr(10)?(\+|plus)?`)
	reDolbyVision = regexp.MustCompile(`\b(dolby[.\s]?vision|dv)\b`)

	// Audio
	reTrueHD = regexp.MustCompile(`\b(truehd|atmos)\b`)
	reDTSHD  = regexp.MustCompile(`\bdts-?hd\b`)
	reDTS    = regexp.MustCompile(`\bdts\b`)
	reEAC3   = regexp.MustCompile(`\b(ddp|eac3|dd\+)`)
	reAC3    = regexp.MustCompile(`\b(ac3|dd5\.?1)\b`)
	reAAC    = regexp.MustCompile(`\baac\b`)
)

type parsedTitle struct {
	cam     bool
	source  string
	tier    Tier
	reduced bool
	codec   string
	hdr     bool
	dv      bool
	audio   string
}

func parseTitle(raw string) parsedTitle {
	title := strings.ToLower(raw)
	var p parsedTitle

	switch {
	case reRemux.MatchString(title):
		p.source = "remux"
	case reBluRay.MatchString(title):
		p.source = "bluray"
	case reWebDL.MatchString(title):
		p.source = "webdl"
	case reWebRip.MatchString(title):
		p.source = "webrip"
	case reHDTV.MatchString(title):
		p.source = "hdtv"
	case reDVDRip.MatchString(title):
		p.source = "dvdrip"
	case reCam.MatchString(title):
		p.cam = true
	}

	switch {
	case reDS4K.MatchString(title):
		p.tier = "fhd"
		if re720p.MatchString(title) {
			p.tier = "hd"
		}
	case re4K.MatchString(title):
		p.tier = "uhd"
	case re1440p.MatchString(title):
		p.tier = "qhd"
	case re1080p.MatchString(title):
		p.tier = "fhd"
	case reM1080p.MatchString(title):
		p.tier, p.reduced = "fhd", true
	case re720p.MatchString(title):
		p.tier = "hd"
	case reM720p.MatchString(title):
		p.tier, p.reduced = "hd", true
	case re480p.MatchString(title):
		p.tier = "sd"
	}

	switch {
	case reAV1.MatchString(title):
		p.codec = "av1"
	case reX265.MatchString(title):
		p.codec = "hevc"
	case reX264.MatchString(title):
		p.codec = "h264"
	}

	p.hdr = reHDR.MatchString(title)
	p.dv = reDolbyVision.MatchString(title)

	switch {
	case reTrueHD.MatchString(title):
		p.audio = "truehd"
	case reDTSHD.MatchString(title):
		p.audio = "dtshd"
	case reDTS.MatchString(title):
		p.audio = "dts"
	case reEAC3.MatchString(title):
		p.audio = "eac3"
	case reAC3.MatchString(title):
		p.audio = "ac3"
	case reAAC.MatchString(title):
		p.audio = "aac"
	}

	return p
}

func (p parsedTitle) fileTier() Tier {
	if p.tier == "" {
		return "fhd"
	}
	return p.tier
}

func (p parsedTitle) effectiveCodec() string {
	if p.codec != "" {
		return p.codec
	}
	if tierRank[p.fileTier()] >= tierRank["uhd"] {
		return "hevc"
	}
	return "h264"
}

func qualityScore(torrentTitle string, client *ClientProfile) int {
	if client == nil {
		client = defaultProfile()
	}
	p := parseTitle(torrentTitle)

	if p.cam || !playable(p, client) {
		return 0
	}

	score := resolutionScore(p, client) +
		sourceScore(p.source) +
		codecScore(p, client) +
		hdrScore(p, client) +
		audioScore(p.audio, client)

	return score
}

func playable(p parsedTitle, c *ClientProfile) bool {
	support, ok := c.Video[p.effectiveCodec()]
	if !ok || support.SmoothMaxTier == nil {
		return false
	}
	return tierRank[p.fileTier()] <= tierRank[*support.SmoothMaxTier]
}

func resolutionScore(p parsedTitle, c *ClientProfile) int {
	if p.tier == "" {
		return 10
	}
	diff := tierRank[p.tier] - tierRank[c.Display.Tier]

	var s int
	switch {
	case diff == 0:
		s = 45
	case diff == 1:
		s = 38
	case diff == 2:
		s = 30
	case diff >= 3:
		s = 29
	case diff == -1:
		s = 25
	case diff == -2:
		s = 15
	case diff == -3:
		s = 5
	default:
		s = 3
	}
	if p.reduced {
		s -= 5
	}
	return s
}

func sourceScore(source string) int {
	switch source {
	case "remux":
		return 20
	case "bluray", "webdl":
		return 16
	case "webrip":
		return 8
	case "hdtv":
		return 4
	case "dvdrip":
		return 2
	default:
		return 6
	}
}

func codecScore(p parsedTitle, c *ClientProfile) int {
	codec := p.effectiveCodec()
	base := map[string]int{"av1": 10, "hevc": 10, "h264": 6}[codec]
	if c.Video[codec].HW {
		base += 5
	}
	return base
}

func hdrScore(p parsedTitle, c *ClientProfile) int {
	if c.Display.HDR {
		switch {
		case p.hdr:
			return 10
		case p.dv:
			return 6
		default:
			return 5
		}
	}
	if p.hdr || p.dv {
		return 4
	}
	return 10
}

func audioScore(audio string, c *ClientProfile) int {
	switch audio {
	case "":
		return 5
	case "aac", "ac3", "eac3":
		if slices.Contains(c.Audio, audio) {
			return 10
		}
		return 6
	default:
		return 7
	}
}
