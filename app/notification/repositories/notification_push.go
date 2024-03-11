package repository

import "github.com/go-redis/redis/v8"

type PushClient struct {
	Name                 string
	TopicHandler         *redis.PubSub
	MessageChannel       chan redis.Message
	IsListening          bool
	StopListeningChannel chan struct{}
}

type PushMessage struct {
	Topic   string `json:"topic"`
	Message string `json:"message"`
}
