package transport

import (
	"os"
	"strconv"
	"sync"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/courier-go"
)

var MqttClient *courier.Client

var lock = &sync.Mutex{}

func DoMqttConnect() (*courier.Client, error) {
	port, err := strconv.Atoi(os.Getenv("MQTT_PORT"))
	utils.ExitIfErr(err, "Failed to get mqtt port")

	host := os.Getenv("MQTT_HOST")

	if MqttClient == nil {
		lock.Lock()
		defer lock.Unlock()

		if MqttClient == nil {
			log.Info("Creating new MQTT connection")
			MqttClient, err = courier.NewClient(
				courier.WithAddress(host, uint16(port)),
			)

			if err != nil {
				return nil, err
			}

			if err := MqttClient.Start(); err != nil {
				log.Fatal("failed to start MQTT client: ", err)
			}

		} else {
			log.Info("Using existing MQTT connection")
		}
	}

	log.Info("MQTT conneection established")
	return MqttClient, nil
}
