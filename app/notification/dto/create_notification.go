package dto

type NotificationType int

const (
	Push NotificationType = iota
	Email
)

type CreateNotificationRequest struct {
	Type    NotificationType `json:"type"`
	Payload interface{}      `json:"payload"`
}

type CreatePushNotificationRequest struct {
	Topic   string      `json:"topic"`
	Payload interface{} `json:"payload"`
}
