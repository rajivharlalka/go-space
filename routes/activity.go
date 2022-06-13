package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"rajivharlalka/go-space/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/slack-go/slack"
)

func ActivityRoute(c *fiber.Ctx) error {
	defer utils.RecoverServer()

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

func downloadFile(permaLink string, file_id string, channel_id string, comment string, user_id string, timestamp string) {
	defer utils.RecoverServer()
	// types, err := os.Create("hello123.png")
	buf := new(bytes.Buffer)
	if err := utils.Api.GetFile(permaLink, buf); err != nil {
		panic(err)
	}
	// upload to imgur
	token := "Client-ID " + os.Getenv("IMGUR_CLIENT_ID")

	data := upload(buf, token)

	text := "Posted By <@" + user_id + ">\n" + data.Data.Link + "\n" + comment

	utils.Api.SendMessage(channel_id, slack.MsgOptionText(text, false))

	// Delete File
	err := utils.UserAuthedApi.DeleteFile(file_id)
	if err != nil {
		panic(err)
	}

	_, _, error := utils.UserAuthedApi.DeleteMessage(channel_id, timestamp)
	if err != nil && error.Error() != "message_not_found" {
		panic(error)
	}
}

func upload(image io.Reader, token string) utils.Imgur_data {
	defer utils.RecoverServer()

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
		panic(err)
	}
	return data
}
