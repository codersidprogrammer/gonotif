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

type mqttTransport struct {
	*Transport
	name string
}

func NewMqttTransport(name string) TransportService {
	return &mqttTransport{
		name: name,
		Transport: &Transport{
			lock: &sync.Mutex{},
		},
	}
}

// DoConnect implements TransportService.
func (t *mqttTransport) DoConnect() error {
	port, err := strconv.Atoi(os.Getenv("MQTT_PORT"))
	utils.ExitIfErr(err, "Failed to get mqtt port")

	host := os.Getenv("MQTT_HOST")

	if MqttClient == nil {
		t.lock.Lock()
		defer t.lock.Unlock()

		if MqttClient == nil {
			log.Infof("[MQTT] Creating New %s Transport Connection", t.name)
			MqttClient, err = courier.NewClient(
				courier.WithAddress(host, uint16(port)),
				courier.WithClientID("xops_system"),
				courier.WithUsername(os.Getenv("MQTT_USERNAME")),
				courier.WithPassword(os.Getenv("MQTT_PASSWORD")),
				courier.WithKeepAlive(60),
				// courier.WithMaxReconnectInterval(30),
			)

			if err != nil {
				return err
			}

			if err := MqttClient.Start(); err != nil {
				log.Fatalf("[MQTT] failed to start %s client: %v", t.name, err)
			}

		} else {
			log.Info("[MQTT] Using existing MQTT connection")
		}
	}

	log.Infof("[MQTT] %s connection established", t.name)
	return nil
}

func (t *mqttTransport) Close() error {
	log.Info("Closing transport connection: ", t.name)
	MqttClient.Stop()
	return nil
}
