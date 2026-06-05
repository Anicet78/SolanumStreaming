package service

import (
	"context"
	"errors"
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

func (s *MovieService) Search(ctx context.Context, params domain.SearchRequestQuery) (domain.SearchResponse, error) {
	var result domain.SearchResponse
	res, err := tmdb.New().
		SetQueryParam("query", params.Title).
		SetQueryParam("include_adult", "true").
		SetQueryParam("page", strconv.Itoa(params.Page)).
		SetResult(&result).
		Get("/search/movie")
	if err != nil {
		return domain.SearchResponse{}, err
	} else if res.IsSuccess() == false {
		return domain.SearchResponse{}, domain.ErrTMDBRequestFailed
	}

	return result, nil
}

func (s *MovieService) GetCollection(ctx context.Context, userId pgtype.UUID) ([]domain.CollectionMovie, error) {
	movies, err := s.store.GetAllMoviesInCollection(ctx, userId)
	if err != nil {
		return nil, err
	}

	res := make([]domain.CollectionMovie, 0, len(movies))

	for _, u := range movies {
		res = append(res, domain.CollectionMovie{
			MovieID:     int(u.MovieID),
			TorrentID:   u.TorrentID,
			Length:      int(u.Length),
			Progression: u.Progression,
		})
	}

	return res, nil
}

func (s *MovieService) GetInCollection(ctx context.Context, userId pgtype.UUID, movieId int) (domain.CollectionMovie, error) {
	movie, err := s.store.GetMovieInCollection(ctx, store.GetMovieInCollectionParams{
		UserID:  userId,
		MovieID: int32(movieId),
	})
	if err != nil {
		return domain.CollectionMovie{}, domain.ErrMovieNotInCollection
	}

	return domain.CollectionMovie{
		MovieID:     int(movie.MovieID),
		TorrentID:   movie.TorrentID,
		Length:      int(movie.Length),
		Progression: movie.Progression,
	}, nil
}

func (s *MovieService) AddToCollection(ctx context.Context, userId pgtype.UUID, movieId int) error {
	exists, err := s.store.GetMovieInCollection(ctx, store.GetMovieInCollectionParams{
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

	var TMDBResult domain.Movie
	res, err := tmdb.New().
		SetResult(&TMDBResult).
		SetPathParam("movie_id", strconv.Itoa(movieId)).
		Get("/movie/{movie_id}")
	if err != nil {
		return err
	} else if res.IsSuccess() == false {
		return domain.ErrTMDBRequestFailed
	}

	var JackettResult jackett.JackettResponse
	res, err = jackett.New().
		SetQueryParam("query", TMDBResult.Title).
		SetQueryParam("apikey", "apikey").
		SetQueryParam("Category[]", "2000").
		SetResult(&JackettResult).
		Get("/indexers/all/results")
	if err != nil {
		return err
	} else if res.IsSuccess() == false {
		return domain.ErrJackettRequestFailed
	}

	jackett.PrintResponse(JackettResult)

	_, err = s.store.UserAddMovieToCollection(ctx, store.UserAddMovieToCollectionParams{
		UserID:    userId,
		MovieID:   int32(movieId),
		TorrentID: "",
		Length:    0,
	})

	return err
}

func (s *MovieService) RemoveFromCollection(ctx context.Context, userId pgtype.UUID, movieId int) error {
	_, err := s.store.GetMovieInCollection(ctx, store.GetMovieInCollectionParams{
		UserID:  userId,
		MovieID: int32(movieId),
	})
	if err != nil {
		return domain.ErrMovieNotInCollection
	}

	_, err = s.store.DeleteMovieFromCollection(ctx, store.DeleteMovieFromCollectionParams{
		UserID:  userId,
		MovieID: int32(movieId),
	})

	return err
}
