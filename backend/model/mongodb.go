package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type DBInfo struct {
	DBString string
	DBName   string
}

type Destinasi struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Nama      string             `bson:"nama" json:"nama"`
	Deskripsi string             `bson:"deskripsi" json:"deskripsi"`
	Kategori  string             `bson:"kategori" json:"kategori"`
	Harga     int                `bson:"harga" json:"harga"`
}
