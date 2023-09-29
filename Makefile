.PHONY: build test lint clean

clean:
	@cd backend && go mod tidy

build: clean
	@cd backend && go build -o collector ./cmd/collector
	@cd backend && go build -o service ./cmd/service
	@cd backend && go build -o search-service ./cmd/search-service

test: clean
	@cd backend && go test ./... -count=1

coverage: clean
	@cd backend && go test ./... -count=1 -coverprofile coverage.out
	@cd backend && go tool cover -html=coverage.out

lint: clean
	@cd backend && golangci-lint run
