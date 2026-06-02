-- name: UserAddFilmToCollection :one
INSERT INTO collection (user_id, movie_id, torrent_id, movie_lenght)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFilmInCollection :one
SELECT * FROM collection
WHERE user_id=$1 AND movie_id=$2;
