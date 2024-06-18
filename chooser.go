package main

import (
	"credit-card-chooser/sql"
	"credit-card-chooser/util"
	"fmt"
	"strings"
	"time"

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

type userMsgs struct {
	UserId     string
	Bot        string
	Msg        string
	UpdateTime time.Time
	CreateTime time.Time
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
	now := time.Now().Local()

	cardsData := []cards{}
	res := Reply{}

	res.ReplyToken = event.ReplyToken
	inputText := event.Message.Text

	userMsg := userMsgs{}

	// 如果沒有挾帶資料
	if event.Postback.Data == "" {
		userMsg.UserId = event.Source.UserID
		userMsg.Bot = "A"
		userMsg.Msg = inputText
		userMsg.UpdateTime = now
		result := db.Debug().Where(`"user_id" = ? AND bot = 'A'`, event.Source.UserID).Updates(&userMsg)

		if result.RowsAffected == 0 {
			userMsg.CreateTime = now
			db.Debug().Create(&userMsg)
		}

		msg := Message{
			Type:    "template",
			AltText: "Choose domestic or overseas",
			Template: Template{
				Type: "confirm",
				Text: "設定商家成功\n請選擇國內或國外",
				Actions: []Action{
					{
						Type:        "postback",
						Label:       "國內",
						DisplayText: "已選擇國內",
						Data:        "D",
					},
					{
						Type:        "postback",
						Label:       "國外",
						DisplayText: "已選擇國外",
						Data:        "O",
					},
				},
			},
		}
		res.Messages = append(res.Messages, msg)

		fmt.Println("回傳的訊息:", res.Messages)

		Response(res, util.ChooserToken)
		return
	}

	// ---------------------------------------------如果有儲存的訊息之後---------------------------------------------
	// 相異的合作商家
	diffPartners := []string{}

	// 取得所有相異的合作商家
	db.Debug().Table(`"partners"`).
		Select(`"partner"`).
		Where(`"partner" ilike ?`, "%"+inputText+"%").
		Group(`"partner"`).
		Find(&diffPartners)

	// 取得所有卡別
	db.Debug().Order(`"d_rewards"`).Find(&cardsData)

	bestCardInfo := ""
	secondCardInfo := ""
	rewardsType := ""
	rankArray := []float64{0.0, 0.0}

	switch event.Postback.Data {
	// --------------------------國內消費--------------------------
	case "D":
		// 如果有找到合作商家各卡別要加上合作商家的回饋
		for _, partner := range diffPartners {
			// 跑過所有卡別
			for _, card := range cardsData {
				rewardsType = fmt.Sprintf("國內消費(%s)", partner)

				partnersData := partners{}
				addonRewards := 0.0
				note := ""
				// 尋找合作商家
				result := db.Debug().
					Where(`"card_no" = ? AND "partner" = ?`, card.CardNo, partner).
					Order(`"rewards" desc`).
					Limit(1).
					Find(&partnersData)
				if result.RowsAffected > 0 {
					addonRewards = partnersData.Rewards
					note = partnersData.Note
				}
				totalRewards := card.DRewards + addonRewards

				decideCards(&rankArray, &bestCardInfo, &secondCardInfo, card.CardNm, note, totalRewards)
			}

			msg := Message{
				Type:       "text",
				Text:       rewardsType + fmt.Sprintf("\n1.\n%s\n2.\n%s", bestCardInfo, secondCardInfo),
				QuoteToken: event.Message.QuoteToken,
			}
			res.Messages = append(res.Messages, msg)
		}

		// 如果沒有找到合作商家就只要找各卡別
		if len(diffPartners) == 0 {
			// 跑過所有卡別
			for _, card := range cardsData {
				rewardsType = "國內消費(無合作商家)"

				decideCards(&rankArray, &bestCardInfo, &secondCardInfo, card.CardNm, "", card.DRewards)
			}

			msg := Message{
				Type:       "text",
				Text:       rewardsType + fmt.Sprintf("\n1.\n%s\n2.\n%s", bestCardInfo, secondCardInfo),
				QuoteToken: event.Message.QuoteToken,
			}
			res.Messages = append(res.Messages, msg)
		}
	// --------------------------國外消費--------------------------
	case "O":
		rewardsType = fmt.Sprintf("國外消費(%s)", inputText)

		// 跑過所有卡別
		for _, card := range cardsData {
			// 吉鶴卡只有日本有優惠
			if card.CardNo == "002" {
				if inputText != "日本" && strings.ToUpper(inputText) != "JP" && strings.ToUpper(inputText) != "JAPAN" {
					card.ORewards = 1
				}
			}
			decideCards(&rankArray, &bestCardInfo, &secondCardInfo, card.CardNm, card.Note, card.ORewards)
		}
		msg := Message{
			Type:       "text",
			Text:       rewardsType + fmt.Sprintf("\n1.\n%s\n2.\n%s", bestCardInfo, secondCardInfo),
			QuoteToken: event.Message.QuoteToken,
		}
		res.Messages = append(res.Messages, msg)
	default:
		msg := Message{
			Type:       "text",
			Text:       "取得失敗",
			QuoteToken: event.Message.QuoteToken,
		}
		res.Messages = append(res.Messages, msg)
	}

	fmt.Println("回傳的訊息:", res.Messages)

	Response(res, util.ChooserToken)
}

// 決定最佳卡片
func decideCards(rankArray *[]float64, bestCardInfo, secondCardInfo *string, cardNm, note string, totalRewards float64) {
	if totalRewards >= (*rankArray)[0] { // 如果比最大的大要改第一+第二
		*secondCardInfo = *bestCardInfo
		*bestCardInfo = fmt.Sprintf("卡別: %s\n總回饋: %.1f%%", cardNm, totalRewards)
		if note != "" {
			*bestCardInfo += fmt.Sprintf("\n備註: %s", note)
		}
		(*rankArray)[1] = (*rankArray)[0]
		(*rankArray)[0] = totalRewards
	} else if totalRewards >= (*rankArray)[1] { // 如果只比第二的大只要改第二
		*secondCardInfo = fmt.Sprintf("卡別: %s\n總回饋: %.1f%%", cardNm, totalRewards)
		if note != "" {
			*bestCardInfo += fmt.Sprintf("\n備註: %s", note)
		}
		(*rankArray)[1] = totalRewards
	}
}
