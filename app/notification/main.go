package notification

import (
	"github.com/codersidprogrammer/gonotif/app/notification/controller"
	"github.com/codersidprogrammer/gonotif/app/notification/service"
	"github.com/gofiber/fiber/v2"
)

type appNotification struct {
	app        *fiber.App
	controller controller.NotificationController
}

type AppNotification interface {
	Route()
}

func NewAppNotification(app *fiber.App) AppNotification {
	return &appNotification{
		app:        app,
		controller: controller.NewNotificationController(service.NewNotificationService()),
	}
}

// Route implements AppNotification.
func (a *appNotification) Route() {
	r := a.app.Group("/v1/notification")

	r.Post("/bucket", a.controller.CreateNotification)
	r.Get("/bucket", a.controller.GetNotifications)
	r.Get("/bucket/:id", a.controller.GetNotification)
}
