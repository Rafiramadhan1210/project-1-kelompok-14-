package controller

import (
	"context"
	"gocroot/config"
	"gocroot/model"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

// RegisterUser untuk mendaftarkan akun baru
func RegisterUser(c *fiber.Ctx) error {
	db := config.Mongoconn
	var user model.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	// Hash password
	bytes, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	user.Password = string(bytes)
	_, err := db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Registrasi Berhasil!"})
}

// LoginUser untuk masuk ke sistem
func LoginUser(c *fiber.Ctx) error {
	db := config.Mongoconn
	var user model.Users
	var dbUser model.Users

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Cari user berdasarkan username
	err := db.Collection("users").FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email atau Password salah!"})
	}

	// Cek password (untuk tahap awal kita bandingkan string langsung)
	if user.Password != dbUser.Password {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email atau Password salah!"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Login Berhasil!"})
}