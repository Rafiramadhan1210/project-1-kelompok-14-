package controller

import (
	"context"
	"fmt"
	"strings"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// namaDestinasiMilikPengelola mengembalikan daftar nama destinasi yang dimiliki
// oleh email pengelola tertentu. Dipakai untuk memfilter booking, karena field
// Booking.Destination menyimpan nama destinasi (bukan ID).
func namaDestinasiMilikPengelola(email string) ([]string, error) {
	db := config.Mongoconn
	cursor, err := db.Collection("destinasi").Find(context.Background(), bson.M{"pengelola_email": email})
	if err != nil {
		return nil, err
	}
	var destinasiList []model.Destinations
	if err := cursor.All(context.Background(), &destinasiList); err != nil {
		return nil, err
	}
	names := make([]string, 0, len(destinasiList))
	for _, d := range destinasiList {
		names = append(names, d.Nama)
	}
	return names, nil
}

// GetPengelolaBookings mengembalikan semua booking yang masuk untuk destinasi
// milik pengelola yang sedang login (bukan semua booking di platform).
// Bisa difilter status lewat query ?status=Pending,Dibayar
func GetPengelolaBookings(c *fiber.Ctx) error {
	db := config.Mongoconn
	email := c.Locals("pengelolaEmail").(string)

	names, err := namaDestinasiMilikPengelola(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	if len(names) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": []model.Booking{}})
	}

	filter := bson.M{"destination": bson.M{"$in": names}}

	statusParam := c.Query("status")
	if statusParam != "" {
		filter["status"] = bson.M{"$in": strings.Split(statusParam, ",")}
	}

	opts := options.Find().SetSort(bson.D{{Key: "_id", Value: -1}})
	var data []model.Booking
	cursor, err := db.Collection("bookings").Find(context.Background(), filter, opts)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	if err := cursor.All(context.Background(), &data); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

// UpdateBookingStatusPengelola mengubah status booking, tapi hanya untuk booking
// di destinasi milik pengelola yang sedang login.
// Body JSON: { "booking_id": "...", "status": "Dibayar" | "Checked-in" | "Selesai" | "Dibatalkan" }
func UpdateBookingStatusPengelola(c *fiber.Ctx) error {
	db := config.Mongoconn
	email := c.Locals("pengelolaEmail").(string)

	var body struct {
		BookingID       string `json:"booking_id"`
		Status          string `json:"status"`
		KeteranganTolak string `json:"keterangan_tolak"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}

	validStatus := map[string]bool{"Dibayar": true, "Checked-in": true, "Selesai": true, "Dibatalkan": true}
	if !validStatus[body.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Status tidak valid"})
	}

	oid, err := primitive.ObjectIDFromHex(body.BookingID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "booking_id tidak valid"})
	}

	var booking model.Booking
	if err := db.Collection("bookings").FindOne(context.Background(), bson.M{"_id": oid}).Decode(&booking); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Booking tidak ditemukan"})
	}

	// Pastikan booking ini untuk destinasi milik pengelola yang login
	names, err := namaDestinasiMilikPengelola(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	owned := false
	for _, n := range names {
		if n == booking.Destination {
			owned = true
			break
		}
	}
	if !owned {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Booking ini bukan untuk destinasi milikmu"})
	}

	// Check-in hanya boleh dari status Dibayar
	if body.Status == "Checked-in" && booking.Status != "Dibayar" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Hanya booking berstatus Dibayar yang bisa di-check-in"})
	}

	updateFields := bson.M{"status": body.Status}
	if body.Status == "Dibatalkan" {
		updateFields["keterangan_tolak"] = body.KeteranganTolak
	} else {
		updateFields["keterangan_tolak"] = ""
	}

	_, err = db.Collection("bookings").UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$set": updateFields},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	statusMessage := map[string]string{
		"Dibayar":    fmt.Sprintf("Pembayaran booking %s telah dikonfirmasi. Selamat berlibur!", booking.Destination),
		"Checked-in": fmt.Sprintf("Selamat datang! Kunjungan kamu ke %s sudah dikonfirmasi hadir.", booking.Destination),
		"Selesai":    fmt.Sprintf("Booking %s telah selesai. Terima kasih telah menggunakan GoTrip!", booking.Destination),
		"Dibatalkan": fmt.Sprintf("Booking %s dibatalkan oleh pengelola.", booking.Destination),
	}
	if msg, ok := statusMessage[body.Status]; ok {
		_ = CreateNotification(booking.Email, "booking", "Status Booking Diperbarui", msg, "booking-saya.html")
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Status booking diperbarui"})
}
