package routes

import (
	"encoding/json"
	"net/http"

	"rajivharlalka/imagery-v2/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

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
	action_1 := slack.AttachmentAction{Name: createResponseString(data.User, data.Channels[0], data.ID, data.URLPrivate, "test_comment", data.Timestamp.String()), Value: "Yes", Text: "Yes,Save Space", Type: "button"}
	action_2 := slack.AttachmentAction{Name: "No", Value: "no", Text: "No,This Image is Private", Type: "button", Style: "danger"}
	ephmeral_attachment_data.Actions = []slack.AttachmentAction{action_1, action_2}

	_, err := utils.Api.PostEphemeral(data.Channels[0], data.User, slack.MsgOptionAttachments(ephmeral_attachment_data))
	if err != nil {
		panic(err)
	}
}
