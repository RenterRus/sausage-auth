build:
	@go build cmd/main.go

SQLC_VERSION := 1.28.0

gen proto-v1: ### generate source files from proto
	protoc --go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		docs/proto/v1/*.proto
.PHONY: proto-v1

update:
	@git pull && make build

install goose:
	@go get -u github.com/pressly/goose/v3/cmd/goose@latest

create migration:
	@go run github.com/pressly/goose/v3/cmd/goose@latest create refreshlist sql -dir internal/repo/migrations


MIGRATIONS_DIR ?= ./internal/repo/migrations
DB_DRIVER ?= postgres

CONN_STR = $(DB_CONNECTION_STR)

migrate-set-conn-str-example:
	export DB_CONNECTION_STR=postgres://user:password@localhost:5432/dbname?sslmode=disable

migrate-up:
	echo $(CONN_STR)
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir $(MIGRATIONS_DIR) $(DB_DRIVER) $(CONN_STR) up

generate-sqlc:
	docker run --rm -v $(shell pwd):/src  sqlc/sqlc:$(SQLC_VERSION) generate -f /src/sqlc.yaml
