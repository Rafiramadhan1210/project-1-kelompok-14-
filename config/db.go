package config
import (
	"gocroot/helper"
	"gocroot/model"
	"os"
	"github.com/joho/godotenv"
)

var _ = godotenv.Load()

var MongoString string = os.Getenv("MONGOSTRING")

var mongoinfo = model.DBInfo{
	DBString: helper.SRVLookup(MongoString),
	DBName:os.Getenv("DBNAME"),
}

var Mongoconn, _ = helper.MongoConnect(mongoinfo)
