package main

import (
	"flag"
	"runtime"

	"github.com/codersidprogrammer/gonotif/app/websocket"
	"github.com/codersidprogrammer/gonotif/pkg/config"
	"github.com/codersidprogrammer/gonotif/pkg/middleware"
	platform "github.com/codersidprogrammer/gonotif/platform/cache"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func main() {

	// Bootstraping configuration
	bootstrap()

	app := fiber.New(fiber.Config{
		ServerHeader: "Notification Server",
		AppName:      "Go Notif",
	})

	// Middleware
	app.Use("/ws", middleware.UseWebsocketHandler)

	// Defining routes
	ws := websocket.NewAppWebsocket(app)
	ws.Route()

	// Run!
	log.Fatal(app.Listen(":8000"))
}

func bootstrap() {

	runtime.GOMAXPROCS(10)

	envFile := flag.String("env", ".env", "Set your dotenv file configuration (.env) as default")
	flag.Parse()

	// Load environment variables
	config.LoadEnvironment(envFile)
	platform.DoConnectRedis()
}
