package util

var Token = "qdDY4MC2VgbCIgslZn+4l2KxCNdHbt414Kt65eXaLiei6sMVUDQCoaPN1zoAoAV//t9Oy+bJ5x2Ae6YaCAChfvQIuWofFVYBrKCf4RMunkmZNsU610fIHF4TDuvZfq97K/B+O/lkE3t12KQpqOWBxgdB04t89/1O/w1cDnyilFU="

type WebhookEvent struct {
	Type            string          `json:"type"`
	Timestamp       int64           `json:"timestamp"`
	Source          Source          `json:"source"`
	ReplyToken      string          `json:"replyToken,omitempty"`
	Mode            string          `json:"mode"`
	WebhookEventID  string          `json:"webhookEventId"`
	DeliveryContext DeliveryContext `json:"deliveryContext"`
	Message         *Message        `json:"message,omitempty"`
}

type Source struct {
	Type   string `json:"type"`
	UserID string `json:"userId"`
}

type DeliveryContext struct {
	IsRedelivery bool `json:"isRedelivery"`
}

type Message struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Text string `json:"text"`
}

type WebhookPayload struct {
	Destination string         `json:"destination"`
	Events      []WebhookEvent `json:"events"`
}
