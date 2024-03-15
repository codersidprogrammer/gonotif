package websocket

import (
	"github.com/codersidprogrammer/gonotif/app/websocket/controller"
	repository "github.com/codersidprogrammer/gonotif/app/websocket/repositories"
	fiber_websocket "github.com/gofiber/contrib/websocket"
)

func (a *Websocket) Route() {
	// go a.WsRegister()

	wc := repository.NewWsClient()
	go wc.Listen()

	dc := controller.NewDummy()
	go dc.ConnListener()
	go dc.MessageListener()

	route := a.app.Group("/ws")
	route.Get("/test", fiber_websocket.New(dc.WebsocketHandler))
	route.Get("/", fiber_websocket.New(wc.WebsocketHandler))
}
