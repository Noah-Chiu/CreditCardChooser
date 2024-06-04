package main

import "time"

// 接收格式
type WebhookPayload struct {
	Destination string         `json:"destination"`
	Events      []WebhookEvent `json:"events"`
}

// 事件格式
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

// 回覆格式
type Reply struct {
	ReplyToken string    `json:"replyToken"`
	Messages   []Message `json:"messages"`
}

// 訊息
type Message struct {
	Type       string  `json:"type"`
	ID         string  `json:"id"`
	QuoteToken string  `json:"quoteToken"`
	Text       string  `json:"text"`
	Emojis     []Emoji `json:"emojis"`
}

type Emoji struct {
	Index     int    `json:"index"`
	ProductId string `json:"productId"`
	EmojiId   string `json:"emojiId"`
}

// 行事曆DB
type Calendar struct {
	UserID   string    `json:"userID" gorm:"column:userID"`
	PushTime time.Time `json:"pushTime" gorm:"column:pushTime"`
	PushMsg  string    `json:"pushMsg" gorm:"column:pushMsg"`
}
