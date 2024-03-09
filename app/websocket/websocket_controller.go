package websocket

import (
	"context"

	platform "github.com/codersidprogrammer/gonotif/platform/cache"
	fiber_websocket "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var ctx = context.Background()

func NewAppWebsocket(app *fiber.App) AppWebsocket {
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
	defer func() {
		a.unregister <- c
		c.Close()
	}()

	a.register <- c

	var (
		mt  int
		msg []byte
		err error
	)
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Warn("read:", err)
			break
		}
		log.Infof("recv: %s", msg)

		if mt == fiber_websocket.TextMessage {
			platform.RedisConnection.Publish(ctx, "publish-test", msg)
			a.broadcast <- msg
		} else {
			log.Infof("websocket message type: %v", mt)
		}

		// if err = c.WriteJSON(msg); err != nil {
		// 	log.Info("write:", err)
		// 	break
		// }
	}
}

func (a *Websocket) WsRegister() {
	go redisHandler()

	for {
		select {
		case connection := <-a.register:
			a.clients[connection] = client{
				channel: connection.Query("channel"),
			}

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

func redisHandler() {
	ps := platform.RedisConnection.Subscribe(ctx, "publish-test")

	for {
		select {
		case msg, ok := <-ps.Channel():
			if !ok {
				break
			}
			log.Infof("Subs receive: %s", msg.Payload)

		}
	}
}
