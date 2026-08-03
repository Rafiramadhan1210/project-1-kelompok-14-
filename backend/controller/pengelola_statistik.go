package controller

import (
	"context"
	"sort"
	"time"

	"gocroot/config"
	"gocroot/model"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

type pendapatanBulanan struct {
	Bulan string `json:"bulan"` // format "2026-07"
	Total int    `json:"total"`
}

type pengunjungBulanan struct {
	Bulan  string `json:"bulan"` // format "2026-07"
	Jumlah int    `json:"jumlah"`
}

type bookingPerDestinasi struct {
	Destinasi string `json:"destinasi"`
	Jumlah    int    `json:"jumlah"`
}

// GetPengelolaStatistik menghitung data untuk dashboard statistik pengelola:
// total pendapatan, pendapatan per bulan, jumlah booking per destinasi,
// dan tren jumlah pengunjung per bulan (berdasarkan tanggal kunjungan).
func GetPengelolaStatistik(c *fiber.Ctx) error {
	db := config.Mongoconn
	email := c.Locals("pengelolaEmail").(string)

	names, err := namaDestinasiMilikPengelola(email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	if len(names) == 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"total_pendapatan":       0,
			"pendapatan_per_bulan":   []pendapatanBulanan{},
			"booking_per_destinasi":  []bookingPerDestinasi{},
			"tren_pengunjung_bulan":  []pengunjungBulanan{},
			"total_booking":          0,
		})
	}

	var bookings []model.Booking
	cursor, err := db.Collection("bookings").Find(context.Background(), bson.M{"destination": bson.M{"$in": names}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}
	if err := cursor.All(context.Background(), &bookings); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	// Status yang dihitung sebagai pendapatan nyata (booking yang sudah pasti bayar)
	pendapatanStatus := map[string]bool{"Dibayar": true, "Checked-in": true, "Selesai": true}

	totalPendapatan := 0
	pendapatanBulanMap := map[string]int{}
	pengunjungBulanMap := map[string]int{}
	destinasiMap := map[string]int{}

	for _, b := range bookings {
		destinasiMap[b.Destination]++

		if pendapatanStatus[b.Status] {
			totalPendapatan += b.TotalBayar

			bulanKey := "-"
			if b.CreatedAt.Time().Unix() > 0 {
				bulanKey = b.CreatedAt.Time().Format("2006-01")
			}
			if bulanKey != "-" {
				pendapatanBulanMap[bulanKey] += b.TotalBayar
			}

			// Tren pengunjung dihitung dari tanggal kunjungan (kapan mereka datang),
			// bukan kapan booking dibuat.
			if tgl, errParse := time.Parse("2006-01-02", b.TanggalKunjungan); errParse == nil {
				kunjunganKey := tgl.Format("2006-01")
				pengunjungBulanMap[kunjunganKey] += b.TotalTiket
			}
		}
	}

	pendapatanPerBulan := make([]pendapatanBulanan, 0, len(pendapatanBulanMap))
	for bulan, total := range pendapatanBulanMap {
		pendapatanPerBulan = append(pendapatanPerBulan, pendapatanBulanan{Bulan: bulan, Total: total})
	}
	sort.Slice(pendapatanPerBulan, func(i, j int) bool { return pendapatanPerBulan[i].Bulan < pendapatanPerBulan[j].Bulan })

	trenPengunjung := make([]pengunjungBulanan, 0, len(pengunjungBulanMap))
	for bulan, jumlah := range pengunjungBulanMap {
		trenPengunjung = append(trenPengunjung, pengunjungBulanan{Bulan: bulan, Jumlah: jumlah})
	}
	sort.Slice(trenPengunjung, func(i, j int) bool { return trenPengunjung[i].Bulan < trenPengunjung[j].Bulan })

	bookingPerDest := make([]bookingPerDestinasi, 0, len(destinasiMap))
	for dest, jumlah := range destinasiMap {
		bookingPerDest = append(bookingPerDest, bookingPerDestinasi{Destinasi: dest, Jumlah: jumlah})
	}
	sort.Slice(bookingPerDest, func(i, j int) bool { return bookingPerDest[i].Jumlah > bookingPerDest[j].Jumlah })

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"total_pendapatan":      totalPendapatan,
		"total_booking":         len(bookings),
		"pendapatan_per_bulan":  pendapatanPerBulan,
		"booking_per_destinasi": bookingPerDest,
		"tren_pengunjung_bulan": trenPengunjung,
	})
}
