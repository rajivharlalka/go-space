package utils

import (
	"os"

	"github.com/slack-go/slack"
)

var Api = slack.New(os.Getenv("SLACKBOT_AUTH_TOKEN"))
var UserAuthedApi = slack.New(os.Getenv("SLACKUSER_AUTH_TOKEN"))
