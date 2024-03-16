package cmd

import (
	"flag"
	"os"
	"os/signal"
	"runtime"

	"github.com/codersidprogrammer/gonotif/app/notification"
	"github.com/codersidprogrammer/gonotif/app/websocket"
	"github.com/codersidprogrammer/gonotif/pkg/config"
	"github.com/codersidprogrammer/gonotif/pkg/queue"
	"github.com/codersidprogrammer/gonotif/pkg/routes"
	"github.com/codersidprogrammer/gonotif/platform/cache"
	"github.com/codersidprogrammer/gonotif/platform/database"
	"github.com/codersidprogrammer/gonotif/platform/transport"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

func Bootstrap() {

	runtime.GOMAXPROCS(10)

	envFile := flag.String("env", ".env", "Set your dotenv file configuration (.env) as default")
	flag.Parse()

	// Load environment variables
	config.LoadEnvironment(envFile)

	// load utilities connection
	cache.DoConnectRedis()
	database.DoConnectEtcd()
	queue.QueueHandler()
	transport.DoMqttConnect()
}

func Route(app *fiber.App) {

	// Registering each apps here
	pr := routes.NewRoute("public")
	pr.Register(notification.NewAppNotification(app))
	pr.Register(websocket.NewAppWebsocket(app))

	// Load all Route() methods
	pr.Handler()
}

func RunServer(a *fiber.App) {
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
