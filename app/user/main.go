package user

import (
	"github.com/codersidprogrammer/gonotif/app/user/controller"
	"github.com/codersidprogrammer/gonotif/pkg/model"
	"github.com/gofiber/fiber/v2"
)

type userApp struct {
	app        *fiber.App
	controller controller.UserController
}

func NewUserApp(app *fiber.App) model.App {
	return &userApp{
		app:        app,
		controller: controller.NewUserAppController(),
	}
}

// Route implements model.App.
func (a *userApp) Route() {
	r := a.app.Group("/v1/user")

	r.Post("/hook", a.controller.OnUserHookHandler)
}
