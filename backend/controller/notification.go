package controller

import (
	"context"
	"time"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CreateNotification membuat satu notifikasi baru di database.
// email kosong ("") berarti notifikasi ini untuk SEMUA user (broadcast/promo).
func CreateNotification(email, notifType, title, message, link string) error {
	db := config.Mongoconn
	notif := model.Notification{
		Email:     email,
		Type:      notifType,
		Title:     title,
		Message:   message,
		Link:      link,
		ReadBy:    []string{},
		CreatedAt: primitive.NewDateTimeFromTime(time.Now()),
	}
	_, err := db.Collection("notifications").InsertOne(context.Background(), notif)
	return err
}

// GetMyNotifications mengembalikan notifikasi milik user yang login,
// digabung dengan notifikasi broadcast (promo) untuk semua user.
func GetMyNotifications(c *fiber.Ctx) error {
	db := config.Mongoconn

	email, err := getSessionEmail(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu"})
	}

	filter := bson.M{"$or": []bson.M{
		{"email": email},
		{"email": ""},
	}}

	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(30)
	cursor, err := db.Collection("notifications").Find(context.Background(), filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	defer cursor.Close(context.Background())

	var notifs []model.Notification
	if err := cursor.All(context.Background(), &notifs); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	result := make([]fiber.Map, 0, len(notifs))
	unread := 0
	for _, n := range notifs {
		isRead := containsEmail(n.ReadBy, email)
		if !isRead {
			unread++
		}
		result = append(result, fiber.Map{
			"_id":        n.ID,
			"type":       n.Type,
			"title":      n.Title,
			"message":    n.Message,
			"link":       n.Link,
			"created_at": n.CreatedAt,
			"is_read":    isRead,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":         result,
		"unread_count": unread,
	})
}

// GetUnreadNotificationCount mengembalikan jumlah notifikasi yang belum dibaca saja
// (dipakai untuk badge angka di icon bell, dipanggil berkala / polling).
func GetUnreadNotificationCount(c *fiber.Ctx) error {
	db := config.Mongoconn

	email, err := getSessionEmail(c)
	if err != nil {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"unread_count": 0})
	}

	filter := bson.M{
		"$or": []bson.M{
			{"email": email},
			{"email": ""},
		},
		"read_by": bson.M{"$ne": email},
	}

	count, err := db.Collection("notifications").CountDocuments(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"unread_count": count})
}

// MarkNotificationRead menandai satu notifikasi sebagai sudah dibaca oleh user yang login.
// Body JSON: { "notification_id": "..." }
func MarkNotificationRead(c *fiber.Ctx) error {
	db := config.Mongoconn

	email, err := getSessionEmail(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu"})
	}

	var body struct {
		NotificationID string `json:"notification_id"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}

	oid, err := primitive.ObjectIDFromHex(body.NotificationID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "notification_id tidak valid"})
	}

	_, err = db.Collection("notifications").UpdateOne(
		context.Background(),
		bson.M{"_id": oid, "$or": []bson.M{{"email": email}, {"email": ""}}},
		bson.M{"$addToSet": bson.M{"read_by": email}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Notifikasi ditandai sudah dibaca"})
}

// MarkAllNotificationsRead menandai semua notifikasi milik user (+ broadcast) sebagai sudah dibaca.
func MarkAllNotificationsRead(c *fiber.Ctx) error {
	db := config.Mongoconn

	email, err := getSessionEmail(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu"})
	}

	filter := bson.M{"$or": []bson.M{
		{"email": email},
		{"email": ""},
	}}

	_, err = db.Collection("notifications").UpdateMany(
		context.Background(),
		filter,
		bson.M{"$addToSet": bson.M{"read_by": email}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Semua notifikasi ditandai sudah dibaca"})
}

// BroadcastPromoNotification dipakai dari Admin Panel untuk mengirim notifikasi promo ke semua user.
// Body JSON: { "title": "...", "message": "...", "link": "..." (opsional) }
func BroadcastPromoNotification(c *fiber.Ctx) error {
	var body struct {
		Title   string `json:"title"`
		Message string `json:"message"`
		Link    string `json:"link"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}

	if body.Title == "" || body.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Judul dan pesan wajib diisi"})
	}

	if err := CreateNotification("", "promo", body.Title, body.Message, body.Link); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Notifikasi promo berhasil dikirim ke semua pengguna"})
}

// containsEmail helper kecil untuk cek apakah email ada di dalam slice
func containsEmail(list []string, email string) bool {
	for _, e := range list {
		if e == email {
			return true
		}
	}
	return false
}