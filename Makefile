run: build
	@./bin/i-go

build:
	@go build -o bin/i-go

test:
	@go test ./... -v
