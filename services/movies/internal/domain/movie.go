package domain

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type SearchRequestQuery struct {
	Title string `query:"title" validate:"required"`
	Page  int    `query:"page" validate:"required,min=1"`
}

type SearchResponse struct {
	Page         int     `json:"page"`
	Results      []Movie `json:"results"`
	TotalPages   int     `json:"total_pages"`
	TotalResults int     `json:"total_results"`
}

type Movie struct {
	ID          int     `json:"id"`
	Title       string  `json:"title"`
	Popularity  float64 `json:"popularity"`
	Adult       bool    `json:"adult"`
	Video       bool    `json:"video"`
	PosterPath  string  `json:"poster_path"`
	ReleaseDate string  `json:"release_date"`
}

type CollectionRequestBody struct {
	MovieID int `json:"movie_id"`
}

type CollectionMovie struct {
	MovieID     int             `json:"movie_id"`
	TorrentID   string          `json:"torrent_id"`
	Length      int             `json:"length"`
	Progression pgtype.Interval `json:"progression"`
}

var ErrMovieAlreadyInCollection = errors.New("Movie already in collection")
var ErrTMDBRequestFailed = errors.New("TMDB request failed")
