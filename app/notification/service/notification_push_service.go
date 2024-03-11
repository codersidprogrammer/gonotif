package service

import (
	"context"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/codersidprogrammer/gonotif/platform/cache"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type pushNotification struct {
	bucketService NotificationBucketService
	redisClient   *redis.Client
	ctx           context.Context
}

type NotificationPushService interface {
	Subscribe(bucketId string, topic string) error
	SubsHandler() error
	SendPushNotification(topic string, msg []byte) error
}

var _topic = make(chan string)

func NewNotificationPushService() NotificationPushService {
	return &pushNotification{
		bucketService: NewNotificationBucketService(),
		redisClient:   cache.RedisConnection,
		ctx:           context.Background(),
	}
}

// SendPushNotification implements NotificationPushService.
func (p *pushNotification) SendPushNotification(topic string, msg []byte) error {
	p.redisClient.Publish(p.ctx, topic, msg)
	return nil
}

// Subscribe implements NotificationPushService.
func (p *pushNotification) Subscribe(bucketId string, topic string) error {
	nb, err := p.bucketService.Get(bucketId)
	if err != nil {
		return err
	}

	var _temp []interface{}
	for _, v := range nb.Topics {
		_temp = append(_temp, v)
	}

	b := utils.IsContainValue(_temp, topic)
	if !b {
		return fiber.NewError(fiber.ErrBadRequest.Code, "your topic is not available")
	}

	_topic <- topic

	return nil
}

func (p *pushNotification) SubsHandler() error {
	ps := p.redisClient.Subscribe(p.ctx, <-_topic)
	defer ps.Close()

	cancel := make(chan bool)
	c := ps.Channel()

	for {
		select {
		case msg := <-c:
			log.Info(msg)
			log.Info(msg.Payload)
		case <-cancel:
			return nil
		}
	}
}
