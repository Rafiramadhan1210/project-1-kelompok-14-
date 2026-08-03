package controller

import (
	"context"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const defaultKomisiPersen = 10.0 // dipakai kalau platform_config belum pernah diset

// GetKomisiPersen mengambil persentase komisi platform yang sedang berlaku.
// Kalau belum pernah diset di database, kembalikan nilai default.
func GetKomisiPersen() float64 {
	db := config.Mongoconn
	var cfg model.PlatformConfig
	err := db.Collection("platform_config").FindOne(context.Background(), bson.M{}).Decode(&cfg)
	if err != nil {
		return defaultKomisiPersen
	}
	return cfg.KomisiPersen
}

// GetPlatformConfig mengembalikan konfigurasi platform saat ini (untuk ditampilkan di Admin Panel).
func GetPlatformConfig(c *fiber.Ctx) error {
	db := config.Mongoconn
	var cfg model.PlatformConfig
	err := db.Collection("platform_config").FindOne(context.Background(), bson.M{}).Decode(&cfg)
	if err != nil {
		// belum pernah diset, kembalikan default
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"komisi_persen":          defaultKomisiPersen,
			"total_komisi_terkumpul": 0,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"komisi_persen":          cfg.KomisiPersen,
		"total_komisi_terkumpul": cfg.TotalKomisiTerkumpul,
	})
}

// TambahKomisiTerkumpul menambah "saldo" komisi platform (pencatatan akumulasi,
// bukan transfer uang riil). Dipanggil setiap kali sebuah booking pertama kali
// dikonfirmasi jadi "Dibayar" dan komisinya dihitung.
func TambahKomisiTerkumpul(nominal int) {
	db := config.Mongoconn
	upsert := true
	_, _ = db.Collection("platform_config").UpdateOne(
		context.Background(),
		bson.M{},
		bson.M{"$inc": bson.M{"total_komisi_terkumpul": nominal}},
		&options.UpdateOptions{Upsert: &upsert},
	)
}

// UpdatePlatformConfig mengubah persentase komisi platform. Hanya admin.
// Body JSON: { "komisi_persen": 10 }
func UpdatePlatformConfig(c *fiber.Ctx) error {
	db := config.Mongoconn

	var body struct {
		KomisiPersen float64 `json:"komisi_persen"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Data tidak valid"})
	}
	if body.KomisiPersen < 0 || body.KomisiPersen > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Komisi harus di antara 0-100"})
	}

	upsert := true
	_, err := db.Collection("platform_config").UpdateOne(
		context.Background(),
		bson.M{},
		bson.M{"$set": bson.M{"komisi_persen": body.KomisiPersen}},
		&options.UpdateOptions{Upsert: &upsert},
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Komisi platform berhasil diperbarui", "komisi_persen": body.KomisiPersen})
}

// hitungKomisi menghitung nominal komisi & pendapatan bersih dari total_bayar,
// memakai persentase komisi yang berlaku SAAT DIPANGGIL (dikunci ke booking, tidak
// berubah lagi meskipun persentase komisi platform diubah di kemudian hari).
func hitungKomisi(totalBayar int) (persen float64, nominal int, bersih int) {
	persen = GetKomisiPersen()
	nominal = int(float64(totalBayar) * persen / 100.0)
	bersih = totalBayar - nominal
	return
}