package main

import (
	"fmt"
	"math/rand"
	"regexp"
	"time"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

func handleIssueCommentEvent(event GithubIssueCommentPayload) error {
	reviewersAskCommentRegexp := regexp.MustCompile(`^@` + BOT_NAME + `[\s]+assign[\s]+([a-z]+)[\s]+reviewers`)

	matches := reviewersAskCommentRegexp.FindStringSubmatch(event.Comment.Body)

	// No match, skip this issue comment
	if matches == nil {
		return nil
	}
	matchedTeamName := matches[1]

	// Find the team that is targeted
	var team *Team = nil
	for _, t := range TEAMS {
		if t.Name == matchedTeamName {
			team = &t
		}
	}
	// The given team name doesn't exist
	if team == nil {
		comment := fmt.Sprintf("'%v' is not a team I know about!", matchedTeamName)
		return createGithubIssueComment(event, comment)
	}

	// Attempt to pick a senior and a junior
	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	senior := ""
	if len(team.Seniors) > 0 {
		senior = team.Seniors[random.Intn(len(team.Seniors))]
	}
	junior := ""
	if len(team.Juniors) > 0 {
		junior = team.Juniors[random.Intn(len(team.Juniors))]
	}

	var comment string
	if senior != "" && junior != "" {
		comment = fmt.Sprintf("Thanks for submitting. @%v and @%v: please review!", senior, junior)
	} else if senior != "" {
		comment = fmt.Sprintf("I propose @%v as reviewer!", senior)
	} else if junior != "" {
		comment = fmt.Sprintf("I propose @%v as reviewer!", junior)
	} else {
		comment = fmt.Sprint("I couldn't find anybody to propose as reviewers!")
	}
	return createGithubIssueComment(event, comment)
}

func createGithubIssueComment(event GithubIssueCommentPayload, body string) error {
	client := createGithubClient()
	_, _, err := client.Issues.CreateComment(
		event.Repository.Owner.Login,
		event.Repository.Name,
		event.Issue.Number,
		&github.IssueComment{Body: &body},
	)
	return err
}

func createGithubClient() *github.Client {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GITHUB_TOKEN},
	)
	tc := oauth2.NewClient(oauth2.NoContext, ts)

	return github.NewClient(tc)
}
