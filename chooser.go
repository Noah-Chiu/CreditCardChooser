package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

type cardInfo struct {
	name            string
	domesticRewards float64
	overseasRewards float64
	partner         map[string]float64
}

func Chooser(g *gin.Context) {
	request := WebhookPayload{}
	cards := []cardInfo{
		{
			name:            "test1",
			domesticRewards: 1.5,
			overseasRewards: 2,
			partner: map[string]float64{
				"麥當勞": 5,
				"肯德基": 2,
				"漢堡王": 10,
			},
		},
		{
			name:            "test2",
			domesticRewards: 1,
			overseasRewards: 3,
			partner: map[string]float64{
				"麥當勞": 1,
				"肯德基": 2,
				"漢堡王": 3,
			},
		},
		{
			name:            "test3",
			domesticRewards: 0.5,
			overseasRewards: 5,
			partner: map[string]float64{
				"麥當勞": 2,
				"肯德基": 2,
				"漢堡王": 2,
			},
		},
	}

	g.ShouldBind(&request)

	req := fmt.Sprintf("%+v", request)

	fmt.Println(req)

	for _, event := range request.Events {
		fmt.Println("收到的訊息:", event.Message.Text)
		event.chooseCard(cards)
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

func (event *WebhookEvent) chooseCard(cards []cardInfo) {
	res := Reply{}

	res.ReplyToken = event.ReplyToken

	msg := Message{
		Type:       "text",
		Text:       "hello, this is test",
		QuoteToken: event.Message.QuoteToken,
	}

	res.Messages = append(res.Messages, msg)

	fmt.Println("回傳的訊息:", msg.Text)

	ResLine(res)
}
