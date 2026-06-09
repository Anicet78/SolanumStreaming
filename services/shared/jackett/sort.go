package jackett

import (
	"fmt"
	"regexp"
	"shared/tmdb"
	"slices"
	"strconv"
	"strings"
	"time"
)

type item struct {
	idx   int
	score int
}

func infoScore(torrent Result, target tmdb.Movie) int {
	score := 0

	fmt.Printf("Torrent ID: %s  |  TMDB ID: %s", strconv.FormatInt(torrent.IMDbId, 10), target.IMDbID)

	if torrent.IMDbId != 0 && target.IMDbID != "" {
		if strconv.FormatInt(torrent.IMDbId, 10) == target.IMDbID {
			score += 100
		} else {
			return -9999
		}
	}

	score += torrent.Seeders * 10
	score += torrent.Peers * 2

	daysOld := time.Since(torrent.PublishDate.Time).Hours() / 24

	if daysOld < 1 {
		score += 40
	} else if daysOld < 7 {
		score += 20
	} else if daysOld < 30 {
		score += 5
	} else {
		score -= 10
	}

	if torrent.MagnetUri != "" {
		score += 5
	}

	if torrent.Seeders == 0 {
		score -= 100
	}

	if strconv.Itoa(torrent.Year) == target.ReleaseDate {
		score += 10
	}

	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	const maxSize = 100 * GB
	const minSize = 500 * MB

	if torrent.Size > maxSize {
		score -= 50
	} else if torrent.Size < minSize {
		score -= 50
	}

	if torrent.Files > 2 {
		score -= 100
	}

	return score
}

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
		score += 100
		return score
	}

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

	yearStr := target.ReleaseDate
	if strings.Contains(torrentTitle, yearStr) {
		score += 10
	}

	yearRegex := regexp.MustCompile(`\b(19|20)\d{2}\b`)
	torrentYear := yearRegex.FindString(torrentTitle)

	var targetYear string
	if len(target.ReleaseDate) >= 4 {
		targetYear = target.ReleaseDate[:4]
	}

	if torrentYear != "" && targetYear != "" {
		if torrentYear == targetYear {
			score += 30
		} else {
			score -= 150
		}
	}

	return score
}

func qualityScore(torrentTitle string) int {
	title := strings.ToLower(torrentTitle)
	score := 0

	// --- Source ---
	switch {
	case strings.Contains(title, "remux"):
		score += 60
	case strings.Contains(title, "bluray") || strings.Contains(title, "blu-ray"):
		score += 50
	case strings.Contains(title, "web-dl") || strings.Contains(title, "webdl"):
		score += 35
	case strings.Contains(title, "webrip"):
		score += 25
	case strings.Contains(title, "hdtv"):
		score += 15
	case strings.Contains(title, "dvdrip") || strings.Contains(title, "dvdscr"):
		score += 5
	case strings.Contains(title, "cam") || strings.Contains(title, "ts") || strings.Contains(title, "telecine"):
		score -= 50
	}

	// --- Resolution ---
	switch {
	case strings.Contains(title, "2160p") || strings.Contains(title, "4k"):
		score += 60
	case strings.Contains(title, "1080p"):
		score += 40
	case strings.Contains(title, "720p"):
		score += 20
	case strings.Contains(title, "480p") || strings.Contains(title, "576p"):
		score -= 10
	}

	// --- Video codec ---
	switch {
	case strings.Contains(title, "av1"):
		score += 15
	case strings.Contains(title, "x265") || strings.Contains(title, "h265") || strings.Contains(title, "hevc"):
		score += 10
	case strings.Contains(title, "x264") || strings.Contains(title, "h264") || strings.Contains(title, "avc"):
		score += 5
	}

	// --- HDR ---
	if strings.Contains(title, "hdr10+") {
		score += 15
	} else if strings.Contains(title, "hdr10") || strings.Contains(title, "hdr") {
		score += 10
	}
	if strings.Contains(title, "dolby vision") || strings.Contains(title, "dv.") || strings.Contains(title, ".dv.") {
		score += 10
	}

	// --- Audio ---
	switch {
	case strings.Contains(title, "truehd") || strings.Contains(title, "atmos"):
		score += 10
	case strings.Contains(title, "dts-hd") || strings.Contains(title, "dtshd"):
		score += 8
	case strings.Contains(title, "dts"):
		score += 5
	case strings.Contains(title, "ddp") || strings.Contains(title, "eac3"):
		score += 4
	case strings.Contains(title, "ac3") || strings.Contains(title, "dd5"):
		score += 2
	case strings.Contains(title, "aac"):
		score += 1
	}

	return score
}

func getScore(torrent Result, target tmdb.Movie) int {
	score := 0

	infoScore := infoScore(torrent, target)
	if infoScore < -100 {
		return -999
	}
	score += infoScore

	titleScore := titleScore(torrent.Title, target)
	if titleScore < -50 {
		return -9999
	}
	score += titleScore

	score += qualityScore(torrent.Title)

	return score
}

func SortResults(resp *JackettResponse, target tmdb.Movie) {
	if len(resp.Results) > 25 {
		resp.Results = resp.Results[:25]
	}

	scores := make([]int, len(resp.Results))
	for i, r := range resp.Results {
		scores[i] = getScore(r, target)
	}

	tmp := make([]item, len(resp.Results))
	for i := range resp.Results {
		tmp[i] = item{i, scores[i]}
	}

	slices.SortFunc(tmp, func(a, b item) int {
		if a.score < b.score {
			return 1
		} else if a.score > b.score {
			return -1
		}
		return 0
	})

	sorted := make([]Result, len(resp.Results))
	for i, it := range tmp {
		sorted[i] = resp.Results[it.idx]
		printResult(sorted[i], i)
		fmt.Printf("       🩻  %d\n", it.score)
	}

	resp.Results = sorted
}
