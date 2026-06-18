package service

import (
	"context"
	"errors"
	"fmt"
	"shared/jackett"
	"shared/tmdb"
	"strconv"

	"github.com/Anicet78/SolanumStreaming/movies/internal/domain"
	"github.com/Anicet78/SolanumStreaming/movies/internal/store"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type MovieService struct {
	store *store.Queries
}

func NewMovieService(store *store.Queries) *MovieService {
	return &MovieService{store: store}
}

func (s *MovieService) Stream(ctx context.Context) (tmdb.SearchResponse, error) {
	var result tmdb.SearchResponse
	res, err := tmdb.New().
		SetQueryParam("query", params.Title).
		SetQueryParam("include_adult", "true").
		SetQueryParam("page", strconv.Itoa(params.Page)).
		SetResult(&result).
		Get("/search/movie")
	if err != nil {
		return tmdb.SearchResponse{}, err
	} else if res.IsSuccess() == false {
		return tmdb.SearchResponse{}, domain.ErrTMDBRequestFailed
	}

	return result, nil
}
