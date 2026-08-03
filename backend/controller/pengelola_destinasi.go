package controller

import (
	"context"
	"strings"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetMyDestinasi mengembalikan semua destinasi milik pengelola yang sedang login.
func GetMyDestinasi(c *fiber.Ctx) error {
	db := config.Mongoconn
	email := c.Locals("pengelolaEmail").(string)

	var data []model.Destinations
	cursor, err := db.Collection("destinasi").Find(context.Background(), bson.M{"pengelola_email": email})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	if err := cursor.All(context.Background(), &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

// CreateDestinasi menambahkan destinasi baru milik pengelola yang sedang login.
// Body JSON: { "nama": "...", "deskripsi": "...", "kategori": "...", "harga": 100000, "gambar": "...", "lokasi": "..." }
func CreateDestinasi(c *fiber.Ctx) error {
	db := config.Mongoconn
	email := c.Locals("pengelolaEmail").(string)

	var body model.Destinations
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}

	body.Nama = strings.TrimSpace(body.Nama)
	body.Deskripsi = strings.TrimSpace(body.Deskripsi)
	body.Kategori = strings.TrimSpace(body.Kategori)
	body.Lokasi = strings.TrimSpace(body.Lokasi)

	if body.Nama == "" || body.Kategori == "" || body.Harga <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Nama, kategori, dan harga wajib diisi dengan benar"})
	}

	// Paksa pemilik destinasi = pengelola yang sedang login, ID selalu baru
	body.ID = primitive.NewObjectID()
	body.PengelolaEmail = email

	_, err := db.Collection("destinasi").InsertOne(context.Background(), body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Destinasi berhasil ditambahkan",
		"data":    body,
	})
}

// UpdateDestinasi mengedit destinasi, hanya boleh kalau pengelola yang login adalah pemiliknya.
// Body JSON: field-field yang mau diubah (nama, deskripsi, kategori, harga, gambar, lokasi)
func UpdateDestinasi(c *fiber.Ctx) error {
	db := config.Mongoconn
	email := c.Locals("pengelolaEmail").(string)

	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID destinasi tidak valid"})
	}

	// Pastikan destinasi ini memang miliknya
	var existing model.Destinations
	err = db.Collection("destinasi").FindOne(context.Background(), bson.M{"_id": oid}).Decode(&existing)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Destinasi tidak ditemukan"})
	}
	if existing.PengelolaEmail != email {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Kamu bukan pemilik destinasi ini"})
	}

	var body model.Destinations
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}

	update := bson.M{}
	if strings.TrimSpace(body.Nama) != "" {
		update["nama"] = strings.TrimSpace(body.Nama)
	}
	if strings.TrimSpace(body.Deskripsi) != "" {
		update["deskripsi"] = strings.TrimSpace(body.Deskripsi)
	}
	if strings.TrimSpace(body.Kategori) != "" {
		update["kategori"] = strings.TrimSpace(body.Kategori)
	}
	if body.Harga > 0 {
		update["harga"] = body.Harga
	}
	if strings.TrimSpace(body.Gambar) != "" {
		update["gambar"] = strings.TrimSpace(body.Gambar)
	}
	if strings.TrimSpace(body.Lokasi) != "" {
		update["lokasi"] = strings.TrimSpace(body.Lokasi)
	}

	if len(update) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Tidak ada data yang diubah"})
	}

	_, err = db.Collection("destinasi").UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$set": update},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Destinasi berhasil diperbarui"})
}

// DeleteDestinasi menghapus destinasi, hanya boleh kalau pengelola yang login adalah pemiliknya.
func DeleteDestinasi(c *fiber.Ctx) error {
	db := config.Mongoconn
	email := c.Locals("pengelolaEmail").(string)

	id := c.Params("id")
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID destinasi tidak valid"})
	}

	var existing model.Destinations
	err = db.Collection("destinasi").FindOne(context.Background(), bson.M{"_id": oid}).Decode(&existing)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Destinasi tidak ditemukan"})
	}
	if existing.PengelolaEmail != email {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Kamu bukan pemilik destinasi ini"})
	}

	_, err = db.Collection("destinasi").DeleteOne(context.Background(), bson.M{"_id": oid})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Destinasi berhasil dihapus"})
}
