package main

import (
	"github.com/codersidprogrammer/gonotif/cmd"
	"github.com/codersidprogrammer/gonotif/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Bootstraping configuration
	cmd.Bootstrap()

	// Configure fiber apps
	app := fiber.New(fiber.Config{
		ServerHeader:      "Notification Server",
		AppName:           "Go Notif",
		EnablePrintRoutes: false,
		Prefork:           false,
	})

	// Middleware
	app.Use(logger.New())
	app.Use("/ws", middleware.UseWebsocketHandler)

	// Defining routes
	cmd.Route(app)

	// Run!
	cmd.RunServer(app)
}
