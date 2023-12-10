run: build
	@./bin/quasark

build:
	@go build -o bin/quasark ./cmd

test:
	@go test ./... -v
