package controller

import (
	"encoding/json"
	"sync"

	"github.com/codersidprogrammer/gonotif/app/websocket/service"
	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2/log"
)

type notificationConnection struct {
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	message    chan []byte
	clients    map[*websocket.Conn]*User
	user       *User
	service    service.MqttWebsocketService
	lock       *sync.RWMutex
}

func NewMqttNotificationController() WebsocketController {
	return &notificationConnection{
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
		message:    make(chan []byte),
		clients:    make(map[*websocket.Conn]*User),
		service:    service.NewMqttWebsocketService(),
		lock:       &sync.RWMutex{},
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

		// Add current connection to pool
		case reg := <-n.register:
			n.lock.Lock()
			n.clients[reg] = n.user
			log.Infof("Registered, total active connections: %d", len(n.clients))
			n.lock.Unlock()

		// Remove current connection from pool
		case unreg := <-n.unregister:
			n.lock.Lock()
			log.Infof("Unregistered, total active connections: %d", len(n.clients))
			delete(n.clients, unreg)
			n.lock.Unlock()

		// Wait for incoming messages and send them to websocket
		// Using channel to make sure message are received one by one
		case msg := <-n.message:
			var _msg service.MqttMessage
			if err := json.Unmarshal(msg, &_msg); err != nil {
				log.Error("Error unmarshaling message, error ", err)
				break
			}

			for conn := range n.clients {
				if conn.Query("name") == _msg.Message.To && conn.Query("channel") == _msg.Topic {
					if err := conn.WriteJSON(_msg); err != nil {
						log.Error("WriteMessage error: ", err)
						n.Close(conn)
						return
					}
				}

				if utils.CheckIfHasSpecifiedSuffix(_msg.Message.To, "/", "all") {
					if err := conn.WriteJSON(_msg); err != nil {
						log.Error("WriteMessage error: ", err)
						n.Close(conn)
						return
					}
				}
			}
		}
	}
}

// MessageListener implements WebsocketController.
func (n *notificationConnection) MessageListener() {
	for {
		select {

		// Wait for incoming messages and send them to specific user
		case msg := <-n.service.MessageChannel():
			msgByte, err := json.Marshal(msg)
			if err != nil {
				log.Error("Error marshaling message, error ", err)
			}
			n.message <- msgByte

			// If connection is closed, then unsubscribe from
			// current topic
			// case conn := <-n.unregister:
			// 	n.service.Unsubscribe(conn.Query("channel"))
		}
	}
}

// WebsocketHandler implements WebsocketController.
func (n *notificationConnection) WebsocketHandler(c *websocket.Conn) {
	channel := c.Query("channel", "/#")
	name := c.Query("name", "system")

	// TODO: add necessary field for utilization
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
			n.Close(c)
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
