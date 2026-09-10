package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ngwebwork/payflow/internal/config"
	"github.com/ngwebwork/payflow/internal/database"
	"go.mongodb.org/mongo-driver/bson"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	client, err := database.NewMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())
	db := client.Database(cfg.MongoDatabase)
	email := "e2e.tester@payflow.local"
	var user struct {
		ID any `bson:"_id"`
	}
	if err := db.Collection("users").FindOne(ctx, bson.M{"email": email}).Decode(&user); err != nil {
		fmt.Println("no test user found, nothing to clean up")
		return
	}
	uid := user.ID
	for _, collection := range []string{"payments", "transactions", "customers"} {
		res, err := db.Collection(collection).DeleteMany(ctx, bson.M{"userId": uid})
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("deleted", res.DeletedCount, "from", collection)
	}
	res, err := db.Collection("users").DeleteOne(ctx, bson.M{"email": email})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("deleted", res.DeletedCount, "test user")
}
