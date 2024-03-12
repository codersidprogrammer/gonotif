package notification

import (
	"github.com/codersidprogrammer/gonotif/app/notification/controller"
	"github.com/codersidprogrammer/gonotif/app/notification/service"
	"github.com/gofiber/fiber/v2"
)

type appNotification struct {
	app         *fiber.App
	pushService service.NotificationPushService
	controller  controller.NotificationController
}

type AppNotification interface {
	Route()
}

func NewAppNotification(app *fiber.App) AppNotification {
	nps := service.NewNotificationPushService()
	return &appNotification{
		app:         app,
		controller:  controller.NewNotificationController(service.NewNotificationBucketService(), nps),
		pushService: nps,
	}
}

// Route implements AppNotification.
func (a *appNotification) Route() {
	go a.pushService.SubsHandler()

	r := a.app.Group("/v1/notification")

	r.Post("/bucket", a.controller.CreateNotification)
	r.Post("/push", a.controller.Publish)

	r.Get("/queue", a.controller.CreateNotificationWorker)
	r.Get("/subscribe/:id", a.controller.Subscribe)
	r.Get("/bucket", a.controller.GetNotifications)
	r.Get("/bucket/:id", a.controller.GetNotification)
}
