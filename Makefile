build:
	go build ./cmd/github-review-bot

run: build
	heroku local web

clean:
	rm github-review-bot

.PHONY: build run clean
