VERSION ?= dev
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "none")
BUILD_DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X github.com/sys-2077/memo-fast/internal/version.Version=$(VERSION) \
           -X github.com/sys-2077/memo-fast/internal/version.Commit=$(COMMIT) \
           -X github.com/sys-2077/memo-fast/internal/version.BuildDate=$(BUILD_DATE)

.PHONY: build clean test install uninstall release snapshot

build:
	go build -ldflags "$(LDFLAGS)" -o bin/memo-fast ./cmd/memo-fast

clean:
	rm -rf bin/ dist/

test:
	go test ./...

install: build
	cp bin/memo-fast /usr/local/bin/memo-fast

uninstall:
	rm -f /usr/local/bin/memo-fast

release:
	goreleaser release --clean

snapshot:
	goreleaser release --snapshot --clean
