package service

import (
	"context"
	"log/slog"
	"os"

	"github.com/Anicet78/SolanumStreaming/movies/internal/domain"
	"github.com/Anicet78/SolanumStreaming/movies/internal/store"
	"github.com/go-resty/resty/v2"
)

type MovieService struct {
	store *store.Queries
}

func NewMovieService(store *store.Queries) *MovieService {
	return &MovieService{store: store}
}

func (s *MovieService) Search(ctx context.Context, params domain.SearchRequestQuery) (domain.SearchResponse, error) {
	client := resty.New().
		SetBaseURL("https://api.themoviedb.org/3").
		SetHeader("Authorization", "Bearer " + os.Getenv("TMDB_TOKEN"))

	var result any
	_, err := client.R().
		SetQueryParam("query", "inception").
		SetResult(&result).
		Get("/search/movie")
	if err != nil {
		return domain.SearchResponse{}, err
	}

	slog.Error("Request result", "result", result)

	return domain.SearchResponse{}, nil
}
