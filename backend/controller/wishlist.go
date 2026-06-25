package controller

import (
	"context"
	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type wishlistRequest struct {
	DestinationID string `json:"destination_id"`
}

// ToggleWishlist menambah/menghapus destinasi dari wishlist user yang sedang login.
// Kalau destinasi sudah ada di wishlist, akan dihapus. Kalau belum ada, akan ditambahkan.
func ToggleWishlist(c *fiber.Ctx) error {
	db := config.Mongoconn

	email, err := getSessionEmail(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu"})
	}

	var body wishlistRequest
	if err := c.BodyParser(&body); err != nil || body.DestinationID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "destination_id wajib diisi"})
	}

	var user model.Users
	err = db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	alreadyWishlisted := false
	for _, id := range user.Wishlist {
		if id == body.DestinationID {
			alreadyWishlisted = true
			break
		}
	}

	var update bson.M
	if alreadyWishlisted {
		update = bson.M{"$pull": bson.M{"wishlist": body.DestinationID}}
	} else {
		update = bson.M{"$addToSet": bson.M{"wishlist": body.DestinationID}}
	}

	_, err = db.Collection("users").UpdateOne(context.Background(), bson.M{"email": email}, update)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Berhasil diperbarui",
		"wishlisted": !alreadyWishlisted,
	})
}

// GetMyWishlist mengembalikan detail destinasi yang sudah di-wishlist user yang login
func GetMyWishlist(c *fiber.Ctx) error {
	db := config.Mongoconn

	email, err := getSessionEmail(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu"})
	}

	var user model.Users
	err = db.Collection("users").FindOne(context.Background(), bson.M{"email": email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "User tidak ditemukan"})
	}

	if len(user.Wishlist) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": []model.Destinations{}})
	}

	var objectIDs []primitive.ObjectID
	for _, idStr := range user.Wishlist {
		oid, err := primitive.ObjectIDFromHex(idStr)
		if err == nil {
			objectIDs = append(objectIDs, oid)
		}
	}

	var destinations []model.Destinations
	cursor, err := db.Collection("destinasi").Find(context.Background(), bson.M{"_id": bson.M{"$in": objectIDs}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	if err := cursor.All(context.Background(), &destinations); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": destinations})
}