package jackett

import (
	"regexp"
	"shared/tmdb"
	"strings"
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
	re1440p  = regexp.MustCompile(`\b1440p\b`)
	re1080p  = regexp.MustCompile(`\b1080p\b`)
	reM1080p = regexp.MustCompile(`\bm1080p\b`)
	re720p   = regexp.MustCompile(`\b720p\b`)
	reM720p  = regexp.MustCompile(`\bm720p\b`)
	re480p   = regexp.MustCompile(`\b(480p|576p)\b`)

	// Video codec
	reAV1  = regexp.MustCompile(`\bav1\b`)
	reX265 = regexp.MustCompile(`\b(x265|h265|hevc)\b`)
	reX264 = regexp.MustCompile(`\b(x264|h264|avc)\b`)

	// HDR
	reHDR10Plus   = regexp.MustCompile(`\bhdr10\+\b`)
	reHDR         = regexp.MustCompile(`\bhdr(10)?\b`)
	reDolbyVision = regexp.MustCompile(`\b(dolby[.\s]?vision|dv)\b`)

	// Audio
	reAtmosTrueHD = regexp.MustCompile(`\b(truehd|atmos)\b`)
	reDTSHD       = regexp.MustCompile(`\bdts-?hd\b`)
	reDTS         = regexp.MustCompile(`\bdts\b`)
	reEAC3        = regexp.MustCompile(`\b(ddp|eac3|dd\+)\b`)
	reAC3         = regexp.MustCompile(`\b(ac3|dd5\.?1)\b`)
	reAAC         = regexp.MustCompile(`\baac\b`)
)

func qualityScore(torrentTitle string) int {
	title := strings.ToLower(torrentTitle)
	score := 0

	// --- Source ---
	switch {
	case reRemux.MatchString(title):
		score += 10
	case reBluRay.MatchString(title):
		score += 8
	case reWebDL.MatchString(title):
		score += 8
	case reWebRip.MatchString(title):
		score += 3
	case reHDTV.MatchString(title):
		score += 1
	case reCam.MatchString(title):
		return 0
	}

	// --- Resolution ---
	switch {
	case reDS4K.MatchString(title):
		switch {
		case re1080p.MatchString(title):
			score += 65
		case re720p.MatchString(title):
			score += 40
		default:
			score += 50
		}
	case re4K.MatchString(title):
		score += 75
	case re1440p.MatchString(title):
		score += 70
	case re1080p.MatchString(title):
		score += 60
	case reM1080p.MatchString(title):
		score += 50
	case re720p.MatchString(title):
		score += 35
	case reM720p.MatchString(title):
		score += 25
	case re480p.MatchString(title):
		score += 10
	}

	// --- Video codec ---
	switch {
	case reAV1.MatchString(title):
		score += 3
	case reX265.MatchString(title):
		score += 2
	case reX264.MatchString(title):
		score += 1
	}

	// --- HDR ---
	switch {
	case reHDR10Plus.MatchString(title):
		score += 4
	case reHDR.MatchString(title):
		score += 2
	}
	if reDolbyVision.MatchString(title) {
		score += 5
	}

	// --- Audio ---
	switch {
	case reAtmosTrueHD.MatchString(title):
		score += 3
	case reDTSHD.MatchString(title):
		score += 3
	case reDTS.MatchString(title):
		score += 2
	case reEAC3.MatchString(title):
		score += 2
	case reAC3.MatchString(title):
		score += 1
	case reAAC.MatchString(title):
		score += 1
	}

	return score
}
