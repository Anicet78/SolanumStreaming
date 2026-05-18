include .env
export

all: up

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
	cd services/movies && sqlc generate

test:
	cd bruno && bru run

deploy: # k8s

up:
	docker compose up --build -d

down:
	docker compose down

clean: down

fclean:
	docker compose down --rmi all -v

re:
	docker compose down -v
	docker compose up --build -d

.PHONY: all front auth movies stream migrates sqlc test deploy up down clean fclean re
