package notification

import (
	"github.com/codersidprogrammer/gonotif/app/notification/controller"
	"github.com/codersidprogrammer/gonotif/pkg/model"
	"github.com/gofiber/fiber/v2"
)

type appNotification struct {
	app        *fiber.App
	controller controller.NotificationController
}

func NewAppNotification(app *fiber.App) model.App {
	return &appNotification{
		app:        app,
		controller: controller.NewNotificationController(),
	}
}

// Route implements AppNotification.
func (a *appNotification) Route() {
	r := a.app.Group("/v1/notification")
	r.Post("/queue", a.controller.CreatePushNotification)
}
