package controller

import (
	"context"
	"gocroot/config"
	"gocroot/helper"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetAllDestinasi(c *fiber.Ctx) error {
	kategoriUser := c.Query("kategori")
	filter := bson.M{}
	if kategoriUser != "" {
		filter = bson.M{"kategori": kategoriUser}
	}
	destinasis, err := helper.GetManyDoc[model.Destinations](config.Mongoconn, "destinasi", filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(destinasis)
}
func GetOneDestinasi(c *fiber.Ctx) error {
	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "ID destinasi tidak valid",
		})
	}

	var destinasi model.Destinations
	err = config.Mongoconn.Collection("destinasi").FindOne(context.Background(), bson.M{"_id": oid}).Decode(&destinasi)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  false,
			"message": "Destinasi tidak ditemukan",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": true,
		"data":   destinasi,
	})
}