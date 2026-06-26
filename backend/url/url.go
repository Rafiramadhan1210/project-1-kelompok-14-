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
	page.Get("/api/my-bookings", controller.GetMyBookings)
	page.Post("/register", controller.RegisterUser)
	page.Post("/Login", controller.LoginUser)
	page.Post("/logout", controller.LogoutUser)
	page.Get("/api/me", controller.GetCurrentUser)
	page.Post("/api/profile/update", controller.UpdateProfile)
	page.Post("/api/profile/photo", controller.UploadProfilePhoto)
	page.Get("/button", controller.GetButton)
	page.Post("/api/wishlist/toggle", controller.ToggleWishlist)
	page.Get("/api/my-wishlist", controller.GetMyWishlist)
	page.Get("/destinasi", controller.GetAllDestinasi)
	page.Get("/destinasi/:id", controller.GetOneDestinasi)
	page.Post("/api/booking/update-status", controller.UpdateBookingStatus)
	page.Delete("/api/my-bookings/:id", controller.CancelMyBooking)
	page.Get("/api/notifications", controller.GetMyNotifications)
	page.Get("/api/notifications/unread-count", controller.GetUnreadNotificationCount)
	page.Post("/api/notifications/read", controller.MarkNotificationRead)
	page.Post("/api/notifications/read-all", controller.MarkAllNotificationsRead)
	page.Post("/api/notifications/broadcast", controller.BroadcastPromoNotification)
}