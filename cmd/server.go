package cmd

import (
	"flag"
	"os"
	"os/signal"
	"runtime"

	"github.com/codersidprogrammer/gonotif/app/notification"
	"github.com/codersidprogrammer/gonotif/app/user"
	"github.com/codersidprogrammer/gonotif/app/websocket"
	"github.com/codersidprogrammer/gonotif/pkg/config"
	"github.com/codersidprogrammer/gonotif/pkg/queue"
	"github.com/codersidprogrammer/gonotif/pkg/routes"
	"github.com/codersidprogrammer/gonotif/platform/database"
	"github.com/codersidprogrammer/gonotif/platform/transport"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

var sigint = make(chan os.Signal, 1)
var mqtt = transport.NewMqttTransport("MQTT Verne")
var redis = transport.NewRedisTransport("redis")
var http = transport.NewHttpClientTransport("HTTP Service")

func Bootstrap() {

	runtime.GOMAXPROCS(10)

	envFile := flag.String("env", ".env", "Set your dotenv file configuration (.env) as default")
	flag.Parse()

	// Load environment variables
	config.LoadEnvironment(envFile)

	// load utilities connection
	database.DoConnectEtcd()

	// Do connect all transports
	mqtt.DoConnect()
	redis.DoConnect()
	http.DoConnect()
	// service.InitConnectionMqtt() // TODO: remove soon

	// Start queue handler
	queue.QueueHandler()
}

func Route(app *fiber.App) {

	// Registering each apps here
	pr := routes.NewRoute("public")
	pr.Register(notification.NewAppNotification(app))
	pr.Register(websocket.NewAppWebsocket(app))
	pr.Register(user.NewUserApp(app))

	// Load all Route() methods
	pr.Handler()
}

func RunServer(a *fiber.App) {
	_idleConn := make(chan struct{})

	go func() {
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// Closing transport connection
		mqtt.Close()
		redis.Close()

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
