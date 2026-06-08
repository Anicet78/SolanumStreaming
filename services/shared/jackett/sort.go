package jackett

import (
	"fmt"
	"regexp"
	"shared/tmdb"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type item struct {
	idx   int
	score int
}

func removeAccents(s string) string {
	t := transform.Chain(norm.NFD, transform.RemoveFunc(func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn = non-spacing marks
	}), norm.NFC)
	result, _, _ := transform.String(t, s)
	return result
}

func titleCoherenceScore(torrentTitle string, target tmdb.Movie) int {
	score := 0

	normalize := func(s string) string {
		s = removeAccents(s)
		s = strings.ToLower(s)
		// Retire articles et ponctuation
		re := regexp.MustCompile(`[^\w\s]`)
		s = re.ReplaceAllString(s, " ")
		s = strings.TrimSpace(s)
		// Retire les articles en début
		for _, article := range []string{"the ", "a ", "an ", "le ", "la ", "les ", "un ", "une "} {
			if strings.HasPrefix(s, article) {
				s = strings.TrimPrefix(s, article)
			}
		}
		return s
	}

	// Retire les tags techniques du titre torrent (qualité, codec, langue...)
	cleanTorrent := func(s string) string {
		// Patterns typiques : 1080p, BluRay, x264, HEVC, HDR, FRENCH, VOSTFR, etc.
		re := regexp.MustCompile(`(?i)\b(1080p?|720p?|2160p?|4k|bluray|blu-ray|webrip|web-dl|webdl|hdtv|dvdrip|dvdscr|bdrip|remux|x264|x265|h264|h265|hevc|avc|aac|dts|ac3|mp3|atmos|truehd|hdr|hdr10|dolby|french|english|vostfr|multi|proper|repack|extended|theatrical|directors\.cut|\d{4})\b.*`)
		return strings.TrimSpace(re.ReplaceAllString(s, ""))
	}

	targetNorm := normalize(target.Title)

	// Utilise aussi le titre original si différent
	targetOrigNorm := normalize(target.Title)

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

	return score
}

func toSet(words []string) map[string]bool {
	set := make(map[string]bool)
	for _, w := range words {
		if len(w) > 2 {
			set[w] = true
		}
	}
	return set
}

func levenshtein(a, b string) int {
	ra, rb := []rune(a), []rune(b)
	la, lb := len(ra), len(rb)
	dp := make([][]int, la+1)
	for i := range dp {
		dp[i] = make([]int, lb+1)
		dp[i][0] = i
	}
	for j := 0; j <= lb; j++ {
		dp[0][j] = j
	}
	for i := 1; i <= la; i++ {
		for j := 1; j <= lb; j++ {
			if ra[i-1] == rb[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], min(dp[i][j-1], dp[i-1][j-1]))
			}
		}
	}
	return dp[la][lb]
}

func getScore(torrent Result, target tmdb.Movie) int {
	score := 0

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

	titleScore := titleCoherenceScore(torrent.Title, target)

	// Seuil minimal : si le titre ne matche pas du tout, on rejette directement
	if titleScore < -50 {
		return -9999 // éliminé
	}

	score += titleScore

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

func SortResults(resp *JackettResponse, target tmdb.Movie) {
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
