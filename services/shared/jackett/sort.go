package jackett

import (
	"errors"
	"fmt"
	"math"
	"shared/tmdb"
	"time"
)

type TorrentScore struct {
	Torrent      Result
	InfoScore    int
	TitleScore   int
	QualityScore int
	Score        int
}

var ErrNoTorrentFound = errors.New("No torrent found")

func seederScore(seeders int) int {
	const cap = 30
	const maxScore = 70

	if seeders <= 0 {
		return 0
	}
	if seeders >= cap {
		return maxScore
	}
	return int(math.Log2(float64(seeders)+1) / math.Log2(float64(cap)+1) * float64(maxScore))
}

func peerScore(peers int) int {
	const cap = 10
	const maxScore = 15

	if peers <= 0 {
		return 0
	}
	if peers >= cap {
		return maxScore
	}
	return int(math.Log2(float64(peers)+1) / math.Log2(float64(cap)+1) * float64(maxScore))
}

const (
	KB = 1024
	MB = KB * 1024
	GB = MB * 1024
)

func infoScore(torrent Result) int {
	score := 0

	seederScore := seederScore(torrent.Seeders)
	if seederScore == 0 {
		return 0
	}
	peerScore := peerScore(torrent.Peers)
	if peerScore == 0 {
		return 0
	}

	score += seederScore + peerScore

	daysOld := time.Since(torrent.PublishDate.Time).Hours() / 24
	if daysOld < 7 {
		score += 5
	} else if daysOld < 30 {
		score += 2
	}

	if torrent.MagnetUri != "" {
		score += 10
	}

	const maxSize = 100 * GB
	if torrent.Size > maxSize {
		return 0
	}

	return score
}

func getScore(currentTorrent TorrentScore, bestTorrent *TorrentScore, target tmdb.Movie) {
	fmt.Printf("%s> ", currentTorrent.Torrent.Title)

	infoScore := infoScore(currentTorrent.Torrent)
	fmt.Printf("INFO: %d | ", infoScore)

	titleScore := titleScore(currentTorrent.Torrent.Title, target)
	fmt.Printf("TITLE: %d | ", titleScore)

	qualityScore := qualityScore(currentTorrent.Torrent.Title)
	fmt.Printf("QUALITY: %d\n", qualityScore)

	currentTorrent.InfoScore = infoScore
	currentTorrent.TitleScore = titleScore
	currentTorrent.QualityScore = qualityScore
	currentTorrent.Score = infoScore + titleScore + qualityScore

	if currentTorrent.Score > bestTorrent.Score {
		*bestTorrent = currentTorrent
	}
}

func FindBestResult(resp *JackettResponse, target tmdb.Movie) (TorrentScore, error) {
	if len(resp.Results) > 150 {
		resp.Results = resp.Results[:150]
	}

	if len(resp.Results) == 0 {
		return TorrentScore{}, ErrNoTorrentFound
	}

	bestTorrent := TorrentScore{
		Torrent:      resp.Results[0],
		InfoScore:    0,
		TitleScore:   0,
		QualityScore: 0,
	}

	if len(resp.Results) == 1 {
		getScore(bestTorrent, &bestTorrent, target)

		if bestTorrent.Score < 150 {
			return TorrentScore{}, ErrNoTorrentFound
		}

		return bestTorrent, nil
	}

	for _, r := range resp.Results {
		getScore(TorrentScore{
			Torrent:      r,
			InfoScore:    0,
			TitleScore:   0,
			QualityScore: 0,
			Score:        0,
		}, &bestTorrent, target)

		if bestTorrent.Score == 300 || (bestTorrent.InfoScore >= 95 && bestTorrent.TitleScore == 100 && bestTorrent.QualityScore >= 95) {
			break
		}
	}

	if bestTorrent.Score < 150 {
		return TorrentScore{}, ErrNoTorrentFound
	}

	return bestTorrent, nil
}
