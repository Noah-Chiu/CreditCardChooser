package main

import (
	"credit-card-chooser/sql"
	"credit-card-chooser/util"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func Chooser(g *gin.Context) {
	request := WebhookPayload{}

	g.ShouldBind(&request)

	req := fmt.Sprintf("%+v", request)

	fmt.Println(req)

	for _, event := range request.Events {
		fmt.Println("收到的訊息:", event.Message.Text)
		event.chooseCard()
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

func (event *WebhookEvent) chooseCard() {
	type cards struct {
		CardNo     string
		CardNm     string
		Bank       string
		DRewards   float64
		ORewards   float64
		StartDate  string
		ExpireDate string
		Note       string
	}

	db := sql.GetDB()

	cardsData := []cards{}
	res := Reply{}

	res.ReplyToken = event.ReplyToken

	db.Debug().Find(&cardsData)
	text := strings.Split(event.Message.Text, " ")
	partner, country := "", ""
	if len(text) > 1 {
		country = text[1]
	}
	partner = text[0]
	fmt.Println("partner", partner)
	bestCard := ""
	max := 0.0
	for _, card := range cardsData {
		if country != "" && country != "台灣" && country != "臺灣" && strings.ToUpper(country) != "TW" && strings.ToUpper(country) != "TAIWAN" {
			if card.CardNo == "002" {
				if country != "日本" && strings.ToUpper(country) != "JP" && strings.ToUpper(country) != "JAPAN" {
					card.ORewards = 1
				}
			}
			if card.ORewards > max {
				bestCard = card.CardNm + " " + card.Note
				max = card.ORewards
			}
			continue
		}
		if card.DRewards > max {
			bestCard = card.CardNm + " " + card.Note
			max = card.DRewards
		}
	}
	msg := Message{
		Type:       "text",
		Text:       bestCard,
		QuoteToken: event.Message.QuoteToken,
	}
	res.Messages = append(res.Messages, msg)

	fmt.Println("回傳的訊息:", res.Messages)

	ResLine(res, util.ChooserToken)
}
