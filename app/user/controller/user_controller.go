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
func (*controller) OnUserHookHandler(ctx *fiber.Ctx) error {
	var body interface{}
	if err := ctx.BodyParser(&body); err != nil {
		utils.ReturnErrMessageIfErr(err, "onUserHookHandler", ctx)
	}

	type Headers struct {
		VerneHook string `reqHeader:"Vernemq-Hook"`
	}
	var header = new(Headers)
	if err := ctx.ReqHeaderParser(header); err != nil {
		utils.ReturnErrMessageIfErr(err, "onUserHookHandler", ctx)
	}
	log.Info(header)
	log.Info(string(ctx.Request().Header.Header()))

	if err := service.MqttClient.Publish(context.Background(), "xops/api/user", body, courier.QOSTwo); err != nil {
		log.Error("Error publishing, error ", err)
	} else {
		log.Debug("Published")
		log.Debug(body)
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"data": nil,
		"meta": &fiber.Map{},
		"time": time.Now(),
	})
}
