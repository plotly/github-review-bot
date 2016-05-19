package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type J map[string]interface{}

type Team struct {
	Name    string
	Seniors []string
	Juniors []string
}

type GithubUser struct {
	Id        int    `json:"id"`
	Login     string `json:"login"`
	AvatarUrl string `json:"avatar_url"`
	Type      string `json:"type"`
}

type GithubComment struct {
	Id   int        `json:"id"`
	Body string     `json:"body"`
	User GithubUser `json:"user"`
}

type GithubIssue struct {
	Id     int        `json:"id"`
	Number int        `json:"number"`
	Title  string     `json:"title"`
	State  string     `json:"state"`
	Body   string     `json:"body"`
	User   GithubUser `json:"user"`
}

type GithubRepository struct {
	Id          int        `json:"id"`
	Name        string     `json:"name"`
	FullName    string     `json:"full_name"`
	Owner       GithubUser `json:"owner"`
	Private     bool       `json:"private"`
	Fork        bool       `json:"fork"`
	Description string     `json:"decription"`
	GitUrl      string     `json:"git_url"`
	SshUrl      string     `json:"ssh_url"`
	CloneUrl    string     `json:"clone_url"`
}

type GithubIssueCommentPayload struct {
	Action     string           `json:"action"`
	Issue      GithubIssue      `json:"issue"`
	Comment    GithubComment    `json:"comment"`
	Repository GithubRepository `json:"repository"`
	Sender     GithubUser       `json:"sender"`
}

var PORT string
var GITHUB_TOKEN string
var BOT_NAME string

func init() {
	PORT = os.Getenv("PORT")
	if PORT == "" {
		log.Fatal("$PORT must be set")
	}

	GITHUB_TOKEN = os.Getenv("GITHUB_TOKEN")
	if GITHUB_TOKEN == "" {
		log.Fatal("$GITHUB_TOKEN must be set")
	}

	BOT_NAME = os.Getenv("BOT_NAME")
	if BOT_NAME == "" {
		log.Fatal("$BOT_NAME must be set")
	}

	parseTeams()
}

func main() {
	router := gin.New()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, J{"teams": TEAMS})
	})

	router.POST("/hooks/github", handleHooksGithub)

	router.Run(":" + PORT)
}

func handleHooksGithub(c *gin.Context) {
	switch c.Request.Header.Get("X-GitHub-Event") {
	case "ping":
		c.JSON(http.StatusOK, J{"message": "Pong!"})
	case "issue_comment":
		payload := GithubIssueCommentPayload{}
		if c.BindJSON(&payload) != nil {
			c.JSON(http.StatusBadRequest, J{"error": "Can't parse JSON payload"})
			return
		}

		if payload.Action != "created" {
			c.JSON(http.StatusOK, J{"message": "event skipped, not create"})
			return
		}

		err := handleIssueCommentEvent(payload)
		if err != nil {
			c.JSON(http.StatusBadRequest, J{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, J{"message": "event handled"})
	default:
		c.AbortWithStatus(http.StatusNotFound)
	}
}

var TEAMS []Team

// Loops on all env vares searching for strings like:
//
// BOT_TEAM_FRONTEND=joe|jack,joey,july
// BOT_TEAM_BACKEND=jody,july|
//
// Where in that case we would have a "frontend" team with "joe" as senior,
// and "jack", "joey" and "july" as juniors. Same pattern for our second
// "backend" team.
func parseTeams() {
	const teamEnvVarPrefix = "BOT_TEAM_"

	TEAMS = []Team{}

	for _, envVar := range os.Environ() {
		if strings.HasPrefix(envVar, teamEnvVarPrefix) {
			teamParts := strings.Split(envVar[len(teamEnvVarPrefix):], "=")
			teamName := strings.ToLower(teamParts[0])

			team := Team{
				Name:    teamName,
				Seniors: []string{},
				Juniors: []string{},
			}

			teamMembersParts := strings.Split(teamParts[1], "|")

			for _, senior := range strings.Split(teamMembersParts[0], ",") {
				if len(senior) > 0 {
					team.Seniors = append(team.Seniors, senior)
				}
			}

			for _, junior := range strings.Split(teamMembersParts[1], ",") {
				if len(junior) > 0 {
					team.Juniors = append(team.Juniors, junior)
				}
			}

			TEAMS = append(TEAMS, team)
		}
	}

	if len(TEAMS) == 0 {
		log.Fatal("At least one '" + teamEnvVarPrefix + "...' variable is required")
	}
}
