package controller

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"log"
	"time"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const sessionCookieName = "session_token"
const sessionDuration = 24 * time.Hour

// generateToken membuat token acak yang aman untuk session login
func generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// RegisterUser untuk mendaftarkan akun baru
func RegisterUser(c *fiber.Ctx) error {
	db := config.Mongoconn
	var user model.Users
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}
	log.Printf("DEBUG register parsed: %+v\n", user)
	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Email dan Password wajib diisi!"})
	}

	// Cek apakah email sudah terdaftar
	var existing model.Users
	err := db.Collection("users").FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&existing)
	if err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "Email sudah terdaftar!"})
	}

	// Hash password
	bytes, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memproses password"})
	}
	user.Password = string(bytes)

	_, err = db.Collection("users").InsertOne(context.Background(), user)
	if err != nil {
		result, err := db.Collection("users").InsertOne(context.Background(), user)
		log.Printf("DEBUG insert result: %+v, err: %v\n", result, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Registrasi Berhasil!"})
}

// LoginUser untuk masuk ke sistem
func LoginUser(c *fiber.Ctx) error {
	db := config.Mongoconn
	var user model.Users
	var dbUser model.Users

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Email dan Password wajib diisi!"})
	}

	// Cari user berdasarkan email
	err := db.Collection("users").FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&dbUser)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email atau Password salah!"})
	}

	// Bandingkan password dengan hash yang tersimpan
	if err := bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Email atau Password salah!"})
	}

	// Buat session token baru
	token, err := generateToken()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal membuat session"})
	}

	now := time.Now()
	session := model.Session{
		Token:     token,
		Email:     dbUser.Email,
		Nama:      dbUser.Nama,
		CreatedAt: primitive.NewDateTimeFromTime(now),
		ExpiresAt: primitive.NewDateTimeFromTime(now.Add(sessionDuration)),
	}

	if _, err := db.Collection("sessions").InsertOne(context.Background(), session); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menyimpan session"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    token,
		Expires:  now.Add(sessionDuration),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login Berhasil!",
		"user": fiber.Map{
			"nama":  dbUser.Nama,
			"email": dbUser.Email,
		},
	})
}

// UpdateUserProfile untuk memperbarui nama dan nomor HP user
func UpdateUserProfile(c *fiber.Ctx) error {
	db := config.Mongoconn
	token := c.Cookies(sessionCookieName)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Belum login"})
	}

	var session model.Session
	err := db.Collection("sessions").FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session tidak valid"})
	}

	if session.ExpiresAt.Time().Before(time.Now()) {
		_, _ = db.Collection("sessions").DeleteOne(context.Background(), bson.M{"token": token})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session telah berakhir"})
	}

	var updatedUser model.Users
	if err := c.BodyParser(&updatedUser); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	// Tambahkan validasi untuk Nama
	if updatedUser.Nama == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Nama tidak boleh kosong"})
	}

	// Pastikan hanya nama dan phone yang diupdate
	updateFields := bson.M{}
	if updatedUser.Nama != "" {
		updateFields["nama"] = updatedUser.Nama
	}
	if updatedUser.Phone != "" {
		updateFields["phone"] = updatedUser.Phone
	}

	if len(updateFields) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Tidak ada data yang diperbarui"})
	}

	_, err = db.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"email": session.Email},
		bson.M{"$set": updateFields},
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memperbarui profil"})
	}

	// Update session name if name was updated
	if updatedUser.Nama != "" {
		_, err = db.Collection("sessions").UpdateMany(
			context.Background(),
			bson.M{"email": session.Email},
			bson.M{"$set": bson.M{"nama": updatedUser.Nama}},
		)
		if err != nil {
			log.Printf("WARNING: Gagal memperbarui nama di session: %v", err)
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profil berhasil diperbarui"})
}

// ChangeUserPassword untuk mengubah password user
func ChangeUserPassword(c *fiber.Ctx) error {
	db := config.Mongoconn
	token := c.Cookies(sessionCookieName)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Belum login"})
	}

	var session model.Session
	err := db.Collection("sessions").FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session tidak valid"})
	}

	if session.ExpiresAt.Time().Before(time.Now()) {
		_, _ = db.Collection("sessions").DeleteOne(context.Background(), bson.M{"token": token})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session telah berakhir"})
	}

	var req struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()})
	}

	if req.CurrentPassword == "" || req.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Kata sandi saat ini dan kata sandi baru wajib diisi"})
	}

	var user model.Users
	err = db.Collection("users").FindOne(context.Background(), bson.M{"email": session.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data user"})
	}

	// Verifikasi kata sandi saat ini
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.CurrentPassword)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Kata sandi saat ini salah"})
	}

	// Hash kata sandi baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal memproses kata sandi baru"})
	}

	_, err = db.Collection("users").UpdateOne(
		context.Background(),
		bson.M{"email": session.Email},
		bson.M{"$set": bson.M{"password": string(hashedPassword)}},
	)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengubah kata sandi"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kata sandi berhasil diubah"})
}

// DeleteUserAccount untuk menghapus akun user
func DeleteUserAccount(c *fiber.Ctx) error {
	db := config.Mongoconn
	token := c.Cookies(sessionCookieName)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Belum login"})
	}

	var session model.Session
	err := db.Collection("sessions").FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session tidak valid"})
	}

	if session.ExpiresAt.Time().Before(time.Now()) {
		_, _ = db.Collection("sessions").DeleteOne(context.Background(), bson.M{"token": token})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session telah berakhir"})
	}

	// Hapus user dari koleksi users
	_, err = db.Collection("users").DeleteOne(context.Background(), bson.M{"email": session.Email})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menghapus akun user"})
	}

	// Hapus semua session user
	_, err = db.Collection("sessions").DeleteMany(context.Background(), bson.M{"email": session.Email})
	if err != nil {
		log.Printf("WARNING: Gagal menghapus session user: %v", err)
	}

	// Hapus cookie di sisi client
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Akun berhasil dihapus"})
}

// LogoutUser untuk keluar dari sistem
func LogoutUser(c *fiber.Ctx) error {
	db := config.Mongoconn
	token := c.Cookies(sessionCookieName)

	if token != "" {
		_, _ = db.Collection("sessions").DeleteOne(context.Background(), bson.M{"token": token})
	}

	// Hapus cookie di sisi client
	c.Cookie(&fiber.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		SameSite: "Lax",
		Path:     "/",
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Logout Berhasil!"})
}

// GetCurrentUser mengembalikan data user yang sedang login berdasarkan session cookie
func GetCurrentUser(c *fiber.Ctx) error {
	db := config.Mongoconn
	token := c.Cookies(sessionCookieName)
	if token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Belum login"})
	}

	var session model.Session
	err := db.Collection("sessions").FindOne(context.Background(), bson.M{"token": token}).Decode(&session)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session tidak valid"})
	}

	if session.ExpiresAt.Time().Before(time.Now()) {
		_, _ = db.Collection("sessions").DeleteOne(context.Background(), bson.M{"token": token})
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Session telah berakhir"})
	}

	var user model.Users
	err = db.Collection("users").FindOne(context.Background(), bson.M{"email": session.Email}).Decode(&user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data user dari database"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"nama":  user.Nama,
			"email": user.Email,
			"phone": user.Phone,
		},
	})
}

func UpdateProfile(c *fiber.Ctx) error {
	db := config.Mongoconn
	// Contoh sederhana: ambil email dari session (sesuaikan dengan cara kamu simpan session)
	// Untuk contoh ini, kita asumsi kirim email lewat body/header untuk identifikasi user
	type RequestUpdate struct {
		Nama  string `json:"nama"`
		Phone string `json:"phone"`
	}

	var req RequestUpdate
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Format data salah"})
	}

	// Update ke database
	// Pastikan filter email sesuai dengan user yang sedang login
	filter := bson.M{"email": c.Locals("email")} // Asumsi session menyimpan email di Locals
	update := bson.M{"$set": bson.M{"nama": req.Nama, "phone": req.Phone}}

	_, err := db.Collection("users").UpdateOne(context.Background(), filter, update)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Gagal update database"})
	}

	return c.JSON(fiber.Map{"message": "Profil berhasil diperbarui"})
}


