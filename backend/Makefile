.PHONY: run build test test-cover test-integration vet fmt lint tidy

run:
	go run ./cmd/api

build:
	go build -o bin/api ./cmd/api

test:
	go test ./...

test-cover:
	go test $$(go list ./... | grep -v /internal/mocks | grep -v /internal/infrastructure/repository | grep -v /internal/infrastructure/database | grep -v /cmd) -coverprofile=coverage.out
	go tool cover -func=coverage.out | tail -1
	go tool cover -html=coverage.out -o coverage.html

test-integration:
	go test -tags=integration ./internal/infrastructure/...

vet:
	go vet ./...

fmt:
	gofmt -w .

tidy:
	go mod tidy
