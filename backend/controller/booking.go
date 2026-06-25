package controller

import (
	"context"
	"fmt"
	"gocroot/config"
	"gocroot/model"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// getSessionEmail mengambil email user yang sedang login dari session cookie.
func getSessionEmail(c *fiber.Ctx) (string, error) {
	db := config.Mongoconn
	token := c.Cookies(sessionCookieName)
	if token == "" {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Belum login")
	}

	var session model.Session
	err := db.Collection("sessions").FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Session tidak valid")
	}

	if session.ExpiresAt.Time().Before(time.Now()) {
		return "", fiber.NewError(fiber.StatusUnauthorized, "Session telah berakhir")
	}

	return session.Email, nil
}

// InsertBooking membuat booking baru untuk user yang sedang login,
// menerima multipart/form-data karena ada file bukti bayar.
func InsertBooking(c *fiber.Ctx) error {
	db := config.Mongoconn

	email, err := getSessionEmail(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu untuk booking"})
	}

	namaUser := strings.TrimSpace(c.FormValue("nama_user"))
	noHP := strings.TrimSpace(c.FormValue("no_hp"))
	destination := strings.TrimSpace(c.FormValue("destination"))
	tanggalKunjungan := strings.TrimSpace(c.FormValue("tanggal_kunjungan"))
	setuju := c.FormValue("setuju_syarat")

	if namaUser == "" || noHP == "" || destination == "" || tanggalKunjungan == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Nama, No. HP, destinasi, dan tanggal kunjungan wajib diisi"})
	}

	if setuju != "true" && setuju != "on" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Kamu harus menyetujui syarat & ketentuan"})
	}

	totalTiket := 1
	if v := c.FormValue("total_tiket"); v != "" {
		fmt.Sscanf(v, "%d", &totalTiket)
	}
	if totalTiket < 1 {
		totalTiket = 1
	}

	totalBayar := 0
	if v := c.FormValue("total_bayar"); v != "" {
		fmt.Sscanf(v, "%d", &totalBayar)
	}

	// Proses upload bukti bayar
	fileHeader, err := c.FormFile("bukti_bayar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bukti bayar wajib diupload"})
	}

	uploadDir := filepath.Join(".", "frontend", "public", "uploads", "bukti-bayar")
	if _, statErr := os.Stat(uploadDir); os.IsNotExist(statErr) {
		// fallback kalau cwd-nya sudah di dalam folder backend
		uploadDir = filepath.Join("..", "frontend", "public", "uploads", "bukti-bayar")
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyiapkan folder upload"})
	}

	ext := filepath.Ext(fileHeader.Filename)
	uniqueName := fmt.Sprintf("%s-%d%s", email, time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, uniqueName)

	if err := c.SaveFile(fileHeader, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan bukti bayar"})
	}

	buktiBayarURL := "/uploads/bukti-bayar/" + uniqueName

	pesanan := model.Booking{
		Email:            email,
		NamaUser:         namaUser,
		NoHP:             noHP,
		Destination:      destination,
		TotalTiket:       totalTiket,
		TotalBayar:       totalBayar,
		TanggalKunjungan: tanggalKunjungan,
		BuktiBayar:       buktiBayarURL,
		Status:           "Pending",
		CreatedAt:        primitive.NewDateTimeFromTime(time.Now()),
	}

	_, err = db.Collection("bookings").InsertOne(context.Background(), pesanan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Booking berhasil dibuat! Menunggu konfirmasi pembayaran.",
		"data":    pesanan,
	})
}

// GetAllBooking tetap ada untuk keperluan admin/internal (semua booking)
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

// GetMyBookings mengembalikan booking milik user yang sedang login saja,
// bisa difilter berdasarkan status lewat query ?status=Pending,Dibayar
func GetMyBookings(c *fiber.Ctx) error {
	db := config.Mongoconn

	email, err := getSessionEmail(c)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu"})
	}

	filter := bson.M{"email": email}

	statusParam := c.Query("status")
	if statusParam != "" {
		statuses := strings.Split(statusParam, ",")
		filter["status"] = bson.M{"$in": statuses}
	}

	var data []model.Booking
	cursor, err := db.Collection("bookings").Find(context.Background(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	err = cursor.All(context.Background(), &data)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": data,
	})
}

// UpdateBookingStatus untuk admin mengonfirmasi/menolak/menyelesaikan booking.
// Body JSON: { "booking_id": "...", "status": "Dibayar" }
func UpdateBookingStatus(c *fiber.Ctx) error {
	db := config.Mongoconn

	var body struct {
		BookingID string `json:"booking_id"`
		Status    string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}

	validStatus := map[string]bool{"Pending": true, "Dibayar": true, "Selesai": true, "Dibatalkan": true}
	if !validStatus[body.Status] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Status tidak valid"})
	}

	oid, err := primitive.ObjectIDFromHex(body.BookingID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "booking_id tidak valid"})
	}

	_, err = db.Collection("bookings").UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{"status": body.Status}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Status booking diperbarui"})
}

// CancelMyBooking memungkinkan user membatalkan booking milik mereka sendiri (hanya yang Pending)
func CancelMyBooking(c *fiber.Ctx) error {
    db := config.Mongoconn

    email, err := getSessionEmail(c)
    if err != nil {
        return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Silakan login dulu"})
    }

    bookingID := c.Params("id")
    oid, err := primitive.ObjectIDFromHex(bookingID)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID booking tidak valid"})
    }

    // Pastikan booking ini milik user yang login dan masih Pending
    var booking model.Booking
    err = db.Collection("bookings").FindOne(context.Background(), bson.M{
        "_id":   oid,
        "email": email,
    }).Decode(&booking)
    if err != nil {
        return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Booking tidak ditemukan"})
    }

    if booking.Status != "Pending" {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "message": "Hanya booking berstatus Pending yang dapat dibatalkan",
        })
    }

    _, err = db.Collection("bookings").UpdateOne(
        context.Background(),
        bson.M{"_id": oid},
        bson.M{"$set": bson.M{"status": "Dibatalkan"}},
    )
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Booking berhasil dibatalkan"})
}