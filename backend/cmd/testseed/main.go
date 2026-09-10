package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ngwebwork/payflow/internal/config"
	"github.com/ngwebwork/payflow/internal/database"
	"github.com/ngwebwork/payflow/internal/models"
	"github.com/ngwebwork/payflow/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	hash, err := utils.HashPassword("TestPass123!")
	if err != nil {
		log.Fatal(err)
	}
	users := db.Collection("users")
	email := "e2e.tester@payflow.local"
	_, _ = users.DeleteOne(ctx, bson.M{"email": email})
	now := time.Now().UTC()
	user := models.User{ID: primitive.NewObjectID(), FirstName: "E2E", LastName: "Tester", Email: email, Password: hash, EmailVerified: true, IsActive: true, CreatedAt: now, UpdatedAt: now}
	if _, err := users.InsertOne(ctx, user); err != nil {
		log.Fatal(err)
	}
	fmt.Println("seeded test user:", user.ID.Hex(), email)
}
