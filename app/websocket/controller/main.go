package controller

import "github.com/gofiber/contrib/websocket"

type User struct {
	Name  string `json:"name"`
	Topic string `json:"topic"`
}

type WebsocketController interface {
	// This method is called by go routines
	// This will handle registration and unregistration
	// while new connection is established
	ConnectionListener()

	// This method is called by go routines
	// This will handle incoming messages that received
	// from clients (MQTT)
	MessageListener()

	// This method is called by fiber websocket library
	// all logic functions are called here
	WebsocketHandler(c *websocket.Conn)

	// This method is called by defer handler
	// This will close all connections
	Close(c *websocket.Conn)
}
