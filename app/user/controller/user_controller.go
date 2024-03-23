package controller

import (
	"context"
	"time"

	"github.com/codersidprogrammer/gonotif/app/user/service"
	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/courier-go"
)

type controller struct{}

type UserController interface {
	OnUserHookHandler(c *fiber.Ctx) error
}

func NewUserAppController() UserController {
	service.InitConnectionMqtt()
	return &controller{}
}

// OnUserHookHandler implements UserController.
func (*controller) OnUserHookHandler(c *fiber.Ctx) error {
	var body interface{}
	if err := c.BodyParser(&body); err != nil {
		utils.ReturnErrMessageIfErr(err, "onUserHookHandler", c)
	}

	log.Debug(body)

	if err := service.MqttClient.Publish(context.Background(), "xops/api/user", body, courier.QOSTwo); err != nil {
		log.Error("Error publishing, error ", err)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": nil,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}
