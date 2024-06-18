package main

import (
	"credit-card-chooser/sql"
	"credit-card-chooser/util"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

func Schedule(g *gin.Context) {
	request := WebhookPayload{}
	g.ShouldBind(&request)

	req := fmt.Sprintf("%+v", request)

	fmt.Println(req)

	for _, event := range request.Events {
		fmt.Println("收到的訊息:", event.Message.Text)
		event.addSchedule()
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

func (event *WebhookEvent) addSchedule() {
	type Calendars struct {
		UserId   string
		PushTime time.Time
		PushMsg  string
	}

	db := sql.GetDB()

	cal := Calendars{
		UserId:   event.Source.UserID,
		PushTime: time.Now().Local(),
		PushMsg:  event.Message.Text,
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

		Response(res, util.CalendarToken)

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

	Response(res, util.CalendarToken)
}
