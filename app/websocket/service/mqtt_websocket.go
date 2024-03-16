package service

import (
	"context"

	"github.com/codersidprogrammer/gonotif/platform/transport"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/courier-go"
)

type Payload struct {
	From string      `json:"from"`
	To   string      `json:"to"`
	Data interface{} `json:"data"`
}

type MqttMessage struct {
	Topic   string  `json:"topic"`
	Message Payload `json:"message"`
}

type mqttWebsocket struct {
	client      *courier.Client
	ctx         context.Context
	messageChan chan *MqttMessage
}

type MqttWebsocketService interface {
	Publish(topic string, message interface{}) error
	Subscribe(topic string) error
	Unsubscribe(topic string) error
	MessageChannel() chan *MqttMessage
}

func NewMqttWebsocketService() MqttWebsocketService {
	return &mqttWebsocket{
		client:      transport.MqttClient,
		ctx:         context.Background(),
		messageChan: make(chan *MqttMessage),
	}
}

// Publish implements MqttWebsocketService.
func (m *mqttWebsocket) Publish(topic string, message interface{}) error {
	if err := m.client.Publish(m.ctx, topic, message); err != nil {
		log.Error("Failed to publish, error: ", err)
		return err
	}

	return nil
}

// Subscribe implements MqttWebsocketService.
func (m *mqttWebsocket) Subscribe(topic string) error {
	if err := m.client.Subscribe(m.ctx, topic, m.message); err != nil {
		log.Error("Subscribe failed, error: ", err)
		return err
	}

	return nil
}

func (m *mqttWebsocket) Unsubscribe(topic string) error {
	if err := m.client.Unsubscribe(m.ctx, topic); err != nil {
		log.Error("Unsubscribe failed, error: ", err)
		return err
	}

	return nil
}

func (m *mqttWebsocket) MessageChannel() chan *MqttMessage {
	return m.messageChan
}

func (m *mqttWebsocket) message(ctx context.Context, ps courier.PubSub, msg *courier.Message) {
	incoming := new(Payload)
	if err := msg.DecodePayload(incoming); err != nil {
		log.Error("Failed to decode payload, error: ", err)
		return
	}

	m.messageChan <- &MqttMessage{
		Topic:   msg.Topic,
		Message: *incoming,
	}
}
