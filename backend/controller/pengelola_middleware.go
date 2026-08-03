package controller

import (
	"context"
	"time"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// RequirePengelola adalah middleware yang memastikan hanya user dengan role
// "pengelola_wisata" yang bisa mengakses endpoint kelola destinasi/booking
// miliknya sendiri. Polanya sama seperti RequireAdmin di admin_middleware.go.
func RequirePengelola(c *fiber.Ctx) error {
	db := config.Mongoconn
	token := c.Cookies(sessionCookieName)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu"})
	}

	var session model.Session
	err := db.Collection("sessions").FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session tidak valid"})
	}

	if session.ExpiresAt.Time().Before(time.Now()) {
		_, _ = db.Collection("sessions").DeleteOne(context.Background(), bson.M{"token": token})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session telah berakhir"})
	}

	var dbUser model.Users
	if err := db.Collection("users").FindOne(context.Background(), bson.M{"email": session.Email}).Decode(&dbUser); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	if dbUser.Role != "pengelola_wisata" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses ditolak, khusus pengelola wisata"})
	}

	// Simpan email pengelola di context supaya handler tahu siapa yang login,
	// dipakai untuk filter data & cek kepemilikan.
	c.Locals("pengelolaEmail", session.Email)

	return c.Next()
}
