.PHONY: run test tidy fmt vuln

BIN := .bin/ping-shepherd

run:
	go run ./cmd/shepherd

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -w .

vuln:
	go run golang.org/x/vuln/cmd/govulncheck@latest ./...

build: $(BIN)

$(BIN):
	mkdir -p .bin
	go build -o $(BIN) ./cmd/shepherd
