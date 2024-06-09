package main

import (
	"credit-card-chooser/sql"
	"credit-card-chooser/util"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

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
type partners struct {
	CardNo     string
	Partner    string
	Rewards    float64
	StartDate  string
	ExpireDate string
	Note       string
}

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
	db := sql.GetDB()

	cardsData := []cards{}
	res := Reply{}

	res.ReplyToken = event.ReplyToken

	text := strings.Split(event.Message.Text, " ")
	partner, country := "", ""
	if len(text) > 1 {
		country = text[1]
	}
	partner = text[0]

	// 取得所有卡別
	db.Debug().Order(`"d_rewards"`).Find(&cardsData)
	rewardsType := "國內消費"
	bestCardInfo := ""
	secondCardInfo := ""
	max := [2]float64{}

	// 跑過所有卡別
	for _, card := range cardsData {
		// --------------------------國外消費--------------------------
		// 非台灣就是國外
		if country != "" && country != "台灣" && country != "臺灣" && strings.ToUpper(country) != "TW" && strings.ToUpper(country) != "TAIWAN" {
			rewardsType = fmt.Sprintf("國外消費(%s)", country)
			// 吉鶴卡只有日本有優惠
			if card.CardNo == "002" {
				if country != "日本" && strings.ToUpper(country) != "JP" && strings.ToUpper(country) != "JAPAN" {
					card.ORewards = 1
				}
			}
			if card.ORewards >= max[0] {
				secondCardInfo = bestCardInfo
				bestCardInfo = fmt.Sprintf("卡別: %s\n總回饋: %.1f%%\n備註: %s", card.CardNm, card.ORewards, card.Note)
				max[1] = max[0]
				max[0] = card.ORewards
			}
			continue
		}

		// --------------------------國內消費--------------------------

		partnersData := partners{}
		addonRewards := 0.0
		note := ""
		// 尋找合作商家
		result := db.Debug().
			Where(`"card_no" = ? AND "partner" ilike ?`, card.CardNo, "%"+partner+"%").
			Order(`"rewards" desc`).
			Limit(1).
			Find(&partnersData)
		if result.RowsAffected > 0 {
			rewardsType = fmt.Sprintf("國內消費(%s)", partnersData.Partner)
			addonRewards = partnersData.Rewards
			note = partnersData.Note
		}
		totalRewards := card.DRewards + addonRewards
		if totalRewards >= max[0] {
			secondCardInfo = bestCardInfo
			bestCardInfo = fmt.Sprintf("卡別: %s\n總回饋: %.1f%%\n備註: %s", card.CardNm, totalRewards, note)

			max[1] = max[0]
			max[0] = totalRewards
		} else if totalRewards >= max[1] {
			secondCardInfo = fmt.Sprintf("卡別: %s\n總回饋: %.1f%%\n備註: %s", card.CardNm, totalRewards, note)

			max[1] = totalRewards
		}
	}

	msg := Message{
		Type:       "text",
		Text:       rewardsType + fmt.Sprintf("\n1.\n%s\n2.\n%s", bestCardInfo, secondCardInfo),
		QuoteToken: event.Message.QuoteToken,
	}
	res.Messages = append(res.Messages, msg)

	fmt.Println("回傳的訊息:", res.Messages)

	ResLine(res, util.ChooserToken)
}

// func AddCard() {
// 	db := sql.GetDB()
// 	partner := partners{
// 		CardNo:     "001",
// 		Partner:    "Pi",
// 		Rewards:    0,
// 		StartDate:  "2024/3/1",
// 		ExpireDate: "2024/12/31",
// 		Note:       "每月回饋上限300 P幣",
// 	}
// }
