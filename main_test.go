package main

import (
	"context"
	"os"
	"os/signal"
	"testing"

	"github.com/codersidprogrammer/gonotif/cmd"
	"github.com/codersidprogrammer/gonotif/platform/transport"
	"github.com/gojek/courier-go"
)

func TestMain(t *testing.T) {
	cmd.Bootstrap()
	ts := transport.NewMqttTransport("MQTT Test")
	if err := ts.DoConnect(); err != nil {
		t.Fatal("Failed to connect to server. err: ", err)
	}

	if transport.MqttClient == nil {
		t.Fatal("MqttClient is nil")
	}

	type chatMessage struct {
		From string      `json:"from"`
		To   string      `json:"to"`
		Data interface{} `json:"data"`
	}

	_msg := &chatMessage{
		From: "test-username-1",
		To:   "test-username-2",
		Data: map[string]string{
			"message": "Hi, User 2!",
		},
	}

	_ = transport.MqttClient.Subscribe(context.Background(), "c/user/a", func(ctx context.Context, ps courier.PubSub, msg *courier.Message) {
		var incoming interface{}
		if err := msg.DecodePayload(incoming); err != nil {
			t.Fatal("Failed to decode payload, error: ", err)
		}

		t.Log(incoming)
		transport.MqttClient.Publish(context.Background(), "a/receive/c", _msg, courier.QOSOne)
	}, courier.QOSOne)

	var sigint = make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}
