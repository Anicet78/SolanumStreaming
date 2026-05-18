package service

import (
	"context"
	"shared/tmdb"
	"strconv"

	"github.com/Anicet78/SolanumStreaming/movies/internal/domain"
	"github.com/Anicet78/SolanumStreaming/movies/internal/store"
)

type MovieService struct {
	store *store.Queries
}

func NewMovieService(store *store.Queries) *MovieService {
	return &MovieService{store: store}
}

func (s *MovieService) Search(ctx context.Context, params domain.SearchRequestQuery) (domain.SearchResponse, error) {
	var result domain.SearchResponse
	res, err := tmdb.New().
		SetQueryParam("query", params.Title).
		SetQueryParam("include_adult", "true").
		SetQueryParam("page", strconv.Itoa(params.Page)).
		SetResult(&result).
		Get("/search/movie")
	if err != nil || res.IsSuccess() == false {
		return domain.SearchResponse{}, domain.ErrTMDBRequestFailed
	}

	return result, nil
}
