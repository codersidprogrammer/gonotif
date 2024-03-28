package controller

import (
	"time"

	"github.com/codersidprogrammer/gonotif/app/notification/dto"
	"github.com/codersidprogrammer/gonotif/app/notification/service"
	"github.com/codersidprogrammer/gonotif/pkg/queue"
	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gojek/work"
)

type controller struct {
	pusher service.PushService
	queue  *queue.Queue
}

type NotificationController interface {
	CreatePushNotification(c *fiber.Ctx) error
}

func NewNotificationController() NotificationController {
	return &controller{
		queue:  queue.NewQueue("development_test"),
		pusher: service.NewPushService(),
	}
}

func (co *controller) CreatePushNotification(c *fiber.Ctx) error {
	var _dto dto.CreatePushNotificationRequest
	err := c.BodyParser(&_dto)
	utils.ReturnHttpErr400MessageIfErr(err, "Unmarshall error", c)

	// Register queue
	result, err := co.queue.Register("send_notification", work.Q{
		"topic":    _dto.Topic,
		"username": _dto.Username,
		"payload":  _dto.Payload,
	})
	utils.ReturnErrMessageIfErr(err, "Failed on queue, error ", c)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": _dto.Payload,
		"meta": &fiber.Map{
			"key": result,
		},
		"time": time.Now(),
	})
}
