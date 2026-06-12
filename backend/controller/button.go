package controller

import (
    "context"
    "gocroot/config"
    "gocroot/model"

    "github.com/gofiber/fiber/v2"
    "go.mongodb.org/mongo-driver/bson"
)

func GetButton(c *fiber.Ctx) error {
    db := config.Mongoconn
    // 1. Ubah model.Destinasi menjadi model.Destinations (Pakai S)
    var destinations []model.Destinations 

    cursor, err := db.Collection("destinasi").Find(context.Background(), bson.M{})
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    // 2. Data cursor akan otomatis di-parsing ke struct model.Destinations yang baru
    err = cursor.All(context.Background(), &destinations)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Daftar Destinasi GoTrip",
        "data":    destinations,
    })
}