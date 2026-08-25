package main

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/ngwebwork/payflow/internal/config"
	"github.com/ngwebwork/payflow/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
)

func main() {
	_ = godotenv.Load()
	if os.Getenv("APP_ENV") != "development" {
		log.Fatal("seed is disabled unless APP_ENV=development")
	}
	userID, err := primitive.ObjectIDFromHex(os.Getenv("SEED_USER_ID"))
	if err != nil {
		log.Fatal("SEED_USER_ID must be a valid MongoDB ObjectID")
	}
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()
	client, err := database.NewMongo(ctx, cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	if err := database.SeedCollections(ctx, client.Database(cfg.MongoDatabase), userID); err != nil {
		log.Fatal(err)
	}
	log.Println("development dashboard data seeded")
}
