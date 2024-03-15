package repository

import (
	"context"
	"encoding/json"

	// "github.com/codersidprogrammer/gonotif/pkg/utils"
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
	clients     map[*websocket.Conn]User
	register    chan *websocket.Conn
	unregister  chan *websocket.Conn
}

func NewWsClient() *WsClient {
	return &WsClient{
		channel:     platform.RedisConnection.Subscribe(ctx, "/#"),
		isListening: make(chan bool),
		message:     make(chan Message),
		register:    make(chan *websocket.Conn),
		unregister:  make(chan *websocket.Conn),
		clients:     make(map[*websocket.Conn]User),
		client:      &User{},
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

func (c *WsClient) Break() error {
	log.Info("Break called")

	c.isListening <- false
	c.unregister <- c.wsconn

	if err := c.channel.Unsubscribe(ctx, c.client.Topic); err != nil {
		log.Fatal("Unsubscribe error: ", err)
	}

	// This lead to locking on the object while closing
	// if err := c.channel.Close(); err != nil {
	// 	log.Fatal("Fatal closing channel: ", err)
	// }

	if err := c.wsconn.Close(); err != nil {
		log.Fatal("Fatal closing wsconn: ", err)
	}

	return nil
}

func (c *WsClient) Send(message interface{}) error {
	return platform.RedisConnection.Publish(ctx, c.client.Topic, message).Err()
}

func (c *WsClient) WebsocketHandler(conn *websocket.Conn) {

	user := conn.Query("user", "sytem")
	topic := conn.Query("topic", "/#")

	c.name = user
	c.channel = platform.RedisConnection.Subscribe(ctx, topic)
	c.client = &User{
		Name:  user,
		Topic: topic,
	}

	defer c.Break()

	// // state as true
	c.isListening <- true
	c.wsconn = conn
	c.register <- conn

	var (
		mt  int
		msg []byte
		err error
	)

loop:
	for {
		if mt, msg, err = conn.ReadMessage(); err != nil {
			log.Warn("read:", err)
			break loop
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

		case cn := <-c.register:
			c.clients[cn] = *c.client
			log.Warnf("Clients: %v ", len(c.clients))

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

			for conn, usr := range c.clients {
				if msg.Channel == usr.Topic {
					if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
						log.Warn("failed to write to websocket", err)
						c.isListening <- false
						conn.Close()
					}
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
				// return
			}

		case ur := <-c.unregister:
			delete(c.clients, ur)
			log.Info("Unregistered, Total Clients registered: ", len(c.clients))
		}
	}
}
