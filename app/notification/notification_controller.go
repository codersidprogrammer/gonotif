package notification

import (
	"context"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	platform "github.com/codersidprogrammer/gonotif/platform/cache"
)

var ctx = context.Background()

func NewNotification(notifType NotificationType) AppNotification {
	return &Notification{
		Type: notifType,
	}
}

// RegisterTopic implements AppNotification.
func (*Notification) RegisterTopic(topic string) {
	_, err := platform.RedisConnection.SAdd(ctx, topic).Result()
	utils.ExitIfErr(err, "Failed to register topic")
}

func (*Notification) UnregisterTopic(topic string) {
	_, err := platform.RedisConnection.SRem(ctx, topic).Result()
	utils.ExitIfErr(err, "Failed to unregister")
}

// Publish implements AppNotification.
func (*Notification) Publish(topic string, message interface{}) {
	panic("unimplemented")
}

// Subscribe implements AppNotification.
func (*Notification) Subscribe(topic string, handler func(message interface{})) {
	panic("unimplemented")
}

// UnregisterTopic implements AppNotification.
