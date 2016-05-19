build:
	go build ./cmd/github-review-bot

run: build
	heroku local web

.PHONY: build run
