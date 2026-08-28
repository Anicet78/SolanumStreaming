package service

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func (s *MovieService) Search(ctx context.Context, params domain.SearchRequestQuery) (tmdb.SearchResponse, error) {
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

func (s *MovieService) GetCollection(ctx context.Context, userId pgtype.UUID) ([]domain.CollectionMovie, error) {
	movies, err := s.store.GetAllMoviesInCollection(ctx, userId)
	if err != nil {
		return nil, err
	}

	res := make([]domain.CollectionMovie, 0, len(movies))

	for _, u := range movies {
		res = append(res, domain.CollectionMovie{
			MovieID:     int(u.MovieID),
			TorrentLink: u.TorrentLink,
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
		TorrentLink: movie.TorrentLink,
		Length:      int(movie.Length),
		Progression: movie.Progression,
	}, nil
}

func (s *MovieService) AddToCollection(ctx context.Context, userId pgtype.UUID, movieId int) (domain.CollectionMovie, error) {
	exists, err := s.store.GetMovieInCollection(ctx, store.GetMovieInCollectionParams{
		UserID:  userId,
		MovieID: int32(movieId),
	})
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return domain.CollectionMovie{}, err
		}
	} else if exists.UserID.Valid {
		return domain.CollectionMovie{}, domain.ErrMovieAlreadyInCollection
	}

	var TMDBResult tmdb.Movie
	res, err := tmdb.New().
		SetResult(&TMDBResult).
		SetPathParam("movie_id", strconv.Itoa(movieId)).
		Get("/movie/{movie_id}")
	if err != nil {
		return domain.CollectionMovie{}, err
	} else if res.IsSuccess() == false {
		return domain.CollectionMovie{}, domain.ErrTMDBRequestFailed
	}

	var JackettResponse jackett.JackettResponse
	res, err = jackett.New().
		SetQueryParam("apikey", "apikey").
		SetQueryParam("Query", TMDBResult.Title).
		SetQueryParam("Category[]", "2000").
		SetQueryParam("Limit", "150").
		SetResult(&JackettResponse).
		Get("/indexers/all/results")
	if err != nil {
		return domain.CollectionMovie{}, err
	} else if res.IsSuccess() == false {
		return domain.CollectionMovie{}, domain.ErrJackettRequestFailed
	}

	for _, idx := range JackettResponse.Indexers {
		if idx.Error != "" || idx.Status != 2 {
			log.Printf("Indexer lent/down: %s | status: %d | error: %s", idx.Name, idx.Status, idx.Error)
		}
	}

	bestResult, err := jackett.FindBestResult(&JackettResponse, TMDBResult)
	if err != nil {
		return domain.CollectionMovie{}, err
	}

	jackett.PrintResult(bestResult.Torrent, 0)

	var link string
	if bestResult.Torrent.MagnetUri != "" {
		link = bestResult.Torrent.MagnetUri
	} else {
		link = bestResult.Torrent.Link
	}

	fmt.Printf("Movie Name: %s | MovieId: %d  |  TorrentLink: %s  |  Lenght: %d\n", TMDBResult.Title, movieId, link, TMDBResult.Runtime)

	movie, err := s.store.UserAddMovieToCollection(ctx, store.UserAddMovieToCollectionParams{
		UserID:      userId,
		MovieID:     int32(movieId),
		TorrentLink: link,
		Length:      int32(TMDBResult.Runtime),
	})

	return domain.CollectionMovie{
		MovieID:     int(movie.MovieID),
		TorrentLink: movie.TorrentLink,
		Length:      int(movie.Length),
		Progression: movie.Progression,
	}, err
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
