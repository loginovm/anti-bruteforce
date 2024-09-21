BIN := "./bin/antibforce"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/antibforce

run: build
	$(BIN) -config ./cmd/antibforce/config.toml

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.60.3

lint: install-lint-deps
	golangci-lint run ./...

test:
	go test -race ./internal/...

.PHONY: run-db
run-db:
	docker run -d --name antibforce_pg \
	-e POSTGRES_USER=otus_user \
	-e POSTGRES_PASSWORD=dev_pass \
	-e POSTGRES_DB=antibforce \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
	-v pg_data_antibforce_pg:/var/lib/postgresql/data \
	-p 5432:5432 \
	postgres:14

format:
	@gofmt -s -w .
	@gci write .
	@gofumpt -w .