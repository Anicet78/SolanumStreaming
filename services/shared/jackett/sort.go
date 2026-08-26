package jackett

import (
	"fmt"
	"math"
	"shared/tmdb"
	"slices"
	"time"
)

type item struct {
	idx   int
	score int
}

func seederScore(seeders int) int {
	const cap = 30

	if seeders <= 0 {
		return 0
	}
	if seeders >= cap {
		seeders = cap
	}
	return int(math.Log2(float64(seeders)+1) / math.Log2(float64(cap)+1) * 100)
}

func peerScore(peers int) int {
	const cap = 10
	const maxScore = 8

	if peers <= 0 {
		return 0
	}
	if peers >= cap {
		return maxScore
	}
	return int(math.Log2(float64(peers)+1) / math.Log2(float64(cap)+1) * float64(maxScore))
}

func infoScore(torrent Result) int {
	score := 0

	score += seederScore(torrent.Seeders)
	score += peerScore(torrent.Peers)

	daysOld := time.Since(torrent.PublishDate.Time).Hours() / 24

	if daysOld < 7 {
		score += 8
	} else if daysOld < 30 {
		score += 2
	}

	if torrent.MagnetUri != "" {
		score += 10
	}

	if torrent.Seeders == 0 {
		score -= 300
	}

	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
	)

	const maxSize = 100 * GB

	if torrent.Size > maxSize {
		score -= 100
	}

	return score
}

func getScore(torrent Result, target tmdb.Movie) int {
	score := 0

	fmt.Printf("%s> ", torrent.Title)

	infoScore := infoScore(torrent)
	/* if infoScore < -100 {
		return -999
	} */
	score += infoScore
	fmt.Printf("INFO: %d | ", infoScore)

	titleScore := titleScore(torrent.Title, target)
	/* if titleScore < -50 {
		return -9999
	} */
	score += titleScore
	fmt.Printf("TITLE: %d | ", titleScore)

	qualityScore := qualityScore(torrent.Title)
	score += qualityScore
	fmt.Printf("QUALITY: %d\n", qualityScore)

	return score
}

func SortResults(resp *JackettResponse, target tmdb.Movie) {
	if len(resp.Results) > 150 {
		resp.Results = resp.Results[:150]
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

	for i, it := range tmp {
		fmt.Printf("%d> [%s] : %d\n", i, resp.Results[it.idx].Title, it.score)
	}

	sorted := make([]Result, len(resp.Results))
	for i, it := range tmp {
		sorted[i] = resp.Results[it.idx]
	}

	resp.Results = sorted
}
