all:
	docker compose up -d

front:
	cd frontend && npm run dev

auth:
	cd services/auth && go run cmd/main.go

users:
	cd services/users && go run cmd/main.go

movies:
	cd services/movies && go run cmd/main.go

stream:
	cd services/movies && npm run dev
