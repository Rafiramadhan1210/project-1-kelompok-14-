package model

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	Nama     string `json:"nama,omitempty" bson:"nama,omitempty"`
	Email    string `json:"email,omitempty" bson:"email,omitempty"`
	Phone    string `json:"phone,omitempty" bson:"phone,omitempty"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

type UsersLogin struct {
	Email    string `json:"email,omitempty" bson:"email,omitempty" query:"email" url:"email,omitempty" reqHeader:"email"`
	Password string `json:"password,omitempty" bson:"password,omitempty"`
}

type Destinations struct {
	Nama        string `json:"nama,omitempty" bson:"nama,omitempty"`
	Description string `json:"description,omitempty" bson:"description,omitempty"`
	Price       int    `json:"price,omitempty" bson:"price,omitempty"`
	Category    string `json:"category,omitempty" bson:"category,omitempty"`
}

type Booking struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	NamaUser    string             `bson:"nama_user" json:"nama_user"`
	Destination string             `bson:"destination" json:"destination"`
	TotalTiket  int                `bson:"total_tiket" json:"total_tiket"`
	TotalBayar  int                `bson:"total_bayar" json:"total_bayar"`
	Status      string             `bson:"status" json:"status"`
	CreatedAt   primitive.DateTime `bson:"created_at" json:"created_at"`
}
