//go:build ignore

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gocroot/model"
)

func main() {
	godotenv.Load(".env")
	uri := os.Getenv("MONGOSTRING")
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(uri))
	if err != nil { log.Fatal(err) }
	
	collection := client.Database("gotrip_db").Collection("destinasi")
	
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil { log.Fatal(err) }
	
	var destinations []model.Destinasi
	if err = cursor.All(context.Background(), &destinations); err != nil {
		log.Fatal("ERROR DECODING:", err)
	}
	
	fmt.Printf("Success! Got %d results\n", len(destinations))
}
