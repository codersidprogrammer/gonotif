package controller

import (
	"github.com/codersidprogrammer/gonotif/app/websocket/service"
	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

type notificationConnection struct {
	register   chan *websocket.Conn
	clients    map[*websocket.Conn]*User
	unregister chan *websocket.Conn
	user       *User
	service    service.MqttWebsocketService
}

func NewMqttNotificationController() WebsocketController {
	return &notificationConnection{
		register:   make(chan *websocket.Conn),
		clients:    make(map[*websocket.Conn]*User),
		unregister: make(chan *websocket.Conn),
		service:    service.NewMqttWebsocketService(),
	}
}

// Close implements WebsocketController.
func (n *notificationConnection) Close(c *websocket.Conn) {
	n.unregister <- c
	c.Close()
}

// ConnectionListener implements WebsocketController.
func (n *notificationConnection) ConnectionListener() {
	log.Debug("starting connection handler")
	for {
		select {
		case reg := <-n.register:
			n.clients[reg] = n.user
			log.Infof("Registered, total active connections: %d", len(n.clients))

		case unreg := <-n.unregister:
			log.Infof("Unregistered, total active connections: %d", len(n.clients))
			delete(n.clients, unreg)
		}
	}
}

// MessageListener implements WebsocketController.
func (n *notificationConnection) MessageListener() {
	for {
		select {

		// Wait for incoming messages and send them to specific user
		case msg := <-n.service.MessageChannel():
			for conn := range n.clients {
				if conn.Query("name") == msg.Message.To && conn.Query("channel") == msg.Topic {
					if err := conn.WriteJSON(msg); err != nil {
						log.Error("WriteMessage error: ", err)
						n.Close(conn)
						return
					}
				}

				if utils.CheckIfHasSpecifiedSuffix(msg.Message.To, "/", "all") {
					if err := conn.WriteJSON(msg); err != nil {
						log.Error("WriteMessage error: ", err)
						n.Close(conn)
						return
					}
				}
			}

		// If connection is closed, then unsubscribe from
		// current topic
		case conn := <-n.unregister:
			n.service.Unsubscribe(conn.Query("channel"))
		}
	}
}

// WebsocketHandler implements WebsocketController.
func (n *notificationConnection) WebsocketHandler(c *websocket.Conn) {
	channel := c.Query("channel", "/#")
	name := c.Query("name", "system")

	n.user = &User{
		Name:  name,
		Topic: channel,
	}

	// Register handler
	n.register <- c
	defer n.Close(c)

	// subscribe handler
	if err := n.service.Subscribe(channel); err != nil {
		log.Error("Unable to subscribe error: ", err)
		n.Close(c)
	}

	// Main loop for websocket connection
	var (
		mt  int
		err error
	)

loop:
	for {
		if mt, _, err = c.ReadMessage(); err != nil {
			log.Warn("read message message: ", err)
			break loop
		}

		if mt == websocket.TextMessage {
			_defaultMessage := `
				Your message was not published. Because this service only
				for one way message`

			if err := c.WriteMessage(websocket.TextMessage, []byte(_defaultMessage)); err != nil {
				log.Warn("write message error: ", err)
				n.Close(c)
				break loop
			}

		} else {
			log.Infof("websocket message type: %v", mt)
		}
	}
}
