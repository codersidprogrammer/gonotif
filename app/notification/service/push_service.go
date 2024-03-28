package service

import (
	"context"

	"github.com/codersidprogrammer/gonotif/app/user/service"
	"github.com/codersidprogrammer/gonotif/platform/transport"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/courier-go"
)

type pushService struct {
	ctx                context.Context
	userSessionService service.OnlineUserService
	mqtt               *courier.Client
}

type PushService interface {
	Send(clientId string, topic string, message interface{}) error
	// Subscribe(bucketId string, topic string) error
	// SubsHandler() error
	// SendPushNotification(topic string, msg []byte) error
}

func NewPushService() PushService {
	return &pushService{
		ctx:                context.Background(),
		userSessionService: service.NewOnlineUserService(),
		mqtt:               transport.MqttClient,
	}
}

// Send implements PushService.
func (ps *pushService) Send(username string, topic string, message interface{}) error {
	uas, err := ps.userSessionService.GetAllActiveSessions()
	if err != nil {
		return err
	}

	for _, ua := range uas.Table {
		if ua.User == username {
			// then send via mqtt
			log.Infof("User: %s active. Sending via mqtt", ua.ClientID)
			if err := ps.mqtt.Publish(ps.ctx, topic, message, courier.QOSOne); err != nil {
				return err
			}

			return nil
		} else {
			// TODO: implement sender non mqtt
			// then send via another protocol
			log.Warnf("User: %s is not active. Sending via another protocol", ua.ClientID)
			return nil
		}
	}

	return nil
}
