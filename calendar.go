package main

import (
	"credit-card-chooser/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Schedule(g *gin.Context) {
	request := WebhookPayload{}
	db := sql.GetDB()
	g.ShouldBind(&request)

	req := fmt.Sprintf("%+v", request)

	fmt.Println(req)

	for _, event := range request.Events {
		fmt.Println("收到的訊息:", event.Message.Text)
		addSchedule(db, event)
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

func addSchedule(db *gorm.DB, event WebhookEvent) {
	cal := Calendar{
		UserID: event.Source.UserID,
	}
	result := db.Debug().Create(&cal)

	if result.Error != nil {
		res := Reply{
			ReplyToken: event.ReplyToken,
			Messages: []Message{{
				Type:       "text",
				Text:       "設定失敗",
				QuoteToken: event.Message.QuoteToken,
			}},
		}

		ResLine(res)

		fmt.Println(result.Error.Error())

		return
	}

	res := Reply{
		ReplyToken: event.ReplyToken,
		Messages: []Message{{
			Type:       "text",
			Text:       "設定完成",
			QuoteToken: event.Message.QuoteToken,
		}},
	}

	ResLine(res)
}
