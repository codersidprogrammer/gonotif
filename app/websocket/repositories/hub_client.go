package repository

import (
	"context"
	"encoding/json"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	platform "github.com/codersidprogrammer/gonotif/platform/cache"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

var ctx = context.Background()

const (
	CONN_KEY = "websocket:connection:active"
	MSG_KEY  = "websocket:message"
)

type Message struct {
	ChannelName string `json:"channel"`
	From        string `json:"from"`
	To          string `json:"to"`
	Message     string `json:"message"`
}

type User struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
}

type WsClient struct {
	channel     *redis.PubSub
	wsconn      *websocket.Conn
	name        string
	isListening chan bool
	message     chan Message
	client      *User
}

func NewWsClient(topic string, name string) *WsClient {
	return &WsClient{
		channel:     platform.RedisConnection.Subscribe(ctx, topic),
		name:        name,
		isListening: make(chan bool),
		message:     make(chan Message),
		client: &User{
			Name:  name,
			Topic: topic,
		},
	}
}

func NewWsClientSpecified(topic string, name string) *WsClient {
	return &WsClient{
		channel:     platform.RedisConnection.Subscribe(ctx, topic),
		name:        name,
		isListening: make(chan bool),
		message:     make(chan Message),
		client: &User{
			Name:  name,
			Topic: topic,
		},
	}
}

func (c *WsClient) Send(message interface{}) {

}

func (c *WsClient) WebsocketHandler(conn *websocket.Conn) {

	s := conn.Query("user", "sytem")
	s2 := conn.Query("topic", "/#")

	if s != "system" && s2 != "/#" {
		c.channel = platform.RedisConnection.Subscribe(ctx, s2)
		c.client = &User{
			Name:  s,
			Topic: s2,
		}
	}

	b, err3 := json.Marshal(c.client)
	if err3 != nil {
		log.Fatal("failed to marshal", err3)
	}

	_, err2 := platform.RedisConnection.SAdd(ctx, CONN_KEY, b).Result()
	if err2 != nil {
		log.Warnf("failed to add %s to redis", CONN_KEY, err2)
		c.isListening <- false
		c.channel.Close()
		utils.ReturnErrorIfErr("", err2)
	}

	defer func() {
		c.isListening <- false
		c.channel.Close()
		conn.Close()
	}()

	// state as true
	c.isListening <- true
	c.wsconn = conn

	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = conn.ReadMessage(); err != nil {
			log.Warn("read:", err)
			break
		}
		log.Infof("recv: %s", msg)

		if mt == websocket.TextMessage {
			var _temp Message
			if errjson := json.Unmarshal(msg, &_temp); errjson != nil {
				log.Fatal("Failed to unmarshal websocket message: ", errjson)
			}
			c.message <- _temp
		} else {
			log.Infof("websocket message type: %v", mt)
		}
	}
}

func (c *WsClient) Listen() {
	for {
		select {

		// TODO: find multiple topics listeners

		// Listening for message
		case msg := <-c.message:
			log.Info("Incoming message from channel", msg.ChannelName)
			platform.RedisConnection.Publish(ctx, msg.ChannelName, msg.Message)

		// Listening for incoming subscription messages
		case msg, ok := <-c.channel.Channel():
			if !ok {
				return
			}
			log.Infof("Subs receive, channer: %s, message: %s \n", msg.Channel, msg.Payload)
			if msg.Channel == c.client.Topic {
				if err := c.wsconn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
					log.Warn("failed to write to websocket", err)
					c.isListening <- false
					c.wsconn.Close()
				}
			}

		// Stop process while client disconnect
		case listen := <-c.isListening:
			if !listen {
				b, err3 := json.Marshal(c.client)
				if err3 != nil {
					log.Fatal("failed to marshal", err3)
				}

				if err := platform.RedisConnection.SRem(ctx, CONN_KEY, b).Err(); err != nil {
					log.Warn("failed to remove listener, ", err)
				}
				log.Warn("Stopping listening")
				return
			}
		}
	}
}
