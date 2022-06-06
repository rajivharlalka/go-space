package utils

import (
	"os"

	"github.com/slack-go/slack"
)

var Api = slack.New(os.Getenv("SLACKBOT_AUTH_TOKEN"))
