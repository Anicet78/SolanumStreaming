package service

import (
	"context"
	"errors"
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

func (s *MovieService) AddToCollection(ctx context.Context, userId pgtype.UUID, movieId int) error {
	exists, err := s.store.GetFilmInCollection(ctx, store.GetFilmInCollectionParams{
		UserID:  userId,
		MovieID: int32(movieId),
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
	} else if exists.UserID.Valid {
		return domain.ErrMovieAlreadyInCollection
	}

	_, err = s.store.UserAddFilmToCollection(ctx, store.UserAddFilmToCollectionParams{
		UserID:      userId,
		MovieID:     int32(movieId),
		TorrentID:   "",
		MovieLenght: 0,
	})

	return err
}
