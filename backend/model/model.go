package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	Nama     string `json:"nama,omitempty" bson:"nama,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
	Foto     string `json:"foto,omitempty" bson:"foto,omitempty"`
	Wishlist []string `json:"wishlist,omitempty" bson:"wishlist,omitempty"`
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
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama      string             `bson:"nama" json:"nama"`
	Deskripsi string             `bson:"deskripsi" json:"deskripsi"`
	Kategori  string             `bson:"kategori" json:"kategori"`
	Harga     int                `bson:"harga" json:"harga"`
	Gambar    string             `bson:"gambar" json:"gambar"`
	Rating    float64            `bson:"rating" json:"rating"`
	Lokasi    string             `bson:"lokasi" json:"lokasi"`
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
	Status        string             `bson:"status" json:"status"`
	CreatedAt     primitive.DateTime `bson:"created_at" json:"created_at"`
}
type Kategori struct {
    ID    primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
    Nama  string             `bson:"nama" json:"nama"`
    Slug  string             `bson:"slug" json:"slug"` // Opsional, untuk URL friendly
}