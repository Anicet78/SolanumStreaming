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

func New() *resty.Request {
	client := resty.New().
		SetBaseURL("http://jackett:9117/api/v2.0")
	return client.R()
}
