.PHONY: dev build seed run

dev:
	go run ./cmd/server

build:
	go build -o bin/server ./cmd/server

seed:
	rm -f data/ai-learning.db
	go run ./cmd/server

run: build
	./bin/server
