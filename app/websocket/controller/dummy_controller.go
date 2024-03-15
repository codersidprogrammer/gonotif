package controller

import (
	repository "github.com/codersidprogrammer/gonotif/app/websocket/repositories"
	platform "github.com/codersidprogrammer/gonotif/platform/cache"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

type DummyController struct {
	register   chan *websocket.Conn
	clients    map[*websocket.Conn]*repository.Tenant
	unregister chan *websocket.Conn
	channel    string
	tenant     *repository.Tenant
}

func NewDummy() *DummyController {
	return &DummyController{
		register:   make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]*repository.Tenant),
		unregister: make(chan *websocket.Conn),
		tenant:     &repository.Tenant{},
	}
}

func (d *DummyController) ConnListener() {
	log.Debug("starting connection handler")
	for {
		select {
		case reg := <-d.register:
			d.clients[reg] = d.tenant
			log.Infof("Registered %s", len(d.clients))

		case unreg := <-d.unregister:
			d.tenant.Unsubscribe(d.channel)
			delete(d.clients, unreg)
		}
	}
}

// move this to internal repository by make function with chan params
func (d *DummyController) MessageListener() {
	log.Debug("starting message handler")
	for {
		select {
		case msg := <-d.tenant.MessageChannel:
			log.Infof("Receving msg from redis on channel: %s", msg.Channel)
			// for cn, _ := range clients {
			// 	if cn.Query("channel") == msg.Channel {
			// 		log.Info("Message is: ", msg.Payload)
			// 	}
			// }
		}
	}
}

func (d *DummyController) WebsocketHandler(c *websocket.Conn) {

	channel := c.Query("channel", "xops_notification")
	name := c.Query("name", "system")
	d.channel = channel

	t := repository.NewTenant(name, platform.RedisConnection)
	d.tenant = t

	t.Subscribe(channel)
	d.register <- c // This section caused lock bugs

	defer func() {
		t.Unsubscribe(channel)
		d.unregister <- c // This section caused lock bugs
	}()

	var (
		mt  int
		msg []byte
		err error
	)

loop:
	for {
		if mt, msg, err = c.ReadMessage(); err != nil {
			log.Warn("read:", err)
			break loop
		}
		log.Infof("recv: %s", msg)

		if mt == websocket.TextMessage {
			t.Send(channel, string(msg))
			log.Infof("websocket message: %v", msg)
		} else {
			log.Infof("websocket message type: %v", mt)
		}
	}
}
