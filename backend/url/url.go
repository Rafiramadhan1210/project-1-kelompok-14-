package url

import (
	"gocroot/controller"

	"github.com/gofiber/fiber/v2"
)

func Web(page *fiber.App) {
	page.Get("/ip", controller.GetIPServer)
	page.Get("/whatsauth/refreshtoken", controller.RefreshWAToken)
	page.Post("/whatsauth/webhook", controller.WhatsAuthReceiver)
	page.Post("/booking", controller.InsertBooking)
	page.Get("/booking", controller.RequireAdmin, controller.GetAllBooking)
	page.Get("/api/my-bookings", controller.GetMyBookings)
	page.Get("/api/my-bookings/:id", controller.GetMyBookingByID)
	page.Post("/register", controller.RegisterUser)
	page.Post("/Login", controller.LoginUser)
	page.Post("/api/auth/google", controller.GoogleLogin)
	page.Post("/logout", controller.LogoutUser)
	page.Get("/api/me", controller.GetCurrentUser)
	page.Post("/api/profile/update", controller.UpdateProfile)
	page.Post("/api/profile/photo", controller.UploadProfilePhoto)
	page.Put("/api/user/password", controller.ChangePassword)
	page.Delete("/api/user", controller.DeleteAccount)
	page.Get("/button", controller.GetButton)
	page.Post("/api/wishlist/toggle", controller.ToggleWishlist)
	page.Get("/api/my-wishlist", controller.GetMyWishlist)
	page.Get("/destinasi", controller.GetAllDestinasi)
	page.Get("/destinasi/:id", controller.GetOneDestinasi)
	page.Post("/api/booking/update-status", controller.RequireAdmin, controller.UpdateBookingStatus)
	page.Delete("/api/my-bookings/:id", controller.CancelMyBooking)
	page.Post("/api/my-bookings/:id/bukti-bayar", controller.ReuploadBuktiBayar)
	page.Get("/api/notifications", controller.GetMyNotifications)
	page.Get("/api/notifications/unread-count", controller.GetUnreadNotificationCount)
	page.Post("/api/notifications/read", controller.MarkNotificationRead)
	page.Post("/api/notifications/read-all", controller.MarkAllNotificationsRead)
	page.Post("/api/notifications/broadcast", controller.RequireAdmin, controller.BroadcastPromoNotification)
	page.Post("/api/support/message", controller.SubmitSupportMessage)
	page.Get("/api/support/messages", controller.RequireAdmin, controller.GetAllSupportMessages)
	page.Post("/api/support/messages/mark-replied", controller.RequireAdmin, controller.MarkSupportMessageReplied)

	// Pengelola Wisata: kelola destinasi milik sendiri
	page.Get("/api/pengelola/destinasi", controller.RequirePengelola, controller.GetMyDestinasi)
	page.Post("/api/pengelola/destinasi", controller.RequirePengelola, controller.CreateDestinasi)
	page.Put("/api/pengelola/destinasi/:id", controller.RequirePengelola, controller.UpdateDestinasi)
	page.Delete("/api/pengelola/destinasi/:id", controller.RequirePengelola, controller.DeleteDestinasi)

	// Pengelola Wisata: booking untuk destinasi milik sendiri + check-in
	page.Get("/api/pengelola/booking", controller.RequirePengelola, controller.GetPengelolaBookings)
	page.Post("/api/pengelola/booking/update-status", controller.RequirePengelola, controller.UpdateBookingStatusPengelola)

	// Pengelola Wisata: statistik pendapatan & tren pengunjung
	page.Get("/api/pengelola/statistik", controller.RequirePengelola, controller.GetPengelolaStatistik)
}