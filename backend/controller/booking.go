package controller

import (
	"context"
	"gocroot/config"
	"gocroot/model"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InsertBooking(c *fiber.Ctx) error {
	db := config.Mongoconn
	var pesanan model.Booking
	
	// Membaca data dari input user
	if err := c.BodyParser(&pesanan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Otomatis mencatat waktu booking
	pesanan.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	// Simpan ke koleksi "bookings" di MongoDB Atlas
	_, err := db.Collection("bookings").InsertOne(context.Background(), pesanan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Booking GoTrip Berhasil!",
		"data":    pesanan,
	})
}
func GetAllBooking(c *fiber.Ctx) error {
	db := config.Mongoconn
	var data []model.Booking

	cursor, err := db.Collection("bookings").Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err = cursor.All(context.Background(), &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(data)
}