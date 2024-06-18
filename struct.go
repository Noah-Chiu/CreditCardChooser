package main

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
	Postback        Postback        `json:"postback"`
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
	Type       string   `json:"type"`
	ID         string   `json:"id"`
	QuoteToken string   `json:"quoteToken"`
	Text       string   `json:"text"`
	Emojis     []Emoji  `json:"emojis"`
	AltText    string   `json:"altText"`
	Template   Template `json:"template"`
}

// 挾帶檔案
type Postback struct {
	Data   string `json:"data"`
	Params Params `json:"params"`
}

// 挾帶檔案
type Params struct {
	Datetime           string `json:"datetime"`
	NewRichMenuAliasId string `json:"newRichMenuAliasId"`
	Status             string `json:"status"`
}

type Emoji struct {
	Index     int    `json:"index"`
	ProductId string `json:"productId"`
	EmojiId   string `json:"emojiId"`
}

type Action struct {
	Type  string `json:"type"`
	Label string `json:"label"`
	Text  string `json:"text"`
	Uri   string `json:"uri"`
	Data  string `json:"data"`
}

type Template struct {
	Type    string   `json:"type"`
	Text    string   `json:"text"`
	Actions []Action `json:"actions"`
}
