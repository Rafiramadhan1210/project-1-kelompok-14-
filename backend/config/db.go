package config

import (
	"fmt"
	"gocroot/helper"
	"gocroot/model"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var Mongoconn *mongo.Database
var mongoinfo model.DBInfo

func InitDB() {
	_ = godotenv.Load()

	MongoString := os.Getenv("MONGOSTRING")
	mongoinfo = model.DBInfo{
		DBString: helper.SRVLookup(MongoString),
		DBName:   os.Getenv("DBNAME"),
	}

	var err error
	Mongoconn, err = helper.MongoConnect(mongoinfo)
	if err != nil {
		panic(fmt.Sprintf("Gagal koneksi ke MongoDB: %v", err))
	}
}
