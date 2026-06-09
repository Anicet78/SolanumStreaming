package tmdb

import (
	"os"

	"github.com/go-resty/resty/v2"
)

type SearchResponse struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

type Movie struct {
	ID            int     `json:"id"`
	Title         string  `json:"title"`
	OriginalTitle string  `json:"original_title"`
	Popularity    float64 `json:"popularity"`
	Adult         bool    `json:"adult"`
	Video         bool    `json:"video"`
	PosterPath    string  `json:"poster_path"`
	ReleaseDate   string  `json:"release_date"`
	IMDbID        string  `json:"imdb_id"`
	Runtime       int     `json:"runtime"`
}

func New() *resty.Request {
	client := resty.New().
		SetBaseURL("https://api.themoviedb.org/3").
		SetHeader("Authorization", "Bearer "+os.Getenv("TMDB_TOKEN"))
	return client.R()
}
