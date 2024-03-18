package websocket

import (
	"github.com/codersidprogrammer/gonotif/app/websocket/controller"
	"github.com/codersidprogrammer/gonotif/pkg/model"
	fiber_websocket "github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
)

type appWebsocket struct {
	app                *fiber.App
	mqttNotifControler controller.WebsocketController
}

func NewAppWebsocket(app *fiber.App) model.App {
	return &appWebsocket{
		app:                app,
		mqttNotifControler: controller.NewMqttNotificationController(),
	}
}

// Route implements model.App.
func (a *appWebsocket) Route() {
	mnc := controller.NewMqttNotificationController()
	go mnc.ConnectionListener()
	go mnc.MessageListener()

	route := a.app.Group("/ws")
	route.Get("/notification", fiber_websocket.New(mnc.WebsocketHandler))
	route.Get("/stream", fiber_websocket.New(mnc.WebsocketHandler))
}
