package websocket

import (
	fiber_websocket "github.com/gofiber/contrib/websocket"
)

func (a *Websocket) Route() {
	go a.WsRegister()

	route := a.app.Group("/ws")

	route.Get("/:id", fiber_websocket.New(a.GetWebsocketHandler))
}
