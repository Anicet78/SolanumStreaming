FROM docker.io/migrate/migrate
COPY auth/migrations /migrations/auth
COPY movies/migrations /migrations/movies
