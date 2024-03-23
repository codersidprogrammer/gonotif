package controller

import (
	"time"

	"github.com/codersidprogrammer/gonotif/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
)

type controller struct{}

type UserController interface {
	OnUserHookHandler(c *fiber.Ctx) error
}

func NewUserAppController() UserController {
	return &controller{}
}

// OnUserHookHandler implements UserController.
func (*controller) OnUserHookHandler(c *fiber.Ctx) error {
	var body interface{}
	if err := c.BodyParser(&body); err != nil {
		utils.ReturnErrMessageIfErr(err, "onUserHookHandler", c)
	}

	log.Debug(body)
	return c.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": nil,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}
