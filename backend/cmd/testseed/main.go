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
	users := db.Collection("users")
	email := "e2e.tester@payflow.local"
	_, err = users.UpdateOne(ctx, bson.M{"email": email}, bson.M{"$set": bson.M{"phone": "+2348011112222", "phoneVerified": true, "updatedAt": time.Now().UTC()}})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("marked phone verified for", email)
}
