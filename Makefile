run: build
	@./bin/qrk $(FILE)

build:
	@go build -o bin/qrk ./cmd

test:
	@go test ./... -v
