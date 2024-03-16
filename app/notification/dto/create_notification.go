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
