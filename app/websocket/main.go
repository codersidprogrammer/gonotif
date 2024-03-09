package websocket

import (
	fiber_websocket "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type client struct {
	channel string
}

type Websocket struct {
	app        *fiber.App
	clients    map[*fiber_websocket.Conn]client
	register   chan *fiber_websocket.Conn
	unregister chan *fiber_websocket.Conn
	broadcast  chan []byte
}

type AppWebsocket interface {
	Route()
	GetWebsocketHandler(c *fiber_websocket.Conn)
}
