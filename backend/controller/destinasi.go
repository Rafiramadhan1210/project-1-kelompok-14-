package controller

import (
	"gocroot/config"
	"gocroot/helper"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
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
