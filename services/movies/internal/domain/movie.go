package domain

import (
	"errors"

	"github.com/jackc/pgx/v5/pgtype"
)

type SearchRequestQuery struct {
	Title string `query:"title" validate:"required"`
	Page  int    `query:"page" validate:"required,min=1"`
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

type MovieIDParam struct {
	MovieID int `param:"movie_id"`
}

var ErrMovieNotInCollection = errors.New("Movie not in collection")
var ErrMovieAlreadyInCollection = errors.New("Movie already in collection")
var ErrTMDBRequestFailed = errors.New("TMDB request failed")
var ErrJackettRequestFailed = errors.New("Jackett request failed")
