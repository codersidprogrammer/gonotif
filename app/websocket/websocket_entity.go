package websocket

type Message struct {
	channel string
	userId  string
	data    []byte
}
