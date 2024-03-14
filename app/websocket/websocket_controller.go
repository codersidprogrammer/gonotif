package websocket

import (
	"context"

	"github.com/codersidprogrammer/gonotif/pkg/model"
	platform "github.com/codersidprogrammer/gonotif/platform/cache"
	fiber_websocket "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var ctx = context.Background()

func NewAppWebsocket(app *fiber.App) model.App {
	return &Websocket{
		app:        app,
		clients:    make(map[*fiber_websocket.Conn]client),
		register:   make(chan *fiber_websocket.Conn),
		unregister: make(chan *fiber_websocket.Conn),
		broadcast:  make(chan []byte),
	}
}

// GetWebsocketHandler implements AppWebsocket.
func (a *Websocket) GetWebsocketHandler(c *fiber_websocket.Conn) {

	// platform.RedisConnection.SAdd(ctx, )
	// defer func() {
	// 	a.unregister <- c
	// 	c.Close()
	// }()

	// a.register <- c

	// ps := platform.RedisConnection.Subscribe(ctx, "/xops/notification")
	// defer ps.Close()

	var (
		mt  int
		msg []byte
		err error
	)

loop:
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Warn("read:", err)
			break loop
		}

		if mt == fiber_websocket.TextMessage {
			c.WriteMessage(mt, msg)
			log.Infof("websocket message: %v", msg)
		} else {
			log.Infof("websocket message type: %v", mt)
		}
	}
}

func (a *Websocket) WsRegister() {
	go redisHandler(a)

	for {
		select {
		case connection := <-a.register:
			a.clients[connection] = client{
				channel: connection.Query("channel"),
			}
			log.Info(a.clients)

		case message := <-a.broadcast:
			for connection := range a.clients {
				if connection.Query("channel") == "xops/notification" {
					if err := connection.WriteMessage(fiber_websocket.TextMessage, message); err != nil {
						log.Warn(err)

						a.unregister <- connection
						log.Info("Broadcasting message to client")
						connection.WriteMessage(fiber_websocket.CloseMessage, message)
						connection.Close()
					}
				}
			}

		case connection := <-a.unregister:
			delete(a.clients, connection)
		}
	}
}

func redisHandler(a *Websocket) {
	ps := platform.RedisConnection.Subscribe(ctx, "/xops/notification")

	for {
		select {
		case msg, ok := <-ps.Channel():
			if !ok {
				break
			}
			log.Infof("Subs receive: %s", msg.Payload)
			a.broadcast <- []byte(msg.Payload)
		}
	}
}
