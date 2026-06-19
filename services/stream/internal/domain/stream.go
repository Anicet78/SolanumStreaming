package domain

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type MovieLinkParam struct {
	MovieLink int `param:"movie_link"`
}

type MovieProgressionQuery struct {
	Progression pgtype.Interval `json:"progression"`
}
