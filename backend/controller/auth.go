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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": fiber.Map{
			"nama":  session.Nama,
			"email": session.Email,
		},
	})
}
