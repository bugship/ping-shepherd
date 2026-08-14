.PHONY: run test tidy fmt

BIN := .bin/coop

run:
	go run ./cmd/coop

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w .

build: $(BIN)

$(BIN):
	mkdir -p .bin
	go build -o $(BIN) ./cmd/coop
