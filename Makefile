.PHONY: build test lint fmt clean

BIN := gts
PKG := ./...

build:
	go build -o $(BIN) ./cmd/gts

test:
	go test -race -count=1 $(PKG)

lint:
	golangci-lint run

fmt:
	gofmt -s -w .

clean:
	rm -f $(BIN)
	rm -rf dist/
