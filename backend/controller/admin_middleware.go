package controller

import (
	"context"
	"time"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// RequireAdmin adalah middleware yang memastikan hanya user dengan role
// "admin" yang bisa mengakses endpoint tertentu (mis. kelola booking,
// broadcast notifikasi). Pasang middleware ini di url.go sebelum handler
// endpoint admin.
//
// Cara set user jadi admin: update manual di MongoDB, contoh lewat
// MongoDB Compass/Shell:
//
//	db.users.updateOne({ email: "email-admin@contoh.com" }, { $set: { role: "admin" } })
func RequireAdmin(c *fiber.Ctx) error {
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

	if dbUser.Role != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses ditolak, khusus admin"})
	}

	// Simpan email admin di context kalau nanti dibutuhkan handler
	c.Locals("adminEmail", session.Email)

	return c.Next()
}