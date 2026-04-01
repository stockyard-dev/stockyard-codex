build:
	CGO_ENABLED=0 go build -o codex ./cmd/codex/

run: build
	./codex

test:
	go test ./...

clean:
	rm -f codex

.PHONY: build run test clean
