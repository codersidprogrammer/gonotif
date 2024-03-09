package notification

type NotificationType int

const (
	Push NotificationType = iota
	Email
)

type Notification struct {
	Message interface{}      `json:"message"`
	Type    NotificationType `json:"type"`
}

type AppNotification interface {
	RegisterTopic(topic string)
	UnregisterTopic(topic string)
	Publish(topic string, message interface{})
	Subscribe(topic string, handler func(message interface{}))
}
