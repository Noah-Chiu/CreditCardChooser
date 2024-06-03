package main

import (
	"bytes"
	"credit-card-chooser/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

type WebhookEvent struct {
	Type            string          `json:"type"`
	Timestamp       int64           `json:"timestamp"`
	Source          Source          `json:"source"`
	ReplyToken      string          `json:"replyToken,omitempty"`
	Mode            string          `json:"mode"`
	WebhookEventID  string          `json:"webhookEventId"`
	DeliveryContext DeliveryContext `json:"deliveryContext"`
	Message         Message         `json:"message"`
}

type Source struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

type DeliveryContext struct {
	IsRedelivery bool `json:"isRedelivery"`
}

type Message struct {
	Type       string  `json:"type"`
	ID         string  `json:"id"`
	QuoteToken string  `json:"quoteToken"`
	Text       string  `json:"text"`
	Emojis     []Emoji `json:"emojis"`
}

type WebhookPayload struct {
	Destination string         `json:"destination"`
	Events      []WebhookEvent `json:"events"`
}

type Reply struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []Message `json:"messages"`
}

type Emoji struct {
	Index     int    `json:"index"`
	ProductId string `json:"productId"`
	EmojiId   string `json:"emojiId"`
}

func main() {
	router := gin.Default()
	router.Use(Cors())
	router.POST("/", ReceiveData)
	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 是否存取cookie
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		method := c.Request.Method
		// 放行OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
		}
		// 處理請求
		c.Next()
	}
}

func ReceiveData(g *gin.Context) {
	request := WebhookPayload{}

	g.ShouldBind(&request)

	req := fmt.Sprintf("%+v", request)

	fmt.Println(req)

	for _, event := range request.Events {
		fmt.Println("收到的訊息:", event.Message.Text)
		event.Response()
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

func (event *WebhookEvent) Response() {
	var buf bytes.Buffer
	res := Reply{}

	res.ReplyToken = event.ReplyToken

	msg := Message{
		Type:       "text",
		Text:       "hello, this is test (heart)",
		QuoteToken: event.Message.QuoteToken,
	}

	if strings.Contains(event.Message.Text, "大湯匙") {
		msg.Text = "小湯匙我愛妳$"
		emoji := Emoji{
			Index:     0,
			ProductId: "5ac1bfd5040ab15980c9b435",
			EmojiId:   "215",
		}
		msg.Emojis = append(msg.Emojis, emoji)
	}

	res.Messages = append(res.Messages, msg)

	fmt.Println("回傳的訊息:", msg.Text)

	json.NewEncoder(&buf).Encode(res)

	req, err := http.NewRequest("POST", util.ReplyUrl, &buf)
	if err != nil {
		return
	}

	// Set the request content type
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer {%s}", util.Token))

	// Send the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	sitemap, _ := ioutil.ReadAll(resp.Body)

	fmt.Println("收到的回傳:", string(sitemap))

	// gjson.Get(string(sitemap), "")
}
