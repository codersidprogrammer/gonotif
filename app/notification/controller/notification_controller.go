package controller

import (
	"time"

	"github.com/codersidprogrammer/gonotif/app/notification/dto"
	repository "github.com/codersidprogrammer/gonotif/app/notification/repositories"
	"github.com/codersidprogrammer/gonotif/app/notification/service"
	"github.com/codersidprogrammer/gonotif/pkg/queue"
	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gojek/work"
)

type controller struct {
	service     service.NotificationBucketService
	pushService service.NotificationPushService
	queue       *queue.Queue
}

type NotificationController interface {
	GetNotification(c *fiber.Ctx) error
	GetNotifications(c *fiber.Ctx) error
	CreateNotification(c *fiber.Ctx) error
	CreateNotificationWorker(c *fiber.Ctx) error

	Publish(c *fiber.Ctx) error
	Subscribe(c *fiber.Ctx) error
}

func NewNotificationController(service service.NotificationBucketService, push service.NotificationPushService) NotificationController {
	return &controller{
		service:     service,
		pushService: push,
		queue:       queue.NewQueue("development_test"),
	}
}

// CreateNotification implements NotificationController.
func (co *controller) CreateNotification(c *fiber.Ctx) error {
	var nb repository.NotificationBucket
	err := c.BodyParser(&nb)
	utils.ReturnHttpErr400MessageIfErr(err, "unmarshal notification", c)

	result, err := co.service.Create(nb)
	utils.ReturnErrMessageIfErr(err, "unmarshal notification", c)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": result,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}

func (co *controller) CreateNotificationWorker(c *fiber.Ctx) error {

	var data dto.CreateNotificationRequest
	err := c.BodyParser(&data)
	utils.ReturnHttpErr400MessageIfErr(err, "unmarshal notification", c)

	switch data.Type {

	// Default handler if no type specified
	default:
		return c.Status(fiber.StatusNotFound).JSON(&fiber.Map{
			"data": "No process found",
			"meta": &fiber.Map{},
			"time": time.Now(),
		})

	case dto.Push:
		log.Debug(data)
		result, err := co.queue.Register("send_notifcation", work.Q{
			"topic":   "/xops/notification",
			"payload": data.Payload,
		})
		utils.ReturnHttpErr400MessageIfErr(err, "queue failed", c)

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"data": data.Payload,
			"meta": &fiber.Map{
				"ID": result,
			},
			"time": time.Now(),
		})

	case dto.Email:
		log.Debug(data)
		result, err := co.queue.Register("send_email", work.Q{
			"topic":   "send mail",
			"payload": data.Payload,
		})
		utils.ReturnHttpErr400MessageIfErr(err, "queue failed", c)

		return c.Status(fiber.StatusOK).JSON(&fiber.Map{
			"data": data.Payload,
			"meta": &fiber.Map{
				"ID": result,
			},
			"time": time.Now(),
		})
	}
}

func (co *controller) Publish(c *fiber.Ctx) error {
	var pm repository.PushMessage
	err := c.BodyParser(&pm)
	utils.ReturnHttpErr400MessageIfErr(err, "unmarshal notification", c)

	co.pushService.SendPushNotification(pm.Topic, []byte(pm.Message))

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": pm,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}

func (co *controller) Subscribe(c *fiber.Ctx) error {
	id := c.Params("id")
	topic := c.Query("topic")

	co.pushService.Subscribe(id, topic)

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": nil,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}

// GetNotification implements NotificationController.
func (co *controller) GetNotification(c *fiber.Ctx) error {
	nb, err := co.service.Get(c.Params("id"))
	utils.ReturnErrMessageIfErr(err, "GetNotification", c)
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": nb,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}

func (co *controller) GetNotifications(c *fiber.Ctx) error {

	var nb []repository.NotificationBucket
	var err error

	if c.Query("prefix") == "" {
		nb, err = co.service.GetAll()
		utils.ReturnErrMessageIfErr(err, "Get all notifications", c)
	} else {
		nb, err = co.service.GetAllWithPrefix(c.Query("prefix"))
		utils.ReturnErrMessageIfErr(err, "Get all notifications prefixed", c)
	}

	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": nb,
		"meta": &fiber.Map{
			"prefix": c.Query("prefix"),
		},
		"time": time.Now(),
	})
}
