run: build
	@./bin/qrk

build:
	@go build -o bin/qrk ./cmd

test:
	@go test ./... -v
