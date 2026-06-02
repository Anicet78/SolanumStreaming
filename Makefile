include .env
export

all: up

front:
	cd frontend && npm run dev

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

clean:
	docker compose down -v

fclean:
	docker compose down --rmi all -v

re: clean up

.PHONY: all front migrates sqlc test deploy up down clean fclean re
