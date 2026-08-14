.PHONY: run test tidy fmt

BIN := .bin/ping-shepherd

run:
	go run ./cmd/shepherd

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w .

build: $(BIN)

$(BIN):
	mkdir -p .bin
	go build -o $(BIN) ./cmd/shepherd
