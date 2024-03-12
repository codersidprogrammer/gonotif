package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"

	"github.com/codersidprogrammer/gonotif/app/notification"
	"github.com/codersidprogrammer/gonotif/app/websocket"
	"github.com/codersidprogrammer/gonotif/pkg/config"
	"github.com/codersidprogrammer/gonotif/pkg/middleware"
	"github.com/codersidprogrammer/gonotif/pkg/queue"
	"github.com/codersidprogrammer/gonotif/platform/cache"
	"github.com/codersidprogrammer/gonotif/platform/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {

	// Bootstraping configuration
	bootstrap()

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
	ws := websocket.NewAppWebsocket(app)
	ws.Route()

	an := notification.NewAppNotification(app)
	an.Route()

	// Run!
	runHandler(app)
}

func bootstrap() {

	runtime.GOMAXPROCS(10)

	envFile := flag.String("env", ".env", "Set your dotenv file configuration (.env) as default")
	flag.Parse()

	// Load environment variables
	config.LoadEnvironment(envFile)
	cache.DoConnectRedis()
	database.DoConnectEtcd()
	queue.QueueHandler()
}

func runHandler(a *fiber.App) {
	_idleConn := make(chan struct{})

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		log.Warn("Shutting down server...")
		if err := a.Shutdown(); err != nil {
			log.Warnf("Failed to shutdown: %v", err)
		}

		log.Warn("Server was landed safely...")
		close(_idleConn)
	}()

	if err := a.Listen(":8000"); err != nil {
		log.Errorf("Server not running!! Error: %v", err)
	}

	<-_idleConn
}
