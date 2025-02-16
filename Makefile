postgres:
	@docker run --name postgres17 -p 5432:5432 -e POSTGRES_DB=cookinator -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:17.2-alpine

createdb:
	@docker exec -it postgres17 createdb --username=root --owner=root cookinator

dropdb:
	@docker exec -it postgres17 dropdb cookinator

migrateup:
	@goose up -env app.env

migratedown:
	@goose down -env app.env

install:
	@go install github.com/pressly/goose/v3/cmd/goose@latest
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	@go install github.com/air-verse/air@latest
	@go install github.com/vektra/mockery/v2@v2.51.1

sqlc:
	@sqlc generate

test:
	@go test ./... -cover

server:
	@air

mock:
	@mockery

.PHONY: createdb dropdb install postgres migrateup migratedown sqlc test server mockery
