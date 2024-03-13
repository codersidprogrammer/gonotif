package websocket

type Message struct {
	ChannelName string `json:"channel"`
	From        string `json:"from"`
	To          string `json:"to"`
	Message     []byte `json:"message"`
}
