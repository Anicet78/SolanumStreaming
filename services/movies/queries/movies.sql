-- name: UserAddMovieToCollection :one
INSERT INTO collection (user_id, movie_id, torrent_id, length)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetMovieInCollection :one
SELECT * FROM collection
WHERE user_id=$1 AND movie_id=$2;

-- name: GetAllMoviesInCollection :many
SELECT movie_id, torrent_id, length, progression FROM collection
WHERE user_id=$1;

-- name: DeleteMovieFromCollection :one
DELETE FROM collection
WHERE user_id=$1 AND movie_id=$2
RETURNING *;