package repository

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2/log"
)

const (
	TENANT_KEY     = "notification:tenant"
	TENANT_CHANNEL = "notification:tenant:channel:%s"
)

type Tenant struct {
	Name              string `json:"name"`
	rdb               *redis.Client
	ChannelHandler    *redis.PubSub
	MessageChannel    chan redis.Message
	StopListenChannel chan struct{}
}

func NewTenant(name string, rd *redis.Client) *Tenant {
	return &Tenant{
		Name:              name,
		rdb:               rd,
		MessageChannel:    make(chan redis.Message),
		StopListenChannel: make(chan struct{}),
	}
}

func (t *Tenant) Send(channel string, msg string) error {
	log.Debug("Sending new message to channel ", channel)
	return t.rdb.Publish(ctx, channel, msg).Err()
}

func (t *Tenant) Subscribe(channel string) error {
	channelKey := fmt.Sprintf(TENANT_CHANNEL, channel)

	// Check tenant name is registered as member of channel
	if t.rdb.SIsMember(ctx, channelKey, t.Name).Val() {
		return nil
	}

	// Add member to channel
	if err := t.rdb.SAdd(ctx, channelKey, t.Name).Err(); err != nil {
		return err
	}

	return t.listen(channel)
}

func (t *Tenant) Unsubscribe(channel string) error {
	channelKey := fmt.Sprintf(TENANT_CHANNEL, channel)

	// Check tenant name, if not registered, return nil
	if !t.rdb.SIsMember(ctx, channelKey, t.Name).Val() {
		return nil
	}

	// Remove member from channel
	if err := t.rdb.SRem(ctx, channelKey, t.Name).Err(); err != nil {
		return err
	}

	t.StopListenChannel <- struct{}{}

	return nil
}

func (t *Tenant) listen(channel string) error {
	pubsub := t.rdb.Subscribe(ctx, channel)

	// The Listener
	go func() {
		log.Info("starting the listener for user:", t.Name, " on channels:", channel)
		for {
			select {
			case msg, ok := <-pubsub.Channel():
				if !ok {
					return
				}
				t.MessageChannel <- *msg

			case <-t.StopListenChannel:
				log.Info("stopping the listener for user:", t.Name)
				return
			}
		}
	}()

	return nil
}
