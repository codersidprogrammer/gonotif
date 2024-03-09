package main

import (
	"github.com/codersidprogrammer/gonotif/app/websocket"
	"github.com/codersidprogrammer/gonotif/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {
	app := fiber.New(fiber.Config{
		ServerHeader: "Notification Server",
		AppName:      "Go Notif",
	})

	// Middleware
	app.Use("/ws", middleware.UseWebsocketHandler)

	//
	ws := websocket.NewAppWebsocket(app)
	ws.Route()

	// Run!
	log.Fatal(app.Listen(":8000"))
}
