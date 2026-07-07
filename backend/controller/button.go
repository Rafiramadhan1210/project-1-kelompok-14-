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

	// println("Database =", db.Name())

	// collections, _ := db.ListCollectionNames(context.Background(), bson.M{})
	// println("Collections:")
	// for _, c := range collections {
	// 	println("-", c)
	// }

	collection := db.Collection("destinasi")
	// println("Collection =", collection.Name())

	var destinations []model.Destinations

	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	err = cursor.All(context.Background(), &destinations)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// println("Jumlah destinasi:", len(destinations))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "INI TEST DARI BUTTON.GO",
		"data":    destinations,
	})
}
