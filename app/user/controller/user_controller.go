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

type controller struct {
	activeUserService service.OnlineUserService
}

type Headers struct {
	VerneHook string `reqHeader:"Vernemq-Hook"`
}

type UserController interface {
	CreateActiveUser(ctx *fiber.Ctx) error
	GetActiveUsers(ctx *fiber.Ctx) error
	OnUserHookHandler(ctx *fiber.Ctx) error
}

func NewUserAppController() UserController {
	// service.InitConnectionMqtt()
	return &controller{
		activeUserService: service.NewOnlineUserService(),
	}
}

// CreateActiveUser implements UserController.
func (c *controller) CreateActiveUser(ctx *fiber.Ctx) error {
	var user service.ActiveUser
	if err := ctx.BodyParser(&user); err != nil {
		utils.ReturnErrMessageIfErr(err, "Creating user failed, error", ctx)
	}

	au, err := c.activeUserService.SetOnlineUser(&user)
	utils.ReturnErrMessageIfErr(err, "Failed to set active user", ctx)

	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"data": au,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}

// GetActiveUsers implements UserController.
func (c *controller) GetActiveUsers(ctx *fiber.Ctx) error {
	key := ctx.Params("key", "")
	au, err := c.activeUserService.GetOnlineUser(key)
	utils.ReturnErrMessageIfErr(err, "Failed to get active users", ctx)

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": au,
		"meta": &fiber.Map{
			"key": key,
		},
		"time": time.Now(),
	})
}

// OnUserHookHandler implements UserController.
func (c *controller) OnUserHookHandler(ctx *fiber.Ctx) error {

	// Parsing user request
	var user service.ActiveUser
	if err := ctx.BodyParser(&user); err != nil {
		utils.ReturnErrMessageIfErr(err, "Creating user failed, error", ctx)
	}

	// Getting header information
	var header = new(Headers)
	if err := ctx.ReqHeaderParser(header); err != nil {
		utils.ReturnErrMessageIfErr(err, "onUserHookHandler", ctx)
	}

	// Handling hook request
	switch header.VerneHook {
	case "on_client_wakeup":
		c.activeUserService.SetOnlineUser(&user)
	case "on_client_offline":
		// TODO: add handler for deleting client
		log.Info(user)
	}

	// Send as monitor topics
	// TODO: change topic handler
	if err := service.MqttClient.Publish(context.Background(), "xops/api/user", user, courier.QOSTwo); err != nil {
		log.Error("Error publishing, error ", err)
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": user,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}
