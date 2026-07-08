package controller

import (
	"context"
	"strings"
	"time"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// SubmitSupportMessage menerima pesan dari form "Hubungi Kami" di bantuan.html.
// Bisa dipakai user yang login maupun belum (tamu), makanya email diambil dari
// body request kalau session tidak ada.
// Body JSON: { "nama": "...", "email": "...", "topik": "...", "pesan": "..." }
func SubmitSupportMessage(c *fiber.Ctx) error {
	db := config.Mongoconn

	var body struct {
		Nama  string `json:"nama"`
		Email string `json:"email"`
		Topik string `json:"topik"`
		Pesan string `json:"pesan"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}

	body.Nama = strings.TrimSpace(body.Nama)
	body.Email = strings.TrimSpace(body.Email)
	body.Topik = strings.TrimSpace(body.Topik)
	body.Pesan = strings.TrimSpace(body.Pesan)

	// Kalau user sedang login, pakai email dari session sebagai sumber kebenaran
	if email, err := getSessionEmail(c); err == nil && email != "" {
		body.Email = email
	}

	if body.Nama == "" || body.Topik == "" || body.Pesan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Nama, topik, dan pesan wajib diisi"})
	}

	pesan := model.SupportMessage{
		Nama:      body.Nama,
		Email:     body.Email,
		Topik:     body.Topik,
		Pesan:     body.Pesan,
		Status:    "Baru",
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err := db.Collection("support_messages").InsertOne(context.Background(), pesan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengirim pesan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Pesan berhasil dikirim"})
}

// GetAllSupportMessages mengembalikan semua pesan bantuan, dipakai di Admin Panel.
// Endpoint ini dilindungi middleware RequireAdmin.
func GetAllSupportMessages(c *fiber.Ctx) error {
	db := config.Mongoconn

	opts := options.Find().SetSort(bson.M{"created_at": -1})
	cursor, err := db.Collection("support_messages").Find(context.Background(), bson.M{}, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	var data []model.SupportMessage
	if err := cursor.All(context.Background(), &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

// MarkSupportMessageReplied menandai satu pesan bantuan sudah dibalas/ditindaklanjuti.
// Body JSON: { "id": "..." }
func MarkSupportMessageReplied(c *fiber.Ctx) error {
	db := config.Mongoconn

	var body struct {
		ID string `json:"id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}

	oid, err := primitive.ObjectIDFromHex(body.ID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID tidak valid"})
	}

	_, err = db.Collection("support_messages").UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": "Dibalas"}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Ditandai sudah dibalas"})
}