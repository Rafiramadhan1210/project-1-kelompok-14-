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

	_ = CreateNotification(
		email,
		"booking",
		"Booking Berhasil Dibuat",
		fmt.Sprintf("Booking kamu untuk %s sedang menunggu konfirmasi pembayaran.", destination),
		"booking-saya.html",
	)

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
		BookingID       string `json:"booking_id"`
		Status          string `json:"status"`
		KeteranganTolak string `json:"keterangan_tolak"`
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

	// Ambil data booking dulu sebelum diupdate, untuk keperluan notifikasi
	var existingBooking model.Booking
	_ = db.Collection("bookings").FindOne(context.Background(), bson.M{"_id": oid}).Decode(&existingBooking)

	updateFields := bson.M{"status": body.Status}
	if body.Status == "Dibatalkan" {
		updateFields["keterangan_tolak"] = body.KeteranganTolak
	} else {
		// status lain dianggap bukan penolakan, bersihkan keterangan lama
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

	if existingBooking.Email != "" {
		dibatalkanMsg := fmt.Sprintf("Maaf, booking %s telah dibatalkan oleh admin.", existingBooking.Destination)
		if body.Status == "Dibatalkan" && body.KeteranganTolak != "" {
			dibatalkanMsg = fmt.Sprintf("Booking %s dibatalkan. Alasan: %s", existingBooking.Destination, body.KeteranganTolak)
		}
		statusMessage := map[string]string{
			"Dibayar":    fmt.Sprintf("Pembayaran booking %s telah dikonfirmasi. Selamat berlibur!", existingBooking.Destination),
			"Selesai":    fmt.Sprintf("Booking %s telah selesai. Terima kasih telah menggunakan GoTrip!", existingBooking.Destination),
			"Dibatalkan": dibatalkanMsg,
			"Pending":    fmt.Sprintf("Status booking %s diubah menjadi menunggu konfirmasi.", existingBooking.Destination),
		}
		message, ok := statusMessage[body.Status]
		if !ok {
			message = fmt.Sprintf("Status booking %s diperbarui menjadi %s.", existingBooking.Destination, body.Status)
		}
		_ = CreateNotification(
			existingBooking.Email,
			"booking",
			"Status Booking Diperbarui",
			message,
			"booking-saya.html",
		)
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

// ReuploadBuktiBayar memungkinkan user mengupload ulang bukti bayar
// untuk booking miliknya yang berstatus Dibatalkan (misal karena bukti
// sebelumnya tidak valid). Setelah upload, status balik jadi Pending
// supaya diverifikasi ulang oleh admin.
// Menerima multipart/form-data dengan field "bukti_bayar".
func ReuploadBuktiBayar(c *fiber.Ctx) error {
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

	// Pastikan booking ini milik user yang login dan statusnya Dibatalkan
	var booking model.Booking
	err = db.Collection("bookings").FindOne(context.Background(), bson.M{
		"_id":   oid,
		"email": email,
	}).Decode(&booking)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Booking tidak ditemukan"})
	}

	if booking.Status != "Dibatalkan" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Hanya booking berstatus Dibatalkan yang bisa diupload ulang bukti bayarnya",
		})
	}

	fileHeader, err := c.FormFile("bukti_bayar")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Bukti bayar wajib diupload"})
	}

	uploadDir := filepath.Join(".", "frontend", "public", "uploads", "bukti-bayar")
	if _, statErr := os.Stat(uploadDir); os.IsNotExist(statErr) {
		uploadDir = filepath.Join("..", "frontend", "public", "uploads", "bukti-bayar")
	}
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyiapkan folder upload"})
	}

	ext := filepath.Ext(fileHeader.Filename)
	uniqueName := fmt.Sprintf("%s-%d%s", strings.ReplaceAll(email, "@", "_"), time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, uniqueName)

	if err := c.SaveFile(fileHeader, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan bukti bayar"})
	}

	buktiBayarURL := "/uploads/bukti-bayar/" + uniqueName

	_, err = db.Collection("bookings").UpdateOne(
		context.Background(),
		bson.M{"_id": oid},
		bson.M{"$set": bson.M{
			"bukti_bayar":      buktiBayarURL,
			"status":           "Pending",
			"keterangan_tolak": "",
		}},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	_ = CreateNotification(
		email,
		"booking",
		"Bukti Bayar Diperbarui",
		fmt.Sprintf("Bukti bayar baru untuk booking %s telah dikirim dan menunggu konfirmasi admin.", booking.Destination),
		"pembayaran.html",
	)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":     "Bukti bayar berhasil diupload ulang, menunggu konfirmasi admin.",
		"bukti_bayar": buktiBayarURL,
	})
}