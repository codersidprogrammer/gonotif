package controller

import (
	"time"

	repository "github.com/codersidprogrammer/gonotif/app/notification/repositories"
	"github.com/codersidprogrammer/gonotif/app/notification/service"
	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type controller struct {
	service service.NotificationService
}

type NotificationController interface {
	GetNotification(c *fiber.Ctx) error
	GetNotifications(c *fiber.Ctx) error
	CreateNotification(c *fiber.Ctx) error
}

func NewNotificationController(service service.NotificationService) NotificationController {
	return &controller{
		service: service,
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
