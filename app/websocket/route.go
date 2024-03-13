package websocket

import (
	repository "github.com/codersidprogrammer/gonotif/app/websocket/repositories"
	fiber_websocket "github.com/gofiber/contrib/websocket"
)

func (a *Websocket) Route() {
	go a.WsRegister()

	wc := repository.NewWsClient("/xops/personal", "582276")
	go wc.Listen()

	route := a.app.Group("/ws")
	// route.Get("/:id", fiber_websocket.New(a.GetWebsocketHandler))
	route.Get("/", fiber_websocket.New(wc.WebsocketHandler))
}
