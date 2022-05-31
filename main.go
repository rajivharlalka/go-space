package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	utils "rajivharlalka/imagery-v2/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"
)

var api = slack.New("xoxp-2609252556242-2609264879891-3466541747172-bdff3b97bf5364fe4b7805ac27a2ce78")

func createResponseString(user_id string, channel_id string, file_id string, file_permalink string, comment string, timestamp string) string {
	response := "response|" + user_id + "|" + channel_id + "|" + file_id + "|" + file_permalink + "|" + comment + "|" + timestamp
	fmt.Print(response, "\n")
	return response
}

func sendEphemeral(data *slack.File) {
	var ephmeral_attachment_data slack.Attachment
	ephmeral_attachment_data.CallbackID = "ephemeral_action"
	ephmeral_attachment_data.Text = "Would you like to upload this image to imgur?"
	fmt.Printf("PermaLink %s", data.Permalink)
	action_1 := slack.AttachmentAction{Name: createResponseString(data.User, data.Channels[0], data.ID, data.URLPrivate, "test_comment", data.Timestamp.String()), Value: "Yes", Text: "Yes,Save Space", Type: "button"}
	action_2 := slack.AttachmentAction{Name: "No", Value: "no", Text: "No,This Image is Private", Type: "button", Style: "danger"}
	ephmeral_attachment_data.Actions = []slack.AttachmentAction{action_1, action_2}

	_, err := api.PostEphemeral(data.Channels[0], data.User, slack.MsgOptionAttachments(ephmeral_attachment_data))
	if err != nil {
		fmt.Print(err)
	}
}

func rootRoute(c *fiber.Ctx) error {
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
			fmt.Print(ev.EventTimestamp, " ", ev.FileID, " ", ev.Type)

			data, _, _, err := api.GetFileInfo(ev.FileID, 0, 0)
			if err != nil {
				fmt.Print(err)
			}

			go sendEphemeral(data)
		}
	}

	return c.SendString("Hello, World 👋!")
}

func formatJSON(data []byte) ([]byte, error) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "    ")
	if err == nil {
		return out.Bytes(), err
	}
	return data, nil
}

func upload(image io.Reader, token string) utils.Imgur_data {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, _ := writer.CreateFormFile("image", "dont care about name")
	io.Copy(part, image)

	writer.Close()
	req, _ := http.NewRequest("POST", "https://api.imgur.com/3/image", buf)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", token)

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	b, _ := ioutil.ReadAll(res.Body)

	var data utils.Imgur_data
	err = json.Unmarshal(b, &data)
	if err != nil {
		fmt.Print(err)
	}
	return data
}

func downloadFile(permaLink string, file_id string, channel_id string, comment string, user_id string, timestamp string) {
	fmt.Printf("PermaLink 2 %s\n", permaLink)
	// types, err := os.Create("hello123.png")
	buf := new(bytes.Buffer)
	if err := api.GetFile(permaLink, buf); err != nil {
		fmt.Print(err)
		return
	}
	// upload to imgur
	token := "Client-ID " + os.Getenv("IMGUR_CLIENT_ID")

	data := upload(buf, token)
	fmt.Print(data)

	text := "Posted By <@" + user_id + ">\n" + data.Data.Link

	api.SendMessage(channel_id, slack.MsgOptionText(text, false))

	// Delete File
	err := api.DeleteFile(file_id)
	if err != nil {
		fmt.Print(err)
	}
	return
}

func activityRoute(c *fiber.Ctx) error {
	var types *slack.InteractionCallback
	if err := json.Unmarshal([]byte(c.FormValue("payload")), &types); err != nil {
		return err
	}

	action := types.ActionCallback.AttachmentActions[0]

	if action.Value == "Yes" {
		// Download and upload Image
		param := strings.Split(action.Name, "|")
		go downloadFile(param[4], param[3], param[2], param[5], param[1], param[6])
	} else {
		fmt.Print("---------NO SELECTED----------")
	}

	res := slack.WebhookMessage{ResponseType: "ephemeral", ReplaceOriginal: true, DeleteOriginal: true}
	return c.JSON(res)
}

func main() {
	app := fiber.New()

	app.Post("/app", rootRoute)
	app.Post("/activity-route", activityRoute)

	log.Fatal(app.Listen(":3000"))
}
