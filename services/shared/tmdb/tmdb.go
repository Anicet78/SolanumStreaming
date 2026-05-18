package tmdb

import (
	"os"

	"github.com/go-resty/resty/v2"
)

func New() *resty.Request {
	client := resty.New().
		SetBaseURL("https://api.themoviedb.org/3").
		SetHeader("Authorization", "Bearer "+os.Getenv("TMDB_TOKEN"))
	return client.R()
}
