package transport

import "sync"

type Transport struct {
	lock *sync.Mutex
}

type TransportService interface {
	// Connect transport provider
	// such as `redis, mqtt, etc`.
	// This function should be called on bootstrap
	// And it should be singleton
	DoConnect() error

	// Close transport provider
	// such as `redis, mqtt, etc`.
	Close() error
}
