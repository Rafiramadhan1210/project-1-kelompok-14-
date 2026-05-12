package url

import (
	"gocroot/controller"

	"github.com/gofiber/fiber/v2"
)

func Web(page *fiber.App) {
	page.Get("/", controller.Homepage)
	page.Get("/ip", controller.GetIPServer)
	page.Get("/whatsauth/refreshtoken", controller.RefreshWAToken)
	page.Post("/whatsauth/webhook", controller.WhatsAuthReceiver)
	page.Post("/booking", controller.InsertBooking)
	page.Get("/booking", controller.GetAllBooking)
	page.Post("/register", controller.RegisterUser)
	page.Post("/Login", controller.LoginUser)
}
