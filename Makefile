run: build
	@./bin/i-go

build:
	@go build -o bin/i-go ./cmd

test:
	@go test ./... -v
