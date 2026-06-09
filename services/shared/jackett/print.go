package jackett

import (
	"fmt"
	"strings"
	"time"
)

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

func printResult(r Result, index int) {
	cyan := func(s string) string { return "\033[36m" + s + "\033[0m" }
	bold := func(s string) string { return "\033[1m" + s + "\033[0m" }
	dim := func(s string) string { return "\033[2m" + s + "\033[0m" }

	// Title + tracker
	fmt.Printf("\n  %s  %s\n",
		bold(fmt.Sprintf("[%d]", index+1)),
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

	// External IDs
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
}

func PrintResponse(resp JackettResponse) {
	sep := strings.Repeat("─", 72)
	bold := func(s string) string { return "\033[1m" + s + "\033[0m" }
	dim := func(s string) string { return "\033[2m" + s + "\033[0m" }

	// Header
	fmt.Printf("\n%s\n", sep)
	fmt.Printf("  %s   %s\n",
		bold(fmt.Sprintf("%d results", len(resp.Results))),
		dim(fmt.Sprintf("across %d indexers", len(resp.Indexers))),
	)
	fmt.Printf("%s\n\n", sep)

	// Indexer summary
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

	// Results
	for i, r := range resp.Results {
		printResult(r, i)
		fmt.Printf("       %s\n", dim(sep))
	}
}
