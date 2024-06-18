package main

import (
	"credit-card-chooser/util"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Test(g *gin.Context) {
	request := WebhookPayload{}

	g.ShouldBind(&request)

	req := fmt.Sprintf("%+v", request)

	fmt.Println(req)

	for _, event := range request.Events {
		fmt.Println("收到的訊息:", event.Message.Text)
		event.testResponse()
	}

	g.JSON(200, struct {
		Status uint16      `json:"status"`
		Msg    string      `json:"msg"`
		Data   interface{} `json:"data"`
	}{
		Status: 200,
		Msg:    "ok",
		Data:   req,
	})
}

func (event *WebhookEvent) testResponse() {
	res := Reply{}

	res.ReplyToken = event.ReplyToken

	msg := Message{
		Type: "text",
		// Text:       "hello, this is test",
		QuoteToken: event.Message.QuoteToken,
	}

	if strings.Contains(event.Message.Text, "大湯匙") {
		msg.Text = "小湯匙我愛妳$"
		emoji := Emoji{
			Index:     6,
			ProductId: "5ac1bfd5040ab15980c9b435",
			EmojiId:   "215",
		}
		msg.Emojis = append(msg.Emojis, emoji)
	}

	// funny tool
	for _, r := range event.Message.Text {
		if r >= 'A' && r <= 'Z' {
			msg.Text += "$"
			msg.Emojis = append(msg.Emojis, Emoji{
				Index:     len(msg.Text) - 1,
				ProductId: "5ac21a8c040ab15980c9b43f",
				EmojiId:   util.IntToDigits(int(r)-64, 3),
			})
		}
		if r >= 'a' && r <= 'z' {
			msg.Text += "$"
			msg.Emojis = append(msg.Emojis, Emoji{
				Index:     len(msg.Text) - 1,
				ProductId: "5ac21a8c040ab15980c9b43f",
				EmojiId:   util.IntToDigits(int(r)-70, 3),
			})
		}
		if r >= '0' && r <= '9' {
			if r == '0' {
				msg.Text += "$"
				msg.Emojis = append(msg.Emojis, Emoji{
					Index:     len(msg.Text) - 1,
					ProductId: "5ac21a8c040ab15980c9b43f",
					EmojiId:   "062",
				})
			} else {
				msg.Text += "$"
				msg.Emojis = append(msg.Emojis, Emoji{
					Index:     len(msg.Text) - 1,
					ProductId: "5ac21a8c040ab15980c9b43f",
					EmojiId:   util.IntToDigits(int(r)+4, 3),
				})
			}
		}
	}

	if msg.Text == "" {
		msg.Text = "hello, this is test"
	}

	res.Messages = append(res.Messages, msg)

	fmt.Println("回傳的訊息:", msg.Text)

	Response(res, util.TestToken)
}
