# Github Review Bot

_A Github bot helping you pick reviewers for Pull Requests_

## Intro

This simple bot will pick reviewers for you in a PR. You can summon him
by writing something like `@ghreviewbot assign frontend reviewers` where
the team name could be any team you configured. I'll respond with 2 randomly
selected reviewers, one senior team member and one junior.

## How it works

This bot is written in Go and requires mainly 3 things to work:

- This repo's small api/webapp to be deployed on a server with some proper environment variables
- A webhook that triggers on the `issue_comment` event in your Github repository's settings
- A Github user (not your personal account about but one for the bot, preferably) to impersonate the bot

Here's what happens on a typical use of this bot:

- A comment is written on an issue or PR of a GitHub repo which you configured to ping the bot via webhook.
- GitHub creates a payload for the event and makes a `POST` request to the configured `https://boturl/hooks/github` from the webhook settings
- The bot take over and make sure the event it got is of the `issue_comment` type and that it was `created` (not `modified` or `deleted`)
- The bot matches the comment against a regex that checks if somebody really is asking him for reviewers
- With the team name matched in the previous regex it looks for a team configured for that name, goes on if it found one or comments that it doesn't know that team otherwise
- Now that it found the concerned team it selects one senior username and one junior username at random
- With those two usernames is now calls the Github API using a personal token from the bot's Github account to create comment on the concerned issue or PR

There we go, a simple bot for narrow and simple use case.

## Configuring

This bot requires to following environment variables to be in a `.env` at the root when developing or on you server configuration in production:

- `GITHUB_TOKEN`: A personal access token to make Github API calls (usually from the bot account)
- `BOT_NAME`: The name you want you bot to answer to (usually the bot's account username)
- `BOT_TEAM_<any-team-name>`: The list of reviewers Github usernames that the bot will select from. Individual usernames separated by commas (`,`) and seniors separated from the juniors list by a pipe (`|`)

Next, you also need a webhook that looks a bit like the following screenshot on the repo you want your bot to be watching for comments:

![Example Webhook](https://raw.githubusercontent.com/kiasaki/github-review-bot/master/example_webhook.png)

## Deploying

This bot is easily deployable to **Heroku** either by using the button below or creating a new app in the dashboard and deploying this repo to it.

[![Deploy](https://www.herokucdn.com/deploy/button.svg)](https://heroku.com/deploy)

## Developing

To develop a new feature for this bot you can use a combination of:

- Running this bot on your computer using `make run` & a filled in `.env` file.
- Running `ngrok` and creating a webhook on a test repo with that ngrok url in it
- Creating a test issue and commenting on it to test the bot answers
- If no answers come in you can check the server logs of the API calls responses in the Github webhook settings page of your test repo

That equates to the following commands:

```
cat > .env <<EOF
>GITHUB_TOKEN=<a valid github personal token>
>BOT_NAME=<the name you want to bot to have>
>BOT_TEAM_SOMENAME=senior1,senior2|junior1,junior2,junior3,junior4
>EOF
make run
ngrok start 5000
```

## Contributing

All contributions are welcomed and will be accepted in the form of pull requests or issues on this repo.

If you want to change how bot behaves in a major way without keeping backwards compatibility for the initial use case consider forking this repository and going from there.

## License

See the `LICENSE` file.
