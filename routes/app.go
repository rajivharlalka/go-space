package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"rajivharlalka/imagery-v2/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

const username_regex = "Posted By <@+.*>"

func RootRoute(c *fiber.Ctx) error {
	body := c.Body()

	eventsAPIEvent, err := slackevents.ParseEvent(json.RawMessage(body), slackevents.OptionNoVerifyToken())
	if err != nil {
		return c.SendStatus(http.StatusInternalServerError)
	}

	if eventsAPIEvent.Type == slackevents.URLVerification {
		var r *slackevents.ChallengeResponse
		err := json.Unmarshal([]byte(body), &r)
		if err != nil {
			return c.SendStatus(http.StatusInternalServerError)
		}
		c.Append("Content-Type", "text")
		return c.Send([]byte(r.Challenge))
	}

	if eventsAPIEvent.Type == slackevents.CallbackEvent {
		innerEvent := eventsAPIEvent.InnerEvent
		switch ev := innerEvent.Data.(type) {
		case *slack.FileSharedEvent:

			data, _, _, err := utils.Api.GetFileInfo(ev.FileID, 0, 0)
			if err != nil {
				panic(err)
			}

			go sendEphemeral(data)
		case *slackevents.ReactionAddedEvent:
			if ev.Reaction == "x" {
				var conversation_history_data slack.GetConversationHistoryParameters
				conversation_history_data.ChannelID = ev.Item.Channel
				conversation_history_data.Inclusive = true
				conversation_history_data.Limit = 1
				conversation_history_data.Latest = ev.Item.Timestamp
				data, err := utils.Api.GetConversationHistory(&conversation_history_data)
				if err != nil {
					fmt.Println(err)
				} else {
					r, _ := regexp.Compile(username_regex)
					match_string := r.FindString(data.Messages[0].Text)
					if len(match_string) == 0 {
						return nil
					}
					posted_user_id := match_string[12 : len(match_string)-1]
					if posted_user_id == ev.User {
						_, _, err := utils.Api.DeleteMessage(ev.Item.Channel, ev.Item.Timestamp)
						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}

	return c.SendString("Hello, World 👋!")
}

func createResponseString(user_id string, channel_id string, file_id string, file_permalink string, comment string, timestamp string) string {
	return "response|" + user_id + "|" + channel_id + "|" + file_id + "|" + file_permalink + "|" + comment + "|" + timestamp
}

func sendEphemeral(data *slack.File) {
	var ephmeral_attachment_data slack.Attachment
	ephmeral_attachment_data.CallbackID = "ephemeral_action"
	ephmeral_attachment_data.Text = "Would you like to upload this image to imgur?"
	// GET COMMENT ON IMAGE
	action_1 := slack.AttachmentAction{Name: createResponseString(data.User, data.Channels[0], data.ID, data.URLPrivate, "test_comment", data.Shares.Public[data.Channels[0]][0].Ts), Value: "Yes", Text: "Yes,Save Space", Type: "button"}
	action_2 := slack.AttachmentAction{Name: "No", Value: "no", Text: "No,This Image is Private", Type: "button", Style: "danger"}
	ephmeral_attachment_data.Actions = []slack.AttachmentAction{action_1, action_2}

	_, err := utils.Api.PostEphemeral(data.Channels[0], data.User, slack.MsgOptionAttachments(ephmeral_attachment_data))
	if err != nil {
		panic(err)
	}
}
