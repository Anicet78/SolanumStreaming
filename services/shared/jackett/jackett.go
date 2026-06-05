package jackett

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

type FlexTime struct {
	time.Time
}

func (ft *FlexTime) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" || s == "0001-01-01T00:00:00" {
		ft.Time = time.Time{}
		return nil
	}

	// Avec timezone
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		ft.Time = t
		return nil
	}

	// Sans timezone
	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		ft.Time = t
		return nil
	}

	return fmt.Errorf("cannot parse date: %q", s)
}

type JackettResponse struct {
	Results  []Result  `json:"Results"`
	Indexers []Indexer `json:"Indexers"`
}

type Result struct {
	// Identification
	Title       string `json:"Title"`
	Tracker     string `json:"Tracker"`
	TrackerId   string `json:"TrackerId"`
	TrackerType string `json:"TrackerType"`

	// Categories
	CategoryDesc string `json:"CategoryDesc"`
	Category     []int  `json:"Category"`

	// Links
	Guid      string `json:"Guid"`
	Link      string `json:"Link"`
	Details   string `json:"Details"`
	MagnetUri string `json:"MagnetUri"`
	InfoHash  string `json:"InfoHash"`
	Poster    string `json:"Poster"`
	IMDbId    int64  `json:"Imdb"`
	TMDbId    int64  `json:"TMDb"`
	TraktId   int64  `json:"Trakt"`
	TVDBId    int64  `json:"TVDBId"`
	TVMazeId  int64  `json:"TVMazeId"`
	RageId    int64  `json:"RageID"`
	DoubanId  int64  `json:"DoubanId"`

	// Metadata
	PublishDate FlexTime `json:"PublishDate"`
	Description string   `json:"Description"`
	Author      string   `json:"Author"`
	BookTitle   string   `json:"BookTitle"`
	Artist      string   `json:"Artist"`
	Album       string   `json:"Album"`
	Label       string   `json:"Label"`
	Track       string   `json:"Track"`
	Year        int      `json:"Year"`
	Genre       string   `json:"Genre"`
	Languages   []string `json:"Languages"`
	Subs        []string `json:"Subs"`

	// Size / Ratio
	Size                 int64   `json:"Size"`
	Files                int64   `json:"Files"`
	Grabs                int64   `json:"Grabs"`
	Seeders              int     `json:"Seeders"`
	Peers                int     `json:"Peers"`
	DownloadVolumeFactor float64 `json:"DownloadVolumeFactor"`
	UploadVolumeFactor   float64 `json:"UploadVolumeFactor"`

	// Flags
	BlackholeLink string `json:"BlackholeLink"`

	// Jackett internal timestamp
	FirstSeen FlexTime `json:"FirstSeen"`
}

type Indexer struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Status  int    `json:"Status"`
	Results int    `json:"Results"`
	Error   string `json:"Error"`
}

func printSize(bytes int64) string {
	switch {
	case bytes >= 1<<30:
		return fmt.Sprintf("%.2f GB", float64(bytes)/float64(1<<30))
	case bytes >= 1<<20:
		return fmt.Sprintf("%.2f MB", float64(bytes)/float64(1<<20))
	case bytes >= 1<<10:
		return fmt.Sprintf("%.2f KB", float64(bytes)/float64(1<<10))
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func printAge(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	d := time.Since(t)
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}

func printOptionalIDs(r Result) string {
	var ids []string
	if r.IMDbId != 0 {
		ids = append(ids, fmt.Sprintf("IMDb:%d", r.IMDbId))
	}
	if r.TMDbId != 0 {
		ids = append(ids, fmt.Sprintf("TMDb:%d", r.TMDbId))
	}
	if r.TVDBId != 0 {
		ids = append(ids, fmt.Sprintf("TVDB:%d", r.TVDBId))
	}
	if r.TraktId != 0 {
		ids = append(ids, fmt.Sprintf("Trakt:%d", r.TraktId))
	}
	if len(ids) == 0 {
		return ""
	}
	return strings.Join(ids, " · ")
}

func PrintResponse(resp JackettResponse) {
	sep := strings.Repeat("─", 72)
	bold := func(s string) string { return "\033[1m" + s + "\033[0m" }
	dim := func(s string) string { return "\033[2m" + s + "\033[0m" }
	cyan := func(s string) string { return "\033[36m" + s + "\033[0m" }

	// ── Header ──────────────────────────────────────────────────────────────
	fmt.Printf("\n%s\n", sep)
	fmt.Printf("  %s   %s\n",
		bold(fmt.Sprintf("%d results", len(resp.Results))),
		dim(fmt.Sprintf("across %d indexers", len(resp.Indexers))),
	)
	fmt.Printf("%s\n\n", sep)

	// ── Indexer summary ─────────────────────────────────────────────────────
	for _, idx := range resp.Indexers {
		status := "✓"
		if idx.Status != 2 {
			status = "✗"
		}
		line := fmt.Sprintf("  %s %-30s %s",
			status,
			idx.Name,
			dim(fmt.Sprintf("%d results", idx.Results)),
		)
		if idx.Error != "" {
			line += "  " + fmt.Sprintf("\033[31m%s\033[0m", idx.Error)
		}
		fmt.Println(line)
	}

	fmt.Printf("\n%s\n", sep)

	// ── Results ─────────────────────────────────────────────────────────────
	for i, r := range resp.Results {
		// Title + tracker
		fmt.Printf("\n  %s  %s\n",
			bold(fmt.Sprintf("[%d]", i+1)),
			r.Title,
		)
		fmt.Printf("       %s  %s\n",
			cyan(r.Tracker),
			dim(r.CategoryDesc),
		)

		// Core stats
		fmt.Printf("       📦 %-12s  🌱 %-6d  👥 %-6d  🕐 %s\n",
			printSize(r.Size),
			r.Seeders,
			r.Peers,
			printAge(r.PublishDate.Time),
		)

		// External IDs (optional)
		if ids := printOptionalIDs(r); ids != "" {
			fmt.Printf("       %s\n", dim(ids))
		}

		// Links
		if r.MagnetUri != "" {
			short := r.MagnetUri
			if len(short) > 60 {
				short = short[:60] + "…"
			}
			fmt.Printf("       🧲 %s\n", dim(short))
		} else if r.Link != "" {
			fmt.Printf("       🔗 %s\n", dim(r.Link))
		}

		fmt.Printf("       %s\n", dim(sep))
	}
}

func New() *resty.Request {
	client := resty.New().
		SetBaseURL("http://jackett:9117/api/v2.0")
	return client.R()
}
