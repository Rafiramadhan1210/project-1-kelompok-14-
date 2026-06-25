package config

import (
	"fmt"
	"gocroot/helper"
	"gocroot/model"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

var Mongoconn *mongo.Database
var mongoinfo model.DBInfo

func InitDB() {
	_ = godotenv.Load()

	MongoString := os.Getenv("MONGOSTRING")
	log.Printf("DEBUG MONGOSTRING: %s\n", MongoString)
	log.Printf("DEBUG DBNAME: %s\n", os.Getenv("DBNAME")) 
	mongoinfo = model.DBInfo{
		DBString: helper.SRVLookup(MongoString),
		DBName:   os.Getenv("DBNAME"),
	}
	log.Printf("DEBUG resolved DBString: %s\n", mongoinfo.DBString)
	var err error
	Mongoconn, err = helper.MongoConnect(mongoinfo)
	if err != nil {
		panic(fmt.Sprintf("Gagal koneksi ke MongoDB: %v", err))
	}
	log.Println("DEBUG koneksi MongoDB sukses")
}
