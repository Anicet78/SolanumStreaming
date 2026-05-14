include .env
export

all:
	docker compose up -d

front:
	cd frontend && npm run dev

auth:
	cd services/auth && go run cmd/main.go

movies:
	cd services/movies && go run cmd/main.go

stream:
	cd services/movies && npm run dev

migrate:
	docker compose up -d migrate

sqlc:
	cd services/auth && sqlc generate

up:
	docker compose up -d

down:
	docker compose down

clean: down

fclean:
	docker compose down --rmi all -v

re: down up

.PHONY: all front auth movies stream migrates sqlc up down clean fclean re
