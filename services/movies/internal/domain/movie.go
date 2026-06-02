package domain

import "errors"

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

type PostCollectionRequestBody struct {
	MovieID int `json:"movie_id"`
}

var ErrMovieAlreadyInCollection = errors.New("Movie already in collection")
var ErrTMDBRequestFailed = errors.New("TMDB request failed")
