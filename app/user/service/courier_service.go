package service

import (
	"sync"

	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/courier-go"
)

var MqttClient *courier.Client
var lock = sync.Mutex{}

func InitConnectionMqtt() error {
	if MqttClient == nil {
		lock.Lock()
		defer lock.Unlock()

		if MqttClient == nil {
			var err error
			MqttClient, err = courier.NewClient(
				courier.WithAddress("172.16.41.73", 1883),
				courier.WithClientID("xops_system"),
				courier.WithUsername("sysadmin"),
				courier.WithPassword("sysadmin"),
			)

			if err != nil {
				return err
			}

			if err := MqttClient.Start(); err != nil {
				log.Fatalf("[MQTT] failed to start %s client: %v", "name", err)
			}

		} else {
			log.Info("[MQTT] Using existing MQTT connection")
		}
	}

	log.Info("[MQTT] Connection VerneMQ established")
	return nil
}
