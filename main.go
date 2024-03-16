package main

import (
	"os"

	"github.com/codersidprogrammer/gonotif/cmd"
	"github.com/codersidprogrammer/gonotif/pkg/middleware"
	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/pprof"
)

func main() {

	// Bootstraping configuration
	cmd.Bootstrap()

	// Configure fiber apps
	app := fiber.New(fiber.Config{
		ServerHeader:      "Notification Server",
		AppName:           "Go Notif",
		EnablePrintRoutes: true,
		Prefork:           false,
	})

	// Middleware
	f, err := os.OpenFile("./request.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	utils.ExitIfErr(err, "Couldn't open log file")
	app.Use(logger.New(logger.Config{
		Output: f,
	}))
	app.Use("/ws", middleware.UseWebsocketHandler)
	app.Use(pprof.New(pprof.Config{
		Prefix: "/profile",
	}))

	// Defining routes
	cmd.Route(app)

	// Run!
	cmd.RunServer(app)
}
