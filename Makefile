BIN := "./bin/antibforce"
DOCKER_IMG="antibforce:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

.PHONY: build
build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/antibforce

.PHONY: run
run: build
	$(BIN) -config ./cmd/antibforce/config.toml

.PHONY: run-db
run-db:
	docker container rm antibforce_pg -f 2>/dev/null || true
	docker run -d --name antibforce_pg \
	-e POSTGRES_USER=otus_user \
	-e POSTGRES_PASSWORD=dev_pass \
	-e POSTGRES_DB=antibforce \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
	-v pg_data_antibforce_pg:/var/lib/postgresql/data \
	-p 5432:5432 \
	postgres:14

.PHONY: build-img
build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f Dockerfile .

.PHONY: run-img
run-img: build-img
	docker run --network=host $(DOCKER_IMG)

.PHONY: swag-gen
swag-gen:
	swagger generate spec -o ./swagger.json

.PHONY: swagger
swagger: swag-gen
	swagger serve -F=swagger -p=58810 swagger.json

.PHONY: test
test:
	go test -race -count 100 ./internal/...

.PHONY: test-integration
test-integration:
	sh ./tests/test_integration.sh

.PHONY: install-lint-deps
install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.61.0

.PHONY: lint
lint: install-lint-deps
	golangci-lint run ./...

.PHONY: format
format:
	@gofmt -s -w .
	@gci write .
	@gofumpt -w .