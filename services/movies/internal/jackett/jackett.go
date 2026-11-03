package jackett

import (
	"encoding/xml"
	"fmt"
	"strconv"
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

	if t, err := time.Parse(time.RFC3339, s); err == nil {
		ft.Time = t
		return nil
	}

	if t, err := time.Parse("2006-01-02T15:04:05", s); err == nil {
		ft.Time = t
		return nil
	}

	return fmt.Errorf("cannot parse date: %q", s)
}

type TorznabResponse struct {
	Results []Result `xml:"channel>item"`
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

type torznabItem struct {
	Title       string        `xml:"title"`
	Guid        string        `xml:"guid"`
	Link        string        `xml:"link"`
	Description string        `xml:"description"`
	PubDate     string        `xml:"pubDate"`
	Attributes  []torznabAttr `xml:"attr"`
	MagnetURI   string        `xml:"magnetURI"`
}

type torznabAttr struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

func (response *TorznabResponse) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	var feed struct {
		Channel struct {
			Items []torznabItem `xml:"item"`
		} `xml:"channel"`
	}
	if err := decoder.DecodeElement(&feed, &start); err != nil {
		return err
	}

	response.Results = make([]Result, 0, len(feed.Channel.Items))
	for _, item := range feed.Channel.Items {
		result := Result{
			Title:       item.Title,
			Guid:        item.Guid,
			Link:        item.Link,
			Details:     item.Guid,
			Description: item.Description,
			MagnetUri:   item.MagnetURI,
			PublishDate: parseTorznabDate(item.PubDate),
		}
		for _, attribute := range item.Attributes {
			setTorznabAttribute(&result, attribute)
		}
		response.Results = append(response.Results, result)
	}

	return nil
}

func parseTorznabDate(value string) FlexTime {
	for _, layout := range []string{time.RFC1123Z, time.RFC1123, time.RFC850, time.ANSIC} {
		if parsed, err := time.Parse(layout, strings.TrimSpace(value)); err == nil {
			return FlexTime{Time: parsed}
		}
	}
	return FlexTime{}
}

func setTorznabAttribute(result *Result, attribute torznabAttr) {
	value := strings.TrimSpace(attribute.Value)
	switch strings.ToLower(attribute.Name) {
	case "category":
		result.CategoryDesc = value
	case "size":
		result.Size, _ = strconv.ParseInt(value, 10, 64)
	case "seeders":
		result.Seeders, _ = strconv.Atoi(value)
	case "peers":
		result.Peers, _ = strconv.Atoi(value)
	case "grabs":
		result.Grabs, _ = strconv.ParseInt(value, 10, 64)
	case "publishdate":
		result.PublishDate = parseTorznabDate(value)
	case "imdbid":
		result.IMDbId, _ = strconv.ParseInt(strings.TrimPrefix(value, "tt"), 10, 64)
	case "tmdbid":
		result.TMDbId, _ = strconv.ParseInt(value, 10, 64)
	case "infohash":
		result.InfoHash = value
	case "downloadvolumefactor":
		result.DownloadVolumeFactor, _ = strconv.ParseFloat(value, 64)
	case "uploadvolumefactor":
		result.UploadVolumeFactor, _ = strconv.ParseFloat(value, 64)
	case "magneturl", "magneturi":
		result.MagnetUri = value
	}
}

func New() *resty.Request {
	client := resty.New().
		SetBaseURL("http://jackett:9117/api/v2.0")
	return client.R()
}
