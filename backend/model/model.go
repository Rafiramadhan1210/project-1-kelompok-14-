package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	Nama     string `json:"nama,omitempty" bson:"nama,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Foto     string `json:"foto,omitempty" bson:"foto,omitempty"`
	Wishlist []string `json:"wishlist,omitempty" bson:"wishlist,omitempty"`
	Provider string `json:"provider,omitempty" bson:"provider,omitempty"` // "local" atau "google"
	Role     string `json:"role,omitempty" bson:"role,omitempty"`         // "admin", "pengelola_wisata", atau kosong/"user"
}

type UsersLogin struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty" query:"email" url:"email,omitempty" reqHeader:"email"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

// Session menyimpan token login user yang aktif
type Session struct {
	Token     string             `bson:"token" json:"token"`
	Email     string             `bson:"email" json:"email"`
	Nama      string             `bson:"nama" json:"nama"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
	ExpiresAt primitive.DateTime `bson:"expires_at" json:"expires_at"`
}

type Destinations struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama           string             `bson:"nama" json:"nama"`
	Deskripsi      string             `bson:"deskripsi" json:"deskripsi"`
	Kategori       string             `bson:"kategori" json:"kategori"`
	Harga          int                `bson:"harga" json:"harga"`
	Gambar         string             `bson:"gambar" json:"gambar"`
	Rating         float64            `bson:"rating" json:"rating"`
	Lokasi         string             `bson:"lokasi" json:"lokasi"`
	PengelolaEmail string             `bson:"pengelola_email,omitempty" json:"pengelola_email,omitempty"` // email akun role "pengelola_wisata" pemilik destinasi ini
}

type Booking struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email         string             `bson:"email" json:"email"`
	NamaUser      string             `bson:"nama_user" json:"nama_user"`
	NoHP          string             `bson:"no_hp" json:"no_hp"`
	Destination   string             `bson:"destination" json:"destination"`
	TotalTiket    int                `bson:"total_tiket" json:"total_tiket"`
	TotalBayar    int                `bson:"total_bayar" json:"total_bayar"`
	TanggalKunjungan string          `bson:"tanggal_kunjungan" json:"tanggal_kunjungan"`
	BuktiBayar    string             `bson:"bukti_bayar" json:"bukti_bayar"`
	Status        string             `bson:"status" json:"status"` // "Pending" | "Dibayar" | "Checked-in" | "Selesai" | "Dibatalkan"
	KeteranganTolak string           `bson:"keterangan_tolak,omitempty" json:"keterangan_tolak,omitempty"`
	CreatedAt     primitive.DateTime `bson:"created_at" json:"created_at"`

	// Komisi platform. Dihitung & dikunci sekali saat status pertama kali jadi
	// "Dibayar", memakai persentase komisi yang berlaku SAAT ITU. Tidak dihitung
	// ulang meskipun persentase komisi platform diubah di kemudian hari, supaya
	// riwayat transaksi lama tetap akurat.
	KomisiPersen     float64 `bson:"komisi_persen,omitempty" json:"komisi_persen,omitempty"`
	KomisiNominal    int     `bson:"komisi_nominal,omitempty" json:"komisi_nominal,omitempty"`
	PendapatanBersih int     `bson:"pendapatan_bersih,omitempty" json:"pendapatan_bersih,omitempty"` // total_bayar - komisi_nominal, ini yang jadi hak pengelola
}

// Notification merepresentasikan satu notifikasi.
// Email kosong ("") berarti notifikasi broadcast (promo) untuk semua user.
// ReadBy menyimpan daftar email yang sudah membaca notifikasi ini.
type Notification struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Email     string             `bson:"email" json:"email"`
	Type      string             `bson:"type" json:"type"` // "booking" | "promo"
	Title     string             `bson:"title" json:"title"`
	Message   string             `bson:"message" json:"message"`
	Link      string             `bson:"link,omitempty" json:"link,omitempty"`
	ReadBy    []string           `bson:"read_by" json:"-"`
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

// SupportMessage menyimpan pesan dari form "Hubungi Kami" di halaman bantuan.
type SupportMessage struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama      string             `bson:"nama" json:"nama"`
	Email     string             `bson:"email" json:"email"`
	Topik     string             `bson:"topik" json:"topik"`
	Pesan     string             `bson:"pesan" json:"pesan"`
	Status    string             `bson:"status" json:"status"` // "Baru" | "Dibalas"
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`
}

type Kategori struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Nama  string             `bson:"nama" json:"nama"`
    Slug  string             `bson:"slug" json:"slug"` // Opsional, untuk URL friendly
}

// PlatformConfig menyimpan pengaturan global platform. Hanya ada 1 dokumen
// di collection "platform_config" (semacam singleton). Kalau belum pernah
// diset, dianggap default 10%.
type PlatformConfig struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	KomisiPersen          float64            `bson:"komisi_persen" json:"komisi_persen"`                     // misal 10 artinya 10%
	TotalKomisiTerkumpul  int                `bson:"total_komisi_terkumpul" json:"total_komisi_terkumpul"`   // akumulasi komisi dari semua booking yang sudah Dibayar (pencatatan, bukan uang riil di rekening)
}