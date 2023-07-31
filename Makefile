.PHONY: build test lint clean

clean:
	@cd backend && go mod tidy

build: clean
	@cd backend && go build -o collector ./cmd/collector

test: clean
	@cd backend && go test ./... -count=1

lint: clean
	@cd backend && golangci-lint run
